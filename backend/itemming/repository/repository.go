package repository

import (
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/ido50/sqlz"
	"go.uber.org/zap"

	"github.com/ireuven89/hello-world/backend/itemming/model"
)

type ItemRepository struct {
	db     *sqlz.DB
	redis  Redis
	logger *zap.Logger
}

type Redis interface {
	Get(key string) (interface{}, error)
	Set(key string, value interface{}, ttl time.Duration) error
}

const redisTtl = time.Minute * 3

func New(db *sqlz.DB, logger *zap.Logger, redis Redis) *ItemRepository {

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
		Select("id", "uuid", "name", "link", "userUuid", "category").
		From("items")

	if input.Name != "" {
		where = append(where, sqlz.Like("name", input.Name))
	}

	if input.Link != "" {
		where = append(where, sqlz.Like("link", input.Link))
	}

	if input.Description != "" {
		where = append(where, sqlz.Like("description", input.Description))
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
	create = item.Uuid == ""

	if create {
		id := uuid.New().String()
		q := r.db.
			InsertInto("items").
			ValueMap(map[string]interface{}{
				"id":        id,
				"name":      item.Name,
				"user_uuid": item.UserUuid,
				"category":  item.Category,
			})

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
			Where(sqlz.Eq("id", item.Uuid))

		if err := q.GetRow(&id); err != nil {
			return "", err
		}

		return id, nil
	}
}

func (r *ItemRepository) Delete(uuid string) error {
	q := r.db.DeleteFrom("items").
		Where(sqlz.Eq("uuid", uuid))

	if _, err := q.Exec(); err != nil {
		r.logger.Error("Delete failed deleting form db: ", zap.Any("error", err))
		return err
	}

	return nil
}
