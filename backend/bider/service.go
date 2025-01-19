package bider

import (
	"go.uber.org/zap"

	"github.com/ireuven89/hello-world/backend/bider/model"
)

type SService interface {
	List(input model.BiddersInput) ([]model.Bidder, error)
	FindOne(uuid string) (model.Bidder, error)
	CreateBidder(input model.BiddersInput) (string, error)
	UpdateBidder(input model.BiddersInput) (string, error)
	Delete(id string) error
}

type BidderRepo interface {
	List(input model.BiddersInput) ([]model.Bidder, error)
	Single(uuid string) (model.Bidder, error)
	Upsert(input model.BiddersInput) (string, error)
	Delete(id string) error
}

type Service struct {
	repo   BidderRepo
	logger *zap.Logger
}

func (s *Service) List(input model.BiddersInput) ([]model.Bidder, error) {
	res, err := s.repo.List(input)

	if err != nil {
		s.logger.Error("bidderService.List failed listing bidders", zap.Error(err))
		return []model.Bidder{}, err
	}

	return res, nil
}
func (s *Service) FindOne(uuid string) (model.Bidder, error) {
	res, err := s.repo.Single(uuid)

	if err != nil {
		s.logger.Error("bidderService.FindOne finding bidder", zap.Error(err))
		return model.Bidder{}, err
	}

	return res, nil
}

func (s *Service) CreateBidder(input model.BiddersInput) (string, error) {
	result, err := s.repo.Upsert(input)

	if err != nil {
		s.logger.Error("BidderService.CreateBidder failed creating bidder", zap.Error(err))
		return "", err
	}

	return result, nil
}
func (s *Service) UpdateBidder(input model.BiddersInput) (string, error) {
	result, err := s.repo.Upsert(input)

	if err != nil {
		s.logger.Error("BidderService.CreateBidder failed creating bidder", zap.Error(err))
		return "", err
	}
	return result, nil
}
func (s *Service) Delete(id string) error {
	s.repo.Delete(id)
	return nil
}
