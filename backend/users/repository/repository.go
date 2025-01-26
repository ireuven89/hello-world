package users

import (
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/ido50/sqlz"
	"go.uber.org/zap"

	dbmodel "github.com/ireuven89/hello-world/backend/db/model"
	"github.com/ireuven89/hello-world/backend/db/utils"
	"github.com/ireuven89/hello-world/backend/users/model"
)

type Redis interface {
	Get(key string) (interface{}, error)
	Set(key string, value interface{}, ttl time.Duration) error
}

type UserRepository struct {
	db     *sqlz.DB
	redis  Redis
	logger *zap.Logger
}

func New(db *sqlz.DB, redis Redis, logger *zap.Logger) *UserRepository {

	return &UserRepository{
		db:     db,
		redis:  redis,
		logger: logger,
	}
}

const redisQueryTTl = time.Minute * 3

// ListUsers - this method queries users from DB
func (r *UserRepository) ListUsers(input model.UserFetchInput) ([]model.User, error) {
	var result []model.User

	cachedQuery := fmt.Sprintf("ListUsers:%s%s%s%v%v", input.Region, input.Name, input.Uuid, input.Page, input.Size)

	//get result from redis
	cachedResult, err := r.redis.Get(cachedQuery)
	if err == nil {
		return cachedResult.([]model.User), nil
	}

	q := r.db.Select(
		"id",
		"uuid",
		"name",
		"region").
		From(dbmodel.Users)

	var whereClauses []sqlz.WhereCondition

	if input.Name != "" {
		whereClauses = append(whereClauses, sqlz.Eq("name", input.Name))
	}

	if input.Region != "" {
		whereClauses = append(whereClauses, sqlz.Eq("region", input.Region))

	}

	q.Where(whereClauses...)

	utils.New().DebugSelect(q, "fetch users")

	err = q.GetAll(&result)

	if err != nil {
		return result, err
	}

	//cache the result
	if err = r.redis.Set(cachedQuery, result, redisQueryTTl); err != nil {
		r.logger.Warn("Failed to cache query: ", zap.Error(err))
	}

	return result, nil
}

// FindUser - this method queries single users from DB
func (r *UserRepository) FindUser(uuid string) (model.User, error) {
	var result model.User

	cachedQuery := fmt.Sprintf("FindUser:%s", uuid)

	//get result from redis
	cachedResult, err := r.redis.Get(cachedQuery)
	if err == nil {
		return cachedResult.(model.User), nil
	}

	q := r.db.Select(
		"id",
		"uuid",
		"name",
		"region").
		From(dbmodel.Users).
		Where(sqlz.Eq("uuid", uuid))

	utils.New().DebugSelect(q, "get users")

	err = q.GetRow(&result)

	if err != nil {
		return result, err
	}

	//cache the result
	if err = r.redis.Set(cachedQuery, result, redisQueryTTl); err != nil {
		r.logger.Warn("Failed to cache query: ", zap.Error(err))
	}

	return result, nil
}

// Upsert - this method upsert users to DB
func (r *UserRepository) Upsert(input model.UserUpsertInput) (string, error) {
	var id string
	var create bool

	if input.Uuid == "" {
		create = true
	}

	if create {
		id = uuid.New().String()
		q := r.db.InsertInto(dbmodel.Users).ValueMap(map[string]interface{}{
			"id":     id,
			"name":   input.Name,
			"region": input.Region,
		}).
			Returning("id")

		utils.New().DebugInsert(q, "insert users")

		if err := q.GetRow(&id); err != nil {
			r.logger.Error("failed to insert model: ", zap.Error(err))
			return id, err
		}
	} else {
		q := r.db.Update(dbmodel.Users).
			SetMap(map[string]interface{}{
				"name":   input.Name,
				"region": input.Region,
			}).
			Where(sqlz.Eq("uuid", input.Uuid)).
			Returning(
				"id",
			)

		utils.New().DebugUpdate(q, "update users")

		if err := q.GetRow(&id); err != nil {
			r.logger.Error("failed to insert model: ", zap.Error(err))
			return id, err
		}
	}

	return id, nil
}

// Delete - this query deletes users from DB
func (r *UserRepository) Delete(uuid string) error {

	q := r.db.DeleteFrom("users").Where(sqlz.Eq("uuid", uuid))

	utils.New().DebugDelete(q, "delete users")

	if _, err := q.Exec(); err != nil {
		return err
	}

	return nil
}
