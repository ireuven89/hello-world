package bider

import (
	"fmt"
	"time"

	"github.com/ido50/sqlz"
	"go.uber.org/zap"

	"github.com/ireuven89/hello-world/backend/bider/model"
	dbmodel "github.com/ireuven89/hello-world/backend/db/model"
	"github.com/ireuven89/hello-world/backend/db/utils"
)

type Bidder struct {
	ID        int64     `json:"-" sql:"id"`
	Uuid      string    `json:"uuid" sql:"uuid"`
	Name      string    `json:"name" sql:"name"`
	CreatedAt time.Time `json:"createdAt" sql:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" sql:"updated_at"`
}

type Redis interface {
	Get(key string) (interface{}, error)
	Set(key string, value interface{}, ttl time.Duration) error
}

const redisQueryTtl = time.Minute * 3

type Repository struct {
	db     *sqlz.DB
	redis  Redis
	logger *zap.Logger
}

func New(db *sqlz.DB, logger *zap.Logger, redis Redis) *Repository {

	return &Repository{
		db:     db,
		logger: logger,
		redis:  redis,
	}
}

func (r *Repository) List(input model.BiddersInput) ([]model.Bidder, error) {
	var result []model.Bidder
	var where []sqlz.WhereCondition

	cachedQuery := fmt.Sprintf("%s%s%s%v%v", input.Uuid, input.Name, input.Item, input.Page.Offset, input.Page.GetLimit())

	cachedResult, err := r.redis.Get(cachedQuery)

	if err == nil {
		r.logger.Debug(fmt.Sprintf("List redis hit on input %s", cachedQuery))
		return cachedResult.([]model.Bidder), nil
	}

	if input.Name != "" {
		where = append(where, sqlz.WhereCondition(sqlz.Eq("name", input.Name)))
	}
	if input.Item != "" {
		where = append(where, sqlz.WhereCondition(sqlz.Eq("item", input.Item)))
	}

	q := r.db.
		Select("uuid", "name", "item", "created_at", "updated_at").From(dbmodel.Bidders).
		Offset(input.Page.Offset, input.Page.GetLimit())

	q.Where(where...)

	utils.New().DebugSelect(q, "select bidders")

	if err = q.GetAll(&result); err != nil {
		r.logger.Error(fmt.Sprintf("BidderRepo.List failed to get db %v", err))
		return nil, err
	}

	if err = r.redis.Set(cachedQuery, result, redisQueryTtl); err != nil {
		r.logger.Warn(fmt.Sprintf("failed to set redis q: %s %v", cachedQuery, err))
	}

	return result, nil
}

func (r *Repository) Single(uuid string) (model.Bidder, error) {
	var result model.Bidder

	q := r.db.Select("uuid", "name", "item", "created_at", "updated_at").From(dbmodel.Bidders).
		Where(sqlz.WhereCondition(sqlz.Eq("uuid", uuid)))

	utils.New().DebugSelect(q, "single bidder")

	if err := q.GetRow(&result); err != nil {
		r.logger.Error("BidderRepo.Single failed finding bidder", zap.Error(err))
		return result, err
	}

	return result, nil
}

func (r *Repository) Upsert(input model.BiddersInput) (string, error) {
	var create bool
	var id string
	if input.Uuid != "" {
		create = true
	}

	if create {
		q := r.db.InsertInto(dbmodel.Bidders).
			ValueMap(map[string]interface{}{
				"uuid":       input.Uuid,
				"item":       input.Item,
				"name":       input.Name,
				"created_at": time.Now(),
				"updated_at": time.Now(),
			}).Returning("id")

		utils.New().DebugInsert(q, "insert bidder")

		if err := q.GetRow(&id); err != nil {
			r.logger.Error("BidderRepo.Upsert failed creating bidder", zap.Error(err))
			return "", err
		}
	} else {
		valuesMap := setValuesMap(input)
		q := r.db.
			Update(dbmodel.Bidders).
			SetMap(valuesMap).
			Where(sqlz.Eq("uuid", input.Uuid))

		utils.New().DebugUpdate(q, "update bidder")

		if err := q.GetRow(&id); err != nil {
			r.logger.Error("BidderRepo.Upsert failed updating bidder", zap.Error(err))
			return "", err
		}
	}

	return id, nil
}

func setValuesMap(input model.BiddersInput) map[string]interface{} {
	var valuesMap map[string]interface{}

	if input.Item != "" {
		valuesMap["item"] = input.Item
	}

	if input.Name != "" {
		valuesMap["name"] = input.Name
	}

	return valuesMap
}

func (r *Repository) Delete(uuid string) error {

	q := r.db.DeleteFrom(dbmodel.Bidders).
		Where(sqlz.Eq("uuid", uuid))

	if _, err := q.Exec(); err != nil {
		r.logger.Error("BidderRepo.Delete failed deleting bidder", zap.Error(err))
	}

	return nil
}
