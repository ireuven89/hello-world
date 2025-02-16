package migration

import (
	"sync"

	"github.com/labstack/gommon/log"

	"github.com/ireuven89/hello-world/backend/migration/model"
)

type Service struct {
	tasks      []model.MigrationTask
	numThreads int
}

func NewMigrationService(numThreads int) *Service {
	return &Service{numThreads: numThreads}
}

func (s *Service) AddTask(task model.MigrationTask) {
	s.tasks = append(s.tasks, task)
}

func (s *Service) worker(tasks <-chan model.MigrationTask, errCh chan<- error, wg *sync.WaitGroup, fn func(t model.MigrationTask) error) {
	defer wg.Done()
	for task := range tasks {
		if err := fn(task); err != nil {
			errCh <- err
		} else {
			log.Printf("Processed migration task: %s\n", task.Name)
		}
	}
}

func (s *Service) execute(fn func(t model.MigrationTask) error) error {
	tasks := make(chan model.MigrationTask, len(s.tasks))
	errCh := make(chan error, len(s.tasks))
	var wg sync.WaitGroup

	for i := 0; i < s.numThreads; i++ {
		wg.Add(1)
		go s.worker(tasks, errCh, &wg, fn)
	}

	for _, task := range s.tasks {
		tasks <- task
	}
	close(tasks)

	wg.Wait()
	close(errCh)

	for err := range errCh {
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Service) Up() error {
	return s.execute(func(t model.MigrationTask) error {
		return t.Execute()
	})
}

func (s *Service) Down() error {
	for _, task := range s.tasks {
		if err := task.Rollback(); err != nil {
			return err
		}
		log.Printf("Rolled back migration task: %T\n", task)
	}
	return nil
}
