package repository

import (
	"time"

	"github.com/google/uuid"
	"github.com/ido50/sqlz"
	"go.uber.org/zap"

	"github.com/ireuven89/hello-world/backend/authenticating/model"
	"github.com/ireuven89/hello-world/backend/db/utils"
)

type Redis interface {
	Get(key string) (interface{}, error)
	Set(key string, value interface{}, ttl time.Duration) error
}

type Repo struct {
	logger *zap.Logger
	db     *sqlz.DB
}

func New(logger *zap.Logger, db *sqlz.DB) *Repo {

	return &Repo{
		db:     db,
		logger: logger,
	}
}

func (r *Repo) Save(username, password string) error {
	id := uuid.New().String()

	q := r.db.InsertInto(model.TableName).
		ValueMap(map[string]interface{}{
			"id":       id,
			"user":     username,
			"password": password,
		})

	utils.New().DebugInsert(q, "create user")

	if _, err := q.Exec(); err != nil {
		r.logger.Error("failed to insert user", zap.Error(err))
		return err
	}

	return nil
}

func (r *Repo) Find(username string) (model.User, error) {
	var result model.User

	q := r.db.Select("user", "password").
		From("users").
		Where(sqlz.WhereCondition(sqlz.Eq("user", username)))

	if err := q.GetRow(&result); err != nil {
		r.logger.Error("failed to insert user")
		return model.User{}, err
	}

	utils.New().DebugSelect(q, "create user")

	return model.User{
		Username: result.Username,
		Password: result.Password,
	}, nil
}
