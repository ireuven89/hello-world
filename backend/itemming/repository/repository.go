package repository

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/ido50/sqlz"
	"go.uber.org/zap"

	"github.com/ireuven89/hello-world/backend/itemming/model"
	"github.com/ireuven89/hello-world/backend/utils"
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
	var result []model.Item

	queryString := fmt.Sprintf("name=%s&price=%s&description=%s", input.Name, input.Price, input.Description)

	cachedResult, err := r.redis.Get(queryString)

	if err == nil {
		r.logger.Debug("redis hit: ", zap.String("query", queryString))
		return cachedResult.([]model.Item), nil
	}

	q := r.db.
		Select("id", "uuid", "name", "link", "itemming", "price", "description").
		From("items")

	if input.Name != "" {
		where = append(where, sqlz.Like("name", input.Name))
	}

	if input.Price != "" {
		where = append(where, sqlz.Like("link", input.Price))
	}

	if input.Description != "" {
		where = append(where, sqlz.Like("description", input.Description))
	}

	q.Where(where...)

	if err = q.GetAll(&result); err != nil {
		r.logger.Error("ItemRepository.ListItems failed query items: ", zap.Any("error", err))
		return nil, err
	}

	//cache the result
	cachedResult, err = json.Marshal(result)
	if err != nil {
		r.logger.Warn("failed caching result", zap.Error(err))
	}
	if err = r.redis.Set(queryString, cachedResult, redisTtl); err != nil {
		r.logger.Warn("failed to set redis key: ", zap.Any("error", err))
	}

	return result, nil
}

func (r *ItemRepository) GetItem(uuid string) (model.Item, error) {

	queryString := fmt.Sprintf("uuid=%s", uuid)

	result, err := r.redis.Get(queryString)

	if err != nil {
		return result.(model.Item), nil
	}

	q := r.db.
		Select("id", "uuid", "mame", "link", "description", "price").
		From("items").
		Where(sqlz.Eq("uuid", uuid))

	if err = q.GetRow(&result); err != nil {
		r.logger.Error("ItemRepository.GetItem failed to get itemming: ", zap.Any("error", err))
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

// InsertBulk - will bulk inset
// retry mechanism on bulk insert
func (r *ItemRepository) BulkInsert(items []model.ItemInput) error {

	multipleMap := prepareBulkInsert(items)

	q := r.db.InsertInto("items").Columns("id", "name", "description", "price").ValueMultiple(multipleMap)

	tx, err := r.db.Begin()

	defer tx.Rollback()

	if err != nil {
		r.logger.Error(fmt.Sprintf("ItemRepository.BulkInsert failed transaction  %v", err))
		return err
	}

	statement, bindings := q.ToSQL(true)

	// exec
	if _, err = tx.Exec(statement, bindings...); err != nil {
		r.logger.Error(fmt.Sprintf("ItemRepository.BulkInsert failed inserting  %v", err))
		return err
	}

	if err = tx.Commit(); err != nil {
		r.logger.Error(fmt.Sprintf("ItemRepository.BulkInsert failed inserting - commit failure  %v", err))
		return err
	}

	return nil
}

func prepareBulkInsert(items []model.ItemInput) [][]interface{} {
	multipleMap := make([][]interface{}, 0, len(items))

	for _, item := range items {
		id := uuid.New().String()
		var input []interface{}
		input = append(input, id)
		input = append(input, item.Name)
		input = append(input, item.Description)
		input = append(input, item.Price)
		multipleMap = append(multipleMap, input)

	}

	return multipleMap
}

func (r *ItemRepository) Insert(item model.ItemInput) (string, error) {
	id := uuid.New().String()

	q := r.db.InsertInto("items").
		Columns("id", "name", "description", "price").ValueMap(map[string]interface{}{
		"id":          id,
		"name":        item.Name,
		"description": item.Description,
		"price":       item.Price,
	})

	if _, err := q.Exec(); err != nil {
		r.logger.Error(fmt.Sprintf("ItemRepository.Insert failed inserting %v", err))
		return "", err
	}

	return id, nil
}

func (r *ItemRepository) Update(item model.ItemInput) error {
	updateMap := getUpdateMap(item)

	q := r.db.Update("items").
		SetMap(updateMap).
		Where(sqlz.Eq("id", item.Uuid))

	if _, err := q.Exec(); err != nil {
		r.logger.Error(fmt.Sprintf("ItemRepository.Update failed inserting %v", err))
		return err
	}

	return nil
}

func getUpdateMap(item model.ItemInput) map[string]interface{} {
	var updateMap map[string]interface{}

	if item.Name != "" {
		updateMap["name"] = item.Name
	}

	if item.Description != "" {
		updateMap["name"] = item.Name
	}

	if item.Price != 0 {
		updateMap["name"] = item.Price
	}

	return updateMap
}

func (r *ItemRepository) Upsert(item model.ItemInput) (string, error) {
	var create bool
	create = item.Uuid == ""

	if create {
		var result model.Item
		id := uuid.New().String()
		q := r.db.
			InsertInto("items").
			ValueMap(map[string]interface{}{
				"id":          id,
				"name":        item.Name,
				"user_uuid":   item.UserUuid,
				"description": item.Description,
			},
			)

		_, err := q.Exec()

		if err != nil {
			r.logger.Error("ItemRepository.Upsert failed creating itemming")
			return "", err
		}

		//query the inserted for return value
		qSelect := r.db.Select("id").
			From("items").
			Where(sqlz.Eq("id", id))

		err = qSelect.GetRow(&result)

		if err != nil {
			return "", err
		}

		return result.ID, err
	} else {
		var id string

		q := r.db.
			Update("items").SetMap(
			map[string]interface{}{
				"user_uuid": item.UserUuid,
				"name":      item.Name,
				"category":  item.Description,
			}).
			Where(sqlz.Eq("id", item.Uuid))

		if err := q.GetRow(&id); err != nil {
			r.logger.Error("ItemRepository.Upsert failed updating itemming")
			return "", err
		}

		return id, nil
	}
}

func (r *ItemRepository) DBstatus() utils.DbStatus {
	var status utils.DbStatus

	status.Name = r.db.DriverName()
	if err := r.db.Ping(); err != nil {
		status.Connected = false
	} else {
		status.Connected = true
	}

	return status
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
