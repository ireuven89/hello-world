package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type MigrationTask struct {
	ID                   primitive.ObjectID `bson:"_id"`
	Name                 string             `bson:"name"`
	Status               string             `bson:"status"`
	ErrorMessage         string             `bson:"errorMessage"`
	Class                string             `bson:"_class"`
	HttpEndpoint         string             `bson:"httpExecute"`
	HttpBody             interface{}        `bson:"httpBody"`
	HttpParams           []string           `bson:"httpParams"`
	HttpMethod           string             `bson:"httpMethod"`
	HttpRollBack         string             `bson:"httpRollBack"`
	HttpRollBackMethod   string             `bson:"httpRollBackMethod"`
	RollbackErrorMessage string             `bson:"rollbackErrorMessage"`
	Params               interface{}        `bson:"-"`
	RollbackParams       interface{}        `bson:"-"`
}
