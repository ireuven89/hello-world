package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

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

type MigrationType int

const (
	Http MigrationType = iota
	Internal
)

func (t MigrationType) String() string {

	return []string{"http", "internal"}[t]
}

type MigrationStatus int

const (
	Pending MigrationStatus = iota
	InProgress
	Stopped
	Finished
)

type TaskStatus int

const (
	TaskPending TaskStatus = iota
	TaskInProgress
	TaskCompleted
	TaskFailed
)

func (ms MigrationStatus) String() string {
	return []string{"pending", "in_progress", "stopped", "finished"}[ms]
}

func (ts TaskStatus) String() string {
	return []string{"pending", "in_progress", "completed", "failed"}[ts]
}
