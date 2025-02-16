package repoistory

import (
	"github.com/ido50/sqlz"
	"go.uber.org/zap"

	"github.com/ireuven89/hello-world/backend/auctioning/model"
	"github.com/ireuven89/hello-world/backend/utils"
)

type Repository struct {
	logger *zap.Logger
	db     *sqlz.DB
}

func New(logger *zap.Logger, db *sqlz.DB) *Repository {

	return &Repository{
		db:     db,
		logger: logger,
	}
}

// FindAll - Search all auctions by request
func (r *Repository) FindAll(req model.AuctionRequest) ([]model.Auction, error) {
	var result []model.Auction
	var where []sqlz.WhereCondition
	q := r.db.Select("id", "item", "price", "winning price", "bidders_price", "bidders_threshold").
		From("auctions")

	if req.Price != 0 {
		where = append(where, sqlz.Eq("price", req.Price))
	}

	if req.WinningPrice != 0 {
		where = append(where, sqlz.Eq("winning_price", req.WinningPrice))
	}

	if req.BiddersCount != 0 {
		where = append(where, sqlz.Eq("bidders_count", req.BiddersCount))
	}

	if req.BiddersThreshold != 0 {
		where = append(where, sqlz.Eq("bidders_threshold", req.BiddersThreshold))
	}

	if req.Status != "" {
		where = append(where, sqlz.Eq("status", req.Status))
	}

	if req.UserUuid != "" {
		where = append(where, sqlz.Eq("user_uuid", req.UserUuid))
	}

	q.Where(where...)

	if err := q.GetAll(&result); err != nil {
		r.logger.Error("Repository.FindAll failed fetching auctions", zap.Error(err))
		return nil, err
	}

	return result, nil
}

// FindOne - find auction by uuid
func (r *Repository) FindOne(uuid string) (model.Auction, error) {
	var result model.Auction
	q := r.db.Select("id", "item", "price", "winning price", "bidders_price", "bidders_threshold").
		From("auctions").
		Where(sqlz.Eq("id", uuid))

	if err := q.GetRow(&result); err != nil {
		r.logger.Error("AuctionRepository.FindOne failed fetching auction", zap.Error(err))
		return model.Auction{}, err
	}

	return model.Auction{}, nil
}

// Delete -
func (r *Repository) Delete(uuid string) error {
	q := r.db.DeleteFrom("auctions").Where(sqlz.Eq("id", uuid))

	if _, err := q.Exec(); err != nil {
		r.logger.Error("failed deleting auction", zap.Error(err))
		return err
	}

	return nil
}

// Update -
func (r *Repository) Update(req model.AuctionRequest) error {
	q := r.db.
		Update("auctions").
		SetMap(prepareSetMap(req)).
		Where(sqlz.Eq("id", req.Id))

	if _, err := q.Exec(); err != nil {
		r.logger.Error("failed deleting auction", zap.Error(err))
		return err
	}

	return nil
}

func prepareSetMap(req model.AuctionRequest) map[string]interface{} {
	var result map[string]interface{}

	if req.Status != "" {
		result["status"] = req.Status
	}

	if req.UserUuid != "" {
		result["user_uuid"] = req.UserUuid
	}

	if req.BiddersCount != 0 {
		result["bidders_count"] = req.BiddersCount
	}

	if req.BiddersThreshold != 0 {
		result["bidders_threshold"] = req.BiddersThreshold
	}

	return result
}

// DbStatus -
func (r *Repository) DbStatus() utils.DbStatus {

	return utils.DbStatus{}
}
