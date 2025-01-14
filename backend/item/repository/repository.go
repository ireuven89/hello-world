package repository

import (
	"fmt"
	"github.com/google/uuid"
	"time"

	"github.com/ido50/sqlz"
	"github.com/ireuven89/hello-world/backend/item/model"
	"github.com/ireuven89/hello-world/backend/redis"
	"go.uber.org/zap"
)

type ItemRepository struct {
	db     *sqlz.DB
	redis  *redis.Service
	logger *zap.Logger
}

const redisTtl = time.Minute * 3

func New(db *sqlz.DB, logger *zap.Logger, redis *redis.Service) *ItemRepository {

	return &ItemRepository{
		db:     db,
		logger: logger,
		redis:  redis,
	}
}

func (r *ItemRepository) ListItems(input model.ListInput) ([]model.Item, error) {
	var where []sqlz.WhereCondition

	queryString := fmt.Sprintf("name=%s&link=%s&description=%s", input.Name, input.Link, input.Description)

	result, err := r.redis.Get(queryString)

	if _, ok := result.([]model.Item); ok {
		return result.([]model.Item), nil
	}

	q := r.db.
		Select("id", "uuid", "mame", "link", "userUuid", "category").
		From("items")

	if input.Name != "" {
		where = append(where, sqlz.Like("name", input.Name))
	}

	if input.Link != "" {
		where = append(where, sqlz.Like("name", input.Link))
	}

	if input.Description != "" {
		where = append(where, sqlz.Like("name", input.Description))
	}

	q.Where(where...)

	if err = q.GetAll(&result); err != nil {
		return nil, err
	}

	if err = r.redis.Set(queryString, result, redisTtl); err != nil {
		r.logger.Warn("failed to set redis key: ", zap.Any("error", err))
	}

	if _, ok := result.([]model.Item); ok {
		return result.([]model.Item), nil
	}

	return []model.Item{}, nil
}

func (r *ItemRepository) GetItem(uuid string) (model.Item, error) {

	queryString := fmt.Sprintf("uuid=%s", uuid)

	result, err := r.redis.Get(queryString)

	if err != nil {
		return result.(model.Item), nil
	}

	q := r.db.
		Select("id", "uuid", "mame", "link", "userUuid", "category").
		From("items").
		Where(sqlz.Eq("uuid", uuid))

	if err = q.GetRow(&result); err != nil {
		r.logger.Error("failed to set redis key: ", zap.Any("error", err))
		return model.Item{}, err
	}

	if err = r.redis.Set(queryString, result, redisTtl); err != nil {
		r.logger.Warn("failed to set redis key: ", zap.Any("error", err))
	}

	if _, ok := result.(model.Item); ok {
		return result.(model.Item), nil
	}

	return model.Item{}, nil
}

func (r *ItemRepository) Upsert(item model.ItemInput) (string, error) {
	var create bool

	if create {
		id := uuid.New().String()
		q := r.db.
			InsertInto("items").
			ValueMap(map[string]interface{}{
				"id":        id,
				"name":      item.Name,
				"user_uuid": item.UserUuid,
				"category":  item.Category,
			}).Returning("id")

		err := q.GetRow(&id)

		if err != nil {
			return "", err
		}

		return id, err
	} else {
		var id string

		q := r.db.
			Update("items").SetMap(
			map[string]interface{}{
				"user_uuid": item.UserUuid,
				"name":      item.Name,
				"category":  item.Category,
			}).
			Where(sqlz.Eq("id", item.Uuid)).
			Returning("id")

		if err := q.GetRow(&id); err != nil {
			return "", err
		}

		return id, nil
	}
}

func (r *ItemRepository) Delete(uuid string) error {

	queryString := fmt.Sprintf("uuid=%s", uuid)

	result, err := r.redis.Get(queryString)

	if err != nil {
		return nil
	}

	q := r.db.DeleteFrom("items").
		Where(sqlz.Eq("uuid", uuid))

	if err = q.GetRow(&result); err != nil {
		r.logger.Error("failed deleting form db: ", zap.Any("error", err))
		return err
	}

	return nil
}
