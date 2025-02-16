package biddering

import (
	"go.uber.org/zap"

	"github.com/ireuven89/hello-world/backend/biddering/model"
)

type Service interface {
	ListBidders(input model.BiddersInput) ([]model.Bidder, error)
	GetBidder(uuid string) (model.Bidder, error)
	CreateBidder(input model.BidderInput) (string, error)
	UpdateBidder(input model.BidderInput) (string, error)
	DeleteBidder(id string) error
}

type BidderRepo interface {
	List(input model.BiddersInput) ([]model.Bidder, error)
	FindOne(uuid string) (model.Bidder, error)
	Upsert(input model.BidderInput) (string, error)
	Delete(id string) error
}

type BidderService struct {
	repo   BidderRepo
	logger *zap.Logger
}

func NewService(repo BidderRepo, logger *zap.Logger) Service {

	return &BidderService{
		repo:   repo,
		logger: logger,
	}
}

func (s *BidderService) ListBidders(input model.BiddersInput) ([]model.Bidder, error) {
	res, err := s.repo.List(input)

	if err != nil {
		s.logger.Error("bidderService.List failed listing bidders", zap.Error(err))
		return []model.Bidder{}, err
	}

	return res, nil
}

func (s *BidderService) GetBidder(uuid string) (model.Bidder, error) {
	res, err := s.repo.FindOne(uuid)

	if err != nil {
		s.logger.Error("bidderService.FindOne finding bidder", zap.Error(err))
		return model.Bidder{}, err
	}

	return res, nil
}

func (s *BidderService) CreateBidder(input model.BidderInput) (string, error) {
	result, err := s.repo.Upsert(input)

	if err != nil {
		s.logger.Error("BidderService.CreateBidder failed creating bidder", zap.Error(err))
		return "", err
	}

	return result, nil
}
func (s *BidderService) UpdateBidder(input model.BidderInput) (string, error) {
	result, err := s.repo.Upsert(input)

	if err != nil {
		s.logger.Error("BidderService.CreateBidder failed creating bidder", zap.Error(err))
		return "", err
	}
	return result, nil
}

func (s *BidderService) DeleteBidder(id string) error {
	err := s.repo.Delete(id)

	if err != nil {
		s.logger.Error("Delete failed deleting  bidder", zap.Error(err))
		return err
	}
	return nil
}
