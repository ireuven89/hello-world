package itemming

import (
	"context"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/sethvargo/go-retry"
	"go.uber.org/zap"

	"github.com/ireuven89/hello-world/backend/itemming/model"
	"github.com/ireuven89/hello-world/backend/utils"
)

type Service interface {
	GetItems(model model.ListInput) ([]model.Item, error)
	GetItem(uuid string) (model.Item, error)
	UpdateItem(item model.ItemInput) error
	CreateItem(item model.ItemInput) (string, error)
	CreateItems(items []model.ItemInput) error
	DeleteItem(uuid string) error
	Health() utils.ServiceHealthCheck
}

type RepositoryItem interface {
	ListItems(input model.ListInput) ([]model.Item, error)
	GetItem(uuid string) (model.Item, error)
	Upsert(input model.ItemInput) (string, error)
	BulkInsert(input []model.ItemInput) error
	Insert(input model.ItemInput) (string, error)
	Update(input model.ItemInput) error
	Delete(uuid string) error
	DBstatus() utils.DbStatus
}

const maxBulk = 500000

type ServiceItem struct {
	repo   RepositoryItem
	logger *zap.Logger
}

func New(repo RepositoryItem, logger *zap.Logger) Service {

	return &ServiceItem{repo: repo, logger: logger}
}

func (s *ServiceItem) Health() utils.ServiceHealthCheck {
	var health utils.ServiceHealthCheck
	health.DBStatus = append(health.DBStatus, s.repo.DBstatus())
	health.ServiceStatus = "UP"

	return health
}

func (s *ServiceItem) GetItems(input model.ListInput) ([]model.Item, error) {
	result, err := s.repo.ListItems(input)

	if err != nil {
		s.logger.Error("failed to execute query", zap.Any("list items", input), zap.Error(err))
		return nil, err
	}

	return result, err
}

func (s *ServiceItem) GetItem(uuid string) (model.Item, error) {
	result, err := s.repo.GetItem(uuid)

	if err != nil {
		s.logger.Error("ServiceItem.GetItem failed getting item", zap.Error(err))
		return model.Item{}, err
	}

	return result, err
}

func (s *ServiceItem) UpdateItem(item model.ItemInput) error {

	if err := s.repo.Update(item); err != nil {
		s.logger.Error("ServiceItem.UpdateItem failed updating item: ", zap.Error(err))
		return err
	}

	return nil
}
func (s *ServiceItem) CreateItem(item model.ItemInput) (string, error) {

	if err := validateMandatoryFields(item); err != nil {
		return "", err
	}

	id, err := s.repo.Insert(item)

	if err != nil {
		s.logger.Error("ServiceItem.CreateItem failed creating item", zap.Any("create itemming", item), zap.Error(err))
		return "", err
	}

	return id, err
}

func validateMandatoryFields(item model.ItemInput) error {

	if item.Name == "" {
		return errors.New("missing mandatory field name")
	}

	if item.Description == "" {
		return errors.New("missing mandatory field description")
	}

	if item.Price == 0 {
		return errors.New("missing mandatory field price")
	}
	return nil
}

// CreateItems - Bulk inserts the
func (s *ServiceItem) CreateItems(items []model.ItemInput) error {
	chunks := splitIntoChunks(items)

	//retry method
	backoff := retry.NewConstant(100 * time.Millisecond)
	maxRetries := retry.WithMaxRetries(5, backoff)
	attemptCount := 1

	//for each chunk do
	for i := 0; i < len(chunks); i++ {
		err := retry.Do(context.Background(), maxRetries, func(ctx context.Context) error {
			s.logger.Info(fmt.Sprintf("inserting chunk %v", i))
			err := s.repo.BulkInsert(items)

			if err != nil {
				s.logger.Warn(fmt.Sprintf("ServiceItem.CreateItems failed creating items retry number %v, error %v", attemptCount, err))
				return retry.RetryableError(err)
			}

			return nil
		})

		if err != nil {
			s.logger.Error("ServiceItem.CreateItems failed to creating items: ", zap.Error(err))
			return err
		}
	}

	return nil
}

func splitIntoChunks(items []model.ItemInput) [][]model.ItemInput {
	rows := math.Ceil(float64(len(items)) / float64(maxBulk))
	slices := make([][]model.ItemInput, int(rows))
	count := 0

	for i := range slices {
		slices[i] = make([]model.ItemInput, 0, maxBulk)
	}

	for i := 0; i < len(items); i++ {
		// Append item to the corresponding slice
		slices[count] = append(slices[count], items[i])

		// If maxBulk is reached, move to the next slice
		if len(slices[count]) == maxBulk {
			count++
		}
	}

	return slices
}

// Delete item - with retry
func (s *ServiceItem) DeleteItem(uuid string) error {
	err := s.repo.Delete(uuid)

	backOff := retry.NewConstant(3 * time.Second)
	maxRetries := retry.WithMaxRetries(3, backOff)
	attemptCount := 1
	err = retry.Do(context.Background(), maxRetries, func(ctx context.Context) error {
		err = s.repo.Delete(uuid)
		if err != nil {
			s.logger.Warn(fmt.Sprintf("failed delete attempt %v uuid %s: %v", attemptCount, uuid, err))
			attemptCount++
			return retry.RetryableError(err)
		}

		return nil
	})

	if err != nil {
		s.logger.Error("failed to delete itemming", zap.Any("delete itemming", uuid), zap.Error(err))
		return err
	}

	return err
}
