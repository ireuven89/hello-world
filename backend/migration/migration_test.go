package migration

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/labstack/gommon/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/testcontainers/testcontainers-go/modules/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap/zaptest"

	"github.com/ireuven89/hello-world/backend/migration/model"
	"github.com/ireuven89/hello-world/backend/utils"
)

func TestMigration(t *testing.T) {
	client, container, err := setupInMemoryMongoDB()
	logger := zaptest.NewLogger(t)
	queueName := "test-queue"
	migrationsDB := client.Database("migrations")
	queuesDB := client.Database("queues")
	queuesDB.Collection(queueName).InsertMany(context.Background(), []interface{}{
		bson.M{"name": queueName, "status": model.Pending.String(), "taskID": "task1"},
		bson.M{"name": queueName, "status": model.Pending.String(), "taskID": "task2"},
	})
	if err != nil {
		t.Fatalf("failed initailize db %v", err)
	}

	defer container.Terminate(context.Background())

	migrationsDB.Collection("migrations").InsertOne(context.Background(),
		bson.M{"name": "test-migration", "type": model.Http.String(), "queueName": queueName, "status": model.Pending.String()})

	service := NewService(logger, migrationsDB, queuesDB, nil)

	tasks := service.getTasks(context.Background(), queueName, model.Pending.String(), 0, batchSize)
	tasks[0].Status = model.TaskCompleted.String()
	service.updateTask(context.Background(), &tasks[0], queueName)

	left := service.getTasksLeft(context.Background(), "test-queue")
	total := service.getTotalTasks(context.Background(), "test-queue")

	assert.NotEmpty(t, tasks)
	assert.NotEqual(t, 0, left)
	assert.NotEqual(t, 0, total)

}

func setupInMemoryMongoDB() (*mongo.Client, *mongodb.MongoDBContainer, error) {
	contianer, err := mongodb.Run(context.Background(), "mongo:latest")

	if err != nil {
		return nil, nil, err
	}

	uri, err := contianer.ConnectionString(context.Background())

	if err != nil {
		return nil, nil, err
	}

	opts := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.Background(), opts)
	if err != nil {
		return nil, nil, err
	}
	return client, contianer, nil
}

func TestFailedConnectToMongo(t *testing.T) {
	_, err := MustNewDB(utils.DataBaseConnection{Host: "empty-host"})

	assert.NotNil(t, err)
}

func TestProcessTasksExternal(t *testing.T) {
	// Set up in-memory MongoDB client
	client, container, err := setupInMemoryMongoDB()
	if err != nil {
		t.Fatalf("Failed to set up in-memory MongoDB: %v", err)
	}
	defer container.Terminate(context.Background())

	// Use a specific database and collection for testing
	migrationDB := client.Database("migrations")
	queuesDB := client.Database("queues")
	queueCollection := queuesDB.Collection("test-queue")
	migrationCollection := migrationDB.Collection("migrations")
	migration := "test-migration"
	queueName := "test-queue"
	logger := zaptest.NewLogger(t)

	// Create a new Service instance with the in-memory MongoDB client
	service := NewService(logger, migrationDB, queuesDB, nil)

	// Insert some test data into the MongoDB collection
	_, err = queueCollection.InsertMany(context.Background(), []interface{}{
		bson.M{"queueName": queueName, "status": model.Pending.String(), "taskID": "task1"},
		bson.M{"queueName": queueName, "status": model.Pending.String(), "taskID": "task2"},
	})
	id, err := migrationCollection.InsertOne(context.Background(), bson.M{"name": migration, "status": model.Pending.String(), "queueName": "test-queue", "type": model.Http.String()})
	if err != nil {
		t.Fatalf("Failed to insert test data: %v", err)
	}

	fmt.Printf("inserted id: %v\n ", id)

	// Define expectations
	processedCount := 0
	// Start the process tasks function
	service.ProcessTasks(context.Background(), migration)

	// Allow some time for the goroutines to process
	time.Sleep(2 * time.Second)

	// Verify the processed task count and check the status updates
	processedCount += 2 // We know that we have inserted 2 tasks

	// Verify the tasks have been processed
	taskCount1, err := queueCollection.CountDocuments(context.Background(), bson.M{})
	tasks, err := queueCollection.Find(context.Background(), bson.M{})
	taskCount, err := queueCollection.CountDocuments(context.Background(), bson.M{"status": model.TaskFailed.String()})
	if err != nil {
		t.Fatalf("Failed to count InProgress tasks: %v", err)
	}

	var parsedTasks []model.MigrationTask
	tasks.All(context.Background(), &parsedTasks)
	assert.True(t, taskCount1 > 0)
	assert.Equal(t, taskCount, int64(processedCount), "Tasks should be marked as 'InProgress'")

	// Verify migration status was updated
	var migrationEntity model.Migration
	err = migrationDB.Collection("migrations").FindOne(context.Background(), bson.M{"name": migration}).Decode(&migrationEntity)
	if err != nil {
		t.Fatalf("Failed to get migration status: %v", err)
	}
	assert.Equal(t, migrationEntity.Status, model.Finished.String(), "Migration should be marked as Finished")

	// Clean up
	_, err = queueCollection.DeleteMany(context.Background(), bson.M{"queueName": queueName})
	if err != nil {
		t.Fatalf("Failed to clean up test data: %v", err)
	}
}

type MockInternalService struct {
	mock mock.Mock
}

func NewIMockInternalService() *MockInternalService {

	return &MockInternalService{}
}

func TestExecute(t *testing.T) {

	logger := zaptest.NewLogger(t)
	task := model.MigrationTask{
		Name: "print",
		Params: map[string]interface{}{
			"fail":  false,
			"print": []string{"string 1", "string 2"},
		},
		RollbackParams: map[string]interface{}{
			"fail":  false,
			"print": []string{"rollback string 1", "rollback string 2"},
		},
	}

	ctx := context.Background()
	service := NewService(logger, nil, nil, NewIMockInternalService())
	service.setSkipUpdate(true)
	service.processTaskInternal(ctx, &task, "")
}

type internalTests struct {
	name           string
	task           model.MigrationTask
	executeCalled  int
	rollbackCalled int
}

var rollbackCalled int
var executeCalled int

func TestInternalServiceExecute(t *testing.T) {
	tests := []internalTests{
		{
			name: "success execute",
			task: model.MigrationTask{
				Name: "print",
				Params: map[string]interface{}{
					"fail":   false,
					"print":  []string{"string 1", "string 2"},
					"called": 0,
				},
				RollbackParams: nil,
			},
			executeCalled:  1,
			rollbackCalled: 0,
		},
		{
			name: "failed execute",
			task: model.MigrationTask{
				Name: "print",
				Params: map[string]interface{}{
					"fail":  true,
					"print": []string{"string 1", "string 2"},
				},
				RollbackParams: map[string]interface{}{
					"fail":  false,
					"print": []string{"string 1", "string 2"},
				},
			},
			executeCalled:  4,
			rollbackCalled: 1,
		},
	}

	logger := zaptest.NewLogger(t)

	service := NewService(logger, nil, nil, NewIMockInternalService())
	service.setSkipUpdate(true)
	ctx := context.Background()

	for _, test := range tests {
		rollbackCalled = 0
		executeCalled = 0
		service.processTaskInternal(ctx, &test.task, "")
		assert.Equal(t, executeCalled, test.executeCalled, "execute call failed")
		assert.Equal(t, rollbackCalled, test.rollbackCalled, "rollback call time failed")
	}

}

func (mis *MockInternalService) Execute(print interface{}) error {
	stringParams, ok := print.(map[string]interface{})
	executeCalled++
	if !ok {
		log.Fatalf("failed parsing")
	}

	if stringParams["fail"].(bool) {
		return errors.New("failing execute")
	}

	for _, param := range stringParams["print"].([]string) {
		fmt.Printf("execute %s\n", param)
	}

	return nil
}

func (mis *MockInternalService) Rollback(print interface{}) error {
	stringParams, ok := print.(map[string]interface{})
	rollbackCalled++

	if !ok {
		log.Fatalf("failed parsing")
	}

	if stringParams["fail"].(bool) {
		return errors.New("failing rollback")
	}
	for _, param := range stringParams["print"].([]string) {
		fmt.Printf("rollabck %vֿ\n", param)
	}

	return nil
}

func TestUrlCreate(t *testing.T) {
	service := NewService(nil, nil, nil, nil)
	endpoint := "http://endpoint:1000"
	params := []string{"param1", "param2"}
	url := service.buildUrl(endpoint, params)

	assert.Equal(t, "http://endpoint:1000/param1/param2", url)
}
