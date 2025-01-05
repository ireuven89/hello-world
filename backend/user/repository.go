package user

import (
	"github.com/ido50/sqlz"
	dbmodel "github.com/ireuven89/hello-world/backend/db/model"
	"github.com/ireuven89/hello-world/backend/db/utils"
	"github.com/ireuven89/hello-world/backend/user/model"
)

type User interface {
	Single(request model.UserFetchInput) (model.User, error)
	List(request model.UserFetchInput) ([]model.User, error)
	Upsert(request model.UserUpsertInput) (model.User, error)
	Delete(request model.DeleteUserInput) error
}

type Repository struct {
	db *sqlz.DB
}

func New(db *sqlz.DB) *Repository {

	return &Repository{
		db: db,
	}
}

// List - this method queries users from DB
func (r *Repository) List(request model.UserFetchInput) ([]model.User, error) {
	var result []model.User

	q := r.db.Select(
		"id",
		"uuid",
		"name",
		"region").
		From(dbmodel.Users)

	var whereClauses []sqlz.WhereCondition

	if request.Name != "" {
		whereClauses = append(whereClauses, sqlz.Eq("name", request.Name))
	}

	if request.Region != "" {
		whereClauses = append(whereClauses, sqlz.Eq("region", request.Region))

	}

	q.Where(whereClauses...)

	utils.New().DebugSelect(q, "fetch users")

	err := q.GetAll(&result)

	if err != nil {
		return result, err
	}

	return result, nil
}

// Single - this method queries single user from DB
func (r *Repository) Single(input model.UserFetchInput) (model.User, error) {
	var result model.User

	q := r.db.Select(
		"id",
		"uuid",
		"name",
		"region").
		From(dbmodel.Users).
		Where(sqlz.Eq("uuid", input.Uuid))

	utils.New().DebugSelect(q, "get user")

	err := q.GetRow(&result)

	if err != nil {
		return result, err
	}

	return result, nil
}

// Upsert - this method upsert user to DB
func (r *Repository) Upsert(input model.UserUpsertInput) (model.User, error) {
	var result model.User
	var create bool

	if input.Uuid == "" {
		create = true
	}

	if create {
		q := r.db.InsertInto(dbmodel.Users).ValueMap(map[string]interface{}{
			"uuid":   input.Uuid,
			"name":   input.Name,
			"region": input.Region,
		}).
			Returning("id",
				"uuid",
				"name",
				"region")

		utils.New().DebugInsert(q, "insert user")

		err := q.GetRow(&result)

		if err != nil {
			return result, err
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
				"uuid",
				"name",
				"region",
			)

		utils.New().DebugUpdate(q, "update user")

		err := q.GetRow(&result)

		if err != nil {
			return result, err
		}
	}

	return result, nil
}

// Delete - this query deletes user from DB
func (r *Repository) Delete(input model.DeleteUserInput) error {

	q := r.db.DeleteFrom("users").Where(sqlz.Eq("uuid", input.Uuid))

	utils.New().DebugDelete(q, "delete user")

	_, err := q.Exec()

	if err != nil {
		return err
	}

	return nil
}
