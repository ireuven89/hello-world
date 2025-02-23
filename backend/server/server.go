package server

import (
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/julienschmidt/httprouter"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/ireuven89/hello-world/backend/itemming"
	itemrepo "github.com/ireuven89/hello-world/backend/itemming/repository"

	"github.com/ireuven89/hello-world/backend/authenticating"
	authrepo "github.com/ireuven89/hello-world/backend/authenticating/repository"
	"github.com/ireuven89/hello-world/backend/aws"
	"github.com/ireuven89/hello-world/backend/db"
	"github.com/ireuven89/hello-world/backend/elastic"
	"github.com/ireuven89/hello-world/backend/environment"
	"github.com/ireuven89/hello-world/backend/publishing"
	"github.com/ireuven89/hello-world/backend/redis"
	"github.com/ireuven89/hello-world/backend/routes"
	"github.com/ireuven89/hello-world/backend/subscribing"
	"github.com/ireuven89/hello-world/backend/userring"
	userrepo "github.com/ireuven89/hello-world/backend/userring/repository"
	"github.com/ireuven89/hello-world/backend/utils"
)

type Server struct {
	UserService userring.Service
	ItemService itemming.Service
	Logger      *zap.Logger
	Echo        *echo.Echo
	Elastic     elastic.Service
	AWSClient   aws.Service
	Pub         publishing.PService
	Sub         subscribing.SService
	Redis       *redis.Service
	Auth        authenticating.Service
}

func New() (*Server, error) {
	err := environment.Load()

	loggerConfig := zap.NewDevelopmentConfig()
	loggerConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	logger, err := loggerConfig.Build()

	config, err := utils.LoadConfig("authenticating", os.Getenv("env"))

	//base services
	redisClient, err := redis.New(logger)
	es, err := elastic.New(logger)
	if err != nil {
		return nil, err
	}
	/*
		mongoCLient, err := migration2.MustNewDB(config.Databases["mongo"])
		migrationDB := mongoCLient.Database("migrations")
		queues := mongoCLient.Database("queues")

		migrationService := migration2.NewService(logger, migrationDB, queues)
		migrationService.ProcessTasks(context.Background(), "test-migration")
	*/
	awsClient, err := aws.New(logger)

	if err != nil {
		return nil, err
	}
	//authenticating
	authDB, dir, err := authenticating.MustNewDB(config.Databases["mysql"])
	if err != nil {
		return nil, err
	}

	authMigration := db.New(authDB, logger, dir)
	if err = authMigration.Run(); err != nil {
		return nil, err
	}

	userStore := authrepo.New(logger, authDB)
	authService := authenticating.NewAuthService(userStore, logger)
	authRouter := httprouter.New()
	authTransport := authenticating.NewTransport(authService, authRouter)
	logger.Info("starting auth service...")
	go authTransport.ListenAndServe(config.ServicePort)

	//itemming
	itemConfig, err := utils.LoadConfig("itemming", os.Getenv("env"))
	if err != nil {
		logger.Error(fmt.Sprintf("failed to initiate db %v", err))
		return nil, err
	}
	itemsDB, itemsMigrationDir, err := itemming.MustNewDB()

	if err != nil {
		logger.Error(fmt.Sprintf("failed to initiate db %v", err))
		return nil, err
	}
	itemsMigration := db.New(itemsDB, logger, itemsMigrationDir)
	if err = itemsMigration.Run(); err != nil {
		return nil, err
	}
	itemRepo := itemrepo.New(itemsDB, logger, redisClient)
	itemService := itemming.New(itemRepo, logger)
	itemRouter := httprouter.New()
	itemTransport := itemming.NewTransport(itemService, itemRouter)
	logger.Info("starting auth service...")
	go itemTransport.ListenAndServe(itemConfig.ServicePort)

	//userring
	usersDB, userMigrationDir, err := userring.MustNewDB()
	if err != nil {
		logger.Error(fmt.Sprintf("failed to initiate userring DB %v", err))
		panic(err)
	}

	userMigration := db.New(usersDB, logger, userMigrationDir)
	if err = userMigration.Run(); err != nil {
		panic(err)
	}

	userRepo := userrepo.New(usersDB, redisClient, logger)
	usersService := userring.New(logger, userRepo)
	userRouter := httprouter.New()
	userTransport := userring.NewTransport(usersService, userRouter)
	go userTransport.ListenAndServe(config.ServicePort)
	logger.Info("starting auth service...")
	//biddering
	/*bidderConfig, err := utils.LoadConfig("biddering", os.Getenv("ENV"))
	if err != nil {
		return nil, err
	}
	bidderDb, dir, err := biddering.MustNewDB(bidderConfig.Databases["mysql"])
	if err != nil {
		return nil, err
	}
	migration := db.New(bidderDb, logger, dir)
	if err = migration.Run(); err != nil {
		return nil, err
	}
	bidderRepo := biddering.NewRepository(bidderDb, logger, redisClient)
	bidderService := biddering.NewService(bidderRepo, logger)
	bidderRoute := httprouter.New()
	bidderTransport := biddering.NewTransport(bidderService, bidderRoute)
	logger.Info("starting auth service...")
	go bidderTransport.ListenAndServe(bidderConfig.ServicePort)*/

	//remoting
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

	return &Server{Auth: authService, Redis: redisClient, ItemService: itemService, UserService: usersService, Logger: logger, Echo: echoServer, AWSClient: awsClient, Elastic: es}, nil
}
