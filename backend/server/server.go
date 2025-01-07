package server

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/ireuven89/hello-world/backend/aws"
	"github.com/ireuven89/hello-world/backend/db"
	"github.com/ireuven89/hello-world/backend/elastic"
	"github.com/ireuven89/hello-world/backend/environment"
	"github.com/ireuven89/hello-world/backend/item"
	"github.com/ireuven89/hello-world/backend/rabbit"
	"github.com/ireuven89/hello-world/backend/redis"
	"github.com/ireuven89/hello-world/backend/routes"
	"github.com/ireuven89/hello-world/backend/user"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Server struct {
	UserDB       *sql.DB
	ItemDB       *sql.DB
	Logger       *zap.Logger
	Echo         *echo.Echo
	Elastic      *elastic.Service
	AWSClient    *aws.Client
	RabbitClient *rabbit.Client
	Redis        *redis.Service
}

func New() (*Server, error) {
	err := environment.Load()

	loggerConfig := zap.NewDevelopmentConfig()
	loggerConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	logger, err := loggerConfig.Build()

	if err != nil {
		return nil, err
	}

	redisClient, err := redis.New(logger)

	if err != nil {
		return nil, err
	}

	usersDB, userMigrationDir, err := user.MustNewDB()
	if err != nil {
		logger.Error(fmt.Sprintf("failed to initiate user DB %v", err))
		return nil, err
	}

	userMigration := db.New(usersDB, logger, userMigrationDir)
	if err = userMigration.Run(); err != nil {
		return nil, err
	}

	itemsDB, itemsMigrationDir, err := item.MustNewDB()

	if err != nil {
		logger.Error(fmt.Sprintf("failed to initiate db %v", err))
		return nil, err
	}

	//migrate
	itemsMigration := db.New(itemsDB, logger, itemsMigrationDir)
	if err = itemsMigration.Run(); err != nil {
		return nil, err
	}

	es, err := elastic.New()
	if err != nil {
		return nil, err
	}
	rabbitClient, err := rabbit.New(logger)

	awsclient, err := aws.New(logger)

	if err != nil {
		return nil, err
	}

	echoServer := echo.New()
	if err != nil {
		return nil, err
	}
	logger.Debug("server initiated DB")

	if err != nil {
		return nil, err
	}

	routes.AssignRoutes(echoServer)

	logger.Info("Server has been initialized")

	return &Server{Redis: redisClient, UserDB: usersDB, ItemDB: itemsDB, Logger: logger, Echo: echoServer, AWSClient: awsclient, Elastic: es, RabbitClient: rabbitClient}, nil
}
