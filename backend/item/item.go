package item

import (
	"github.com/ireuven89/hello-world/backend/item/model"
	"go.uber.org/zap"
)

type Service interface {
	GetItems(model model.ListInput) ([]model.Item, error)
	GetItem(uuid string) (model.Item, error)
	UpdateItem(item model.ItemInput) error
	CreateItem(item model.ItemInput) (string, error)
	DeleteItem(uuid string) error
}

type RepositoryItem interface {
	ListItems(input model.ListInput) ([]model.Item, error)
	GetItem(uuid string) (model.Item, error)
	Upsert(input model.ItemInput) (string, error)
	Delete(uuid string) error
}

type ServiceItem struct {
	repo   RepositoryItem
	logger *zap.Logger
}

func New(repo RepositoryItem, logger *zap.Logger) Service {

	return &ServiceItem{repo: repo, logger: logger}
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
		s.logger.Error("failed to get query", zap.Any("get items", uuid), zap.Error(err))
		return model.Item{}, err
	}

	return result, err
}
func (s *ServiceItem) UpdateItem(item model.ItemInput) error {
	_, err := s.repo.Upsert(item)

	if err != nil {
		s.logger.Error("failed to update item", zap.Any("update item", item), zap.Error(err))
		return err
	}

	return err
}
func (s *ServiceItem) CreateItem(item model.ItemInput) (string, error) {
	id, err := s.repo.Upsert(item)

	if err != nil {
		s.logger.Error("failed to create item", zap.Any("create item", item), zap.Error(err))
		return "", err
	}

	return id, err
}
func (s *ServiceItem) DeleteItem(uuid string) error {
	err := s.repo.Delete(uuid)

	if err != nil {
		s.logger.Error("failed to delete item", zap.Any("delete item", uuid), zap.Error(err))
		return err
	}

	return err
}
