package userring

import (
	"go.uber.org/zap"

	"github.com/ireuven89/hello-world/backend/userring/model"
)

type Service interface {
	ListUsers(input model.UserFetchInput) ([]model.User, error)
	GetUser(uuid string) (model.User, error)
	CreateUser(input model.UserUpsertInput) (string, error)
	UpdateUser(input model.UserUpsertInput) error
	DeleteUser(uuid string) error
}

type UserRepository interface {
	ListUsers(input model.UserFetchInput) ([]model.User, error)
	FindUser(uuid string) (model.User, error)
	Upsert(input model.UserUpsertInput) (string, error)
	Delete(uuid string) error
}

type service struct {
	logger         *zap.Logger
	userRepository UserRepository
}

func New(logger *zap.Logger, repo UserRepository) Service {

	return &service{
		logger:         logger,
		userRepository: repo,
	}
}

func (s *service) ListUsers(input model.UserFetchInput) ([]model.User, error) {
	result, err := s.userRepository.ListUsers(input)

	if err != nil {
		s.logger.Error("failed to retrieve userring", zap.Error(err))
		return nil, err
	}

	return result, nil
}
func (s *service) GetUser(uuid string) (model.User, error) {
	result, err := s.userRepository.FindUser(uuid)

	if err != nil {
		s.logger.Error("failed to retrieve userring", zap.Error(err))
		return model.User{}, err
	}

	return result, nil
}
func (s *service) CreateUser(input model.UserUpsertInput) (string, error) {
	result, err := s.userRepository.Upsert(input)

	if err != nil {
		s.logger.Error("failed to retrieve userring", zap.Error(err))
		return "", err
	}

	return result, nil
}

func (s *service) UpdateUser(input model.UserUpsertInput) error {
	_, err := s.userRepository.Upsert(input)

	if err != nil {
		s.logger.Error("failed to retrieve userring", zap.Error(err))
		return err
	}

	return nil
}

func (s *service) DeleteUser(uuid string) error {
	if err := s.userRepository.Delete(uuid); err != nil {
		s.logger.Error("failed to retrieve userring", zap.Error(err))
		return err
	}

	return nil
}
