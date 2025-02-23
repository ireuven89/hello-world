package auctioning

import (
	"go.uber.org/zap"

	"github.com/ireuven89/hello-world/backend/auction/model"
	"github.com/ireuven89/hello-world/backend/utils"
)

type Service interface {
	Search(req model.AuctionRequest) ([]model.Auction, error)
	Find(uuid string) (model.Auction, error)
	Delete(uuid string) error
	Update(req model.AuctionRequest) error
	Health() utils.ServiceHealthCheck
	Run()
}

type AuctRepo interface {
	FindAll(req model.AuctionRequest) ([]model.Auction, error)
	FindOne(Uuid string) (model.Auction, error)
	Update(req model.AuctionRequest) error
	Delete(uuid string) error
	DbStatus() utils.DbStatus
}

// AuctionService is the core authenticating service
type AuctionService struct {
	repository AuctRepo
	logger     *zap.Logger
}

// NewAuctionService creates a new AuthService
func NewAuctionService(repo AuctRepo, logger *zap.Logger) *AuctionService {
	return &AuctionService{repository: repo, logger: logger}
}

// Search registers a new user
func (service *AuctionService) Search(request model.AuctionRequest) ([]model.Auction, error) {
	result, err := service.repository.FindAll(request)

	if err != nil {
		return nil, err
	}

	return result, nil
}

// Find - finds auction by uuid
func (service *AuctionService) Find(uuid string) (model.Auction, error) {

	auction, err := service.repository.FindOne(uuid)
	if err != nil {
		service.logger.Error("failed get auction", zap.Error(err))
		return model.Auction{}, err
	}

	return auction, nil
}

// Delete - deletes an auction by uuid
func (service *AuctionService) Delete(uuid string) error {

	err := service.repository.Delete(uuid)
	if err != nil {
		service.logger.Error("AuctionService.Delete failed deleting auction", zap.Error(err))
		return err
	}

	return nil
}

func (service *AuctionService) Do(param interface{}) error {
	service.logger.Info("do nothing...")

	return nil
}
func (service *AuctionService) Health() utils.ServiceHealthCheck {
	var healthCheck utils.ServiceHealthCheck

	healthCheck.ServiceStatus = "UP"
	dbStats := service.repository.DbStatus()
	healthCheck.DBStatus = append(healthCheck.DBStatus, dbStats)

	return healthCheck
}
