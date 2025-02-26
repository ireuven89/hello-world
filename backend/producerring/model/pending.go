package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MigrationStatus int
type TaskStatus int

const (
	Pending MigrationStatus = iota
	InProgress
	Stopped
	Finished
)

const (
	PendingTask TaskStatus = iota
	InProgressTask
	Completed
	Failed
)

func (ms MigrationStatus) String() string {
	return []string{"Pending", "InProgress", "Stopped", "Finished"}[ms]
}

func (ts TaskStatus) String() string {
	return []string{"Pending", "InProgress", "Stopped", "Finished"}[ts]
}

type Migration struct {
	ID                 primitive.ObjectID `bson:"_id"`
	Type               string             `bson:"type"`
	MigrationName      string             `bson:"name"`
	QueueName          string             `bson:"queueName"`
	NumOfThreads       int                `bson:"numOfThreads"`
	HttpEndpoint       string             `bson:"httpExecute"`
	HttpMethod         string             `bson:"httpMethod"`
	HttpRollBack       string             `bson:"httpRollBack"`
	HttpRollBackMethod string             `bson:"httpRollBackMethod"`
	Status             string             `bson:"status"`
	CreatedAt          time.Time          `bson:"createdAt"`
}

type MigrationTask struct {
	ID                   primitive.ObjectID `bson:"_id"`
	Name                 string             `bson:"name"`
	Status               string             `bson:"status"`
	ErrorMessage         string             `bson:"errorMessage"`
	Class                string             `bson:"_class"`
	HttpBody             interface{}        `bson:"httpBody"`
	HttpParams           []string           `bson:"httpParams"`
	RollbackErrorMessage string             `bson:"rollbackErrorMessage"`
	Params               interface{}        `bson:"-"`
	RollbackParams       interface{}        `bson:"-"`
}
