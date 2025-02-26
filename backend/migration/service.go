package migration

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/sethvargo/go-retry"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"

	"github.com/ireuven89/hello-world/backend/migration/model"
)

type InternalService interface {
	Execute(interface{}) error
	Rollback(interface{}) error
}

type Service struct {
	tasks           []model.MigrationTask
	numThreads      int
	migrationDB     *mongo.Database
	queuesDB        *mongo.Database
	internalService InternalService
	httpclient      *http.Client
	logger          *zap.Logger
}

var stopMigration = false
var mu sync.Mutex
var migrationsCollection = "migrations"
var httpStatues = []int{http.StatusOK, http.StatusNoContent, http.StatusAccepted, http.StatusCreated}
var skipUpdate bool

const (
	batchSize  = 10
	maxRetries = 3
)

func NewService(logger *zap.Logger, migrationsDB *mongo.Database, queuesDB *mongo.Database, internalService InternalService) *Service {
	httpClient := http.Client{Timeout: 10 * time.Millisecond}

	return &Service{
		logger:          logger,
		migrationDB:     migrationsDB,
		queuesDB:        queuesDB,
		internalService: internalService,
		httpclient:      &httpClient,
	}
}

// ProcessTasks -process all tasks
func (s *Service) ProcessTasks(ctx context.Context, migrationName string) {
	processedCount := 0
	pageNumber := int32(0)
	startTime := time.Now()
	migration, err := s.getMigration(ctx, migrationName)

	if err != nil {
		s.logger.Error("failed to start migration", zap.Error(err))
		return
	}

	fmt.Printf("Started migration %s\n", migration.MigrationName)
	s.updateMigration(ctx, migrationName, "status", model.InProgress.String())
	queueName := migration.QueueName

	for !stopMigration {
		tasks := s.getTasks(ctx, queueName, model.Pending.String(), pageNumber, batchSize)
		batchStartTime := time.Now()

		if len(tasks) == 0 {
			fmt.Println("No more tasks to process")
			s.updateMigration(ctx, migrationName, "status", model.Finished.String())
			fmt.Printf("Finished migration total count is %d\n", s.getTotalTasks(ctx, queueName))
			break
		}

		var wg sync.WaitGroup
		for _, task := range tasks {
			wg.Add(1)
			go func(t model.MigrationTask) {
				defer wg.Done()
				t.Status = model.TaskInProgress.String()
				s.updateTask(ctx, &t, queueName)
				if migration.Type == model.Http.String() {
					s.processTaskHttp(ctx, &t, migration)
				} else {
					s.processTaskInternal(ctx, &t, queueName)
				}
			}(task)
		}
		wg.Wait()

		batchEndTime := time.Now()
		processedCount += len(tasks)
		pageNumber++
		timeLeft := s.calcTimeLeft(ctx, batchStartTime, batchEndTime, queueName)
		s.updateMigration(ctx, migrationName, "timeLeft", timeLeft)
		fmt.Printf("Processed %d tasks...\n", processedCount)

		if s.getMigrationStatus(ctx, migrationName) == model.Stopped.String() {
			mu.Lock()
			stopMigration = true
			mu.Unlock()
			fmt.Printf("Stopping migration, processed %d\n", processedCount)
		}
	}

	endTime := time.Now()
	fmt.Printf("Task processing complete, took %d ms\n", endTime.Sub(startTime).Milliseconds())
}

// calcTimeLeft - calculates the time left for the migration to end
func (s *Service) calcTimeLeft(ctx context.Context, startTime time.Time, endTime time.Time, queue string) int64 {
	timeLeft := (endTime.Second() - startTime.Second()) * batchSize

	return int64(timeLeft) * s.getTasksLeft(ctx, queue)
}

// updateMigration - updates the migration
func (s *Service) updateMigration(ctx context.Context, migration, fieldName string, fieldValue interface{}) {
	result, err := s.migrationDB.
		Collection(migrationsCollection).
		UpdateOne(ctx, bson.M{"name": migration}, bson.M{"$set": bson.M{fieldName: fieldValue}})

	if err != nil {
		s.logger.Error("failed updating migration", zap.Error(err))
	}

	if result.ModifiedCount == 0 {
		s.logger.Error("failed updating migration  with field", zap.Any(fieldName, fieldValue))
	}
}

// getMigrationStatus - get the migration status
func (s *Service) getMigrationStatus(ctx context.Context, migrationName string) string {
	var result model.Migration
	err := s.migrationDB.
		Collection(migrationsCollection).
		FindOne(ctx, bson.M{"name": migrationName}).
		Decode(&result)

	if err != nil {
		return ""
	}
	return result.Status
}

// getTasks - get tasks
func (s *Service) getTasks(ctx context.Context, queueName, status string, page, batchSize int32) []model.MigrationTask {
	var result []model.MigrationTask

	cursor, err := s.queuesDB.
		Collection(queueName).
		Find(ctx, bson.M{"status": status}, options.Find().SetBatchSize(batchSize))

	if err != nil {
		s.logger.Error("MigrationService.getTasks failed getting tasks ", zap.Error(err))
		return []model.MigrationTask{}
	}
	err = cursor.All(ctx, &result)

	if err != nil {
		s.logger.Error("failed decoding results into migration tasks", zap.Error(err))
	}

	return result
}

// getTotalTasks - returns total tasks left
func (s *Service) getTotalTasks(ctx context.Context, queue string) int64 {
	count, err := s.queuesDB.
		Collection(queue).
		CountDocuments(ctx, bson.M{})
	if err != nil {
		s.logger.Error("failed counting documents", zap.String("queue", queue), zap.Error(err))
		return 0
	}

	return count
}

// getTasksLeft - returns int64 tasks left
func (s *Service) getTasksLeft(ctx context.Context, queue string) int64 {
	count, err := s.queuesDB.
		Collection(queue).
		CountDocuments(ctx, bson.M{"status": model.TaskPending.String()})
	if err != nil {
		s.logger.Error("failed counting documents", zap.String("queue", queue), zap.Error(err))
		return 0
	}

	return count
}

// processTask - process with max retries
func (s *Service) processTaskHttp(ctx context.Context, task *model.MigrationTask, migration model.Migration) {
	//execute
	err := retry.Do(ctx, retry.WithMaxRetries(maxRetries, retry.NewConstant(500*time.Millisecond)), func(ctx context.Context) error {
		body, err := json.Marshal(task.HttpBody)
		if err != nil {
			return retry.RetryableError(err)
		}

		req, err := http.NewRequest(migration.HttpMethod, s.buildUrl(migration.HttpEndpoint, task.HttpParams), bytes.NewReader(body))
		if err != nil {
			return retry.RetryableError(err)
		}
		resp, err := s.httpclient.Do(req)

		if err != nil {
			s.logger.Error("failed executing task retry", zap.Error(err), zap.String("task", task.Name))
			return retry.RetryableError(err)
		}

		if !slices.Contains(httpStatues, resp.StatusCode) {
			s.logger.Error("failed executing task ", zap.String("task", task.Name), zap.Any("status code", resp.StatusCode))
			return retry.RetryableError(fmt.Errorf("failed execute with status code: %v", resp.StatusCode))
		}

		//if finished ok - update completed
		task.Status = model.TaskCompleted.String()
		s.updateTask(ctx, task, migration.QueueName)
		return nil
	})

	//if failed  - rollback
	if err != nil {
		task.ErrorMessage = err.Error()
		task.Status = model.TaskFailed.String()
		s.updateTask(ctx, task, migration.QueueName)
		if migration.HttpRollBack != "" {
			req, err := http.NewRequest(migration.HttpRollBackMethod, migration.HttpRollBack, nil)
			if err != nil {
				s.logger.Warn("failed rollback")
			}
			resp, err := s.httpclient.Do(req)
			if err != nil || !slices.Contains(httpStatues, resp.StatusCode) {
				s.logger.Warn("failed rollback task")
			}
		}
		s.logger.Error("failed executing task", zap.Error(err))
		return
	}

}

func (s *Service) buildUrl(endpoint string, pathParams []string) string {
	var url strings.Builder

	if len(pathParams) == 0 {
		return endpoint
	}

	url.WriteString(endpoint)
	for _, param := range pathParams {
		url.WriteString("/" + param)
	}

	return url.String()
}

func (s *Service) processTaskInternal(ctx context.Context, task *model.MigrationTask, queue string) {
	err := retry.Do(ctx, retry.WithMaxRetries(maxRetries, retry.NewConstant(10*time.Millisecond)), func(ctx context.Context) error {

		if err := s.internalService.Execute(task.Params); err != nil {
			s.logger.Error("failed executing task", zap.Error(err), zap.String("name", task.Name), zap.Any("params", task.Params))
			return retry.RetryableError(err)
		}

		//update task completed - if finished
		task.Status = model.TaskCompleted.String()
		s.updateTask(ctx, task, queue)
		return nil
	})

	if err != nil {
		err = s.internalService.Rollback(task.RollbackParams)
		task.Status = model.TaskFailed.String()
		s.updateTask(ctx, task, queue)
		task.ErrorMessage = err.Error()
		if err != nil {
			s.logger.Error("failed rollback", zap.String("task", task.Name), zap.Error(err), zap.Any("params", task.RollbackParams))
			task.RollbackErrorMessage = err.Error()
		}
	}
}

// updateTask - task with status
func (s *Service) updateTask(ctx context.Context, task *model.MigrationTask, queueName string) {
	//in case tasks are not in mongo
	if skipUpdate {
		return
	}

	updatedResult, err := s.queuesDB.
		Collection(queueName).
		UpdateMany(ctx, bson.M{"_id": task.ID}, bson.M{"$set": bson.M{"status": task.Status}})

	if err != nil {
		s.logger.Error("failed updating task", zap.Error(err), zap.String("name", task.Name))
	}

	if updatedResult.ModifiedCount == 0 {
		s.logger.Error("task not found")
	}
}

func (s *Service) setSkipUpdate(skip bool) {
	skipUpdate = skip
}

func (s *Service) stopMigration(ctx context.Context, name string) {
	s.updateMigration(ctx, name, "status", model.Stopped.String())
}

// getMigration - fetch migration from DB  if exists
func (s *Service) getMigration(ctx context.Context, name string) (model.Migration, error) {
	var migration model.Migration

	dbMigration := s.migrationDB.
		Collection(migrationsCollection).
		FindOne(ctx, bson.M{"name": name})
	if dbMigration == nil {
		s.logger.Error("migration not found")
		return model.Migration{}, errors.New("migration not exists")
	}
	err := dbMigration.Decode(&migration)

	if err != nil {
		return model.Migration{}, err
	}

	return migration, nil
}

func (s *Service) createMigration(ctx context.Context, migration model.Migration) error {
	result, err := s.migrationDB.
		Collection(migrationsCollection).
		InsertOne(ctx, bson.M{
			"name":       migration.MigrationName,
			"status":     model.Pending.String(),
			"type":       migration.Type,
			"queueName":  migration.QueueName,
			"created_at": time.Now(),
		})

	if err != nil {
		s.logger.Error("failed inserting migration")
		return err
	}

	if result.InsertedID == "" {
		return errors.New("failed inserting migration")
	}

	return nil
}

func (s *Service) prepareTasks() {
	page := int64(0)
	pageSize := int64(50)
	var currentId string

	cursor, err := s.queuesDB.Collection("test-queue").
		Find(context.Background(), bson.M{}, options.Find().
			SetLimit(pageSize).
			SetSkip(page).
			SetProjection(bson.M{"_id": 1}))

	if err != nil {
		s.logger.Error("failed getting tasks")
		return
	}

	for cursor.Next(context.Background()) {
		if err = cursor.Decode(&currentId); err != nil {
			fmt.Printf("failed decoding...")
		}
		var task model.MigrationTask
		task.Name = currentId
		task.Params = currentId
		task.Status = model.TaskPending.String()
	}
}
