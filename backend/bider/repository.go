package bider

import (
	"fmt"
	"time"

	"github.com/ido50/sqlz"
	"github.com/ireuven89/hello-world/backend/bider/model"
	dbmodel "github.com/ireuven89/hello-world/backend/db/model"
	"github.com/ireuven89/hello-world/backend/db/utils"
	"github.com/ireuven89/hello-world/backend/redis"
	"go.uber.org/zap"
)

type Bidder struct {
	ID        int64     `json:"-" sql:"id"`
	Uuid      string    `json:"uuid" sql:"uuid"`
	Name      string    `json:"name" sql:"name"`
	CreatedAt time.Time `json:"createdAt" sql:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" sql:"updated_at"`
}

type Repo interface {
	List(input model.BiddersInput) ([]model.Bidder, error)
	Single(uuid string) (model.Bidder, error)
	Upsert(input model.BiddersInput) (string, error)
	Delete(id string) error
}

type Repository struct {
	db     *sqlz.DB
	redis  *redis.Service
	logger *zap.Logger
}

func New(db *sqlz.DB, logger *zap.Logger, redis *redis.Service) *Repository {

	return &Repository{
		db:     db,
		logger: logger,
		redis:  redis,
	}
}

func (r *Repository) List(input model.BiddersInput) ([]model.Bidder, error) {
	var result []model.Bidder
	var where []sqlz.WhereCondition

	if input.Uuid != "" {
		where = append(where, sqlz.WhereCondition(sqlz.Eq("uuid", input.Uuid)))
	}
	if input.Name != "" {
		where = append(where, sqlz.WhereCondition(sqlz.Eq("uuid", input.Uuid)))
	}
	if input.Item != "" {
		where = append(where, sqlz.WhereCondition(sqlz.Eq("uuid", input.Uuid)))
	}
	if input.Uuid != "" {
		where = append(where, sqlz.WhereCondition(sqlz.Eq("uuid", input.Uuid)))
	}
	q := r.db.
		Select("uuid", "name", "created_at", "updated_at").From(dbmodel.Bidders)

	q.Where(where...)

	utils.New().DebugSelect(q, "select bidders")

	if err := q.GetAll(&result); err != nil {
		r.logger.Error(fmt.Sprintf("failed to get db %v", err))
		return nil, err
	}

	return result, nil
}

func (r *Repository) Single(uuid string) (model.Bidder, error) {
	var result model.Bidder

	q := r.db.Select("uuid", "name", "created_at", "updated_at").From(dbmodel.Bidders).
		Where(sqlz.WhereCondition(sqlz.Eq("uuid", uuid)))

	utils.New().DebugSelect(q, "single bidder")

	if err := q.GetRow(&result); err != nil {
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
			return "", err
		}
	} else {
		valuesMap := setValuesMap(input)
		q := r.db.
			Update(dbmodel.Bidders).
			SetMap(valuesMap).
			Where(sqlz.Eq("uuid", input.Uuid)).Returning("id")

		utils.New().DebugUpdate(q, "update bidder")

		if err := q.GetRow(&id); err != nil {
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
		r.logger.Error("failed to delete bidder", zap.String("uuid", uuid))
	}

	return nil
}
