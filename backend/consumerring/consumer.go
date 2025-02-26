package consumerring

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"go.uber.org/zap"

	"github.com/ireuven89/hello-world/backend/migration/model"
)

// WorkerService handles task execution
type WorkerService struct {
	consumer   *kafka.Consumer
	numWorkers int
	logger     *zap.Logger
	executor   *Executor
}

type Executor struct {
	name     string
	execute  func(param interface{}) error
	rollback func(param interface{}) error
}

func NewExecutor(name string, execute, rollback func(param interface{}) error) *Executor {
	return &Executor{name: name, execute: execute, rollback: rollback}
}

func NewWorkerService(consumer *kafka.Consumer, executor *Executor, numWorkers int) *WorkerService {
	return &WorkerService{consumer: consumer, numWorkers: numWorkers, executor: executor}
}

func (w *WorkerService) ProcessTasks(ctx context.Context, topic string) {
	err := w.consumer.Subscribe(topic, nil)
	if err != nil {
		w.logger.Error("Failed to subscribe to Kafka topic: %v", zap.Error(err))
	}

	// Task queue (buffered to handle bursts)
	taskQueue := make(chan model.MigrationTask, 100)
	var wg sync.WaitGroup

	// Start worker goroutines
	for i := 0; i < w.numWorkers; i++ {
		wg.Add(1)
		go w.worker(ctx, taskQueue, &wg)
	}

	// Read messages from Kafka and send them to the task queue
	for {
		msg, err := w.consumer.ReadMessage(-1)
		if err != nil {
			w.logger.Error("Kafka error", zap.Error(err))
			continue
		}

		var task model.MigrationTask
		if err := json.Unmarshal(msg.Value, &task); err != nil {
			w.logger.Error("Failed to decode task", zap.Error(err))
			continue
		}

		taskQueue <- task // Send task to the queue
	}

	// Close task queue and wait for workers to finish
	close(taskQueue)
	wg.Wait()
}

// Worker function to process tasks concurrently
func (w *WorkerService) worker(ctx context.Context, taskQueue <-chan model.MigrationTask, wg *sync.WaitGroup) {
	defer wg.Done()
	for task := range taskQueue {
		err := w.executor.execute(task.Params)
		if err != nil {
			w.logger.Error("failed execute params")
			if w.executor.rollback != nil {
				err = w.executor.rollback(task.RollbackParams)
				if err != nil {
					w.logger.Error("failed rollback...", zap.Error(err), zap.Any("task", task.ID))
				}
			}
			task.Status = model.TaskFailed.String()
		}
	}
}
