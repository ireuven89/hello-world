package server

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/ireuven89/hello-world/backend/aws"
	"github.com/ireuven89/hello-world/backend/db"
	"github.com/ireuven89/hello-world/backend/elastic"
	"github.com/ireuven89/hello-world/backend/environment"
	"github.com/ireuven89/hello-world/backend/item"
	itemrepo "github.com/ireuven89/hello-world/backend/item/repository"
	"github.com/ireuven89/hello-world/backend/rabbit"
	"github.com/ireuven89/hello-world/backend/redis"
	"github.com/ireuven89/hello-world/backend/routes"
	"github.com/ireuven89/hello-world/backend/users"
	userrepo "github.com/ireuven89/hello-world/backend/users/repository"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Server struct {
	UserService  users.Service
	ItemService  item.Service
	Logger       *zap.Logger
	Echo         *echo.Echo
	Elastic      *elastic.Service
	AWSClient    aws.Service
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

	//items service
	itemsDB, itemsMigrationDir, err := item.MustNewDB()

	if err != nil {
		logger.Error(fmt.Sprintf("failed to initiate db %v", err))
		return nil, err
	}

	itemsMigration := db.New(itemsDB.DB.DB, logger, itemsMigrationDir)
	if err = itemsMigration.Run(); err != nil {
		return nil, err
	}
	itemRepo := itemrepo.New(itemsDB, logger, redisClient)
	itemService := item.New(itemRepo, logger)

	//users service
	usersDB, userMigrationDir, err := users.MustNewDB()
	if err != nil {
		logger.Error(fmt.Sprintf("failed to initiate users DB %v", err))
		panic(err)
	}

	userMigration := db.New(usersDB.DB.DB, logger, userMigrationDir)
	if err = userMigration.Run(); err != nil {
		panic(err)
	}

	userRepo := userrepo.New(usersDB, redisClient, logger)
	usersService := users.New(logger, userRepo)

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

	return &Server{Redis: redisClient, ItemService: itemService, UserService: usersService, Logger: logger, Echo: echoServer, AWSClient: awsclient, Elastic: es, RabbitClient: rabbitClient}, nil
}
