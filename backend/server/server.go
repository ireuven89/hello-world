package server

import (
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/julienschmidt/httprouter"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/ireuven89/hello-world/backend/authenticating"
	authrepo "github.com/ireuven89/hello-world/backend/authenticating/repository"
	"github.com/ireuven89/hello-world/backend/aws"
	"github.com/ireuven89/hello-world/backend/biddering"
	"github.com/ireuven89/hello-world/backend/db"
	"github.com/ireuven89/hello-world/backend/elastic"
	"github.com/ireuven89/hello-world/backend/environment"
	"github.com/ireuven89/hello-world/backend/itemming"
	itemrepo "github.com/ireuven89/hello-world/backend/itemming/repository"
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
	Redis       redis.Redis
	Auth        authenticating.Service
}

func New() (*Server, error) {
	err := environment.Load()

	loggerConfig := zap.NewDevelopmentConfig()
	loggerConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	logger, err := loggerConfig.Build()

	config, err := utils.LoadConfig("authenticating", os.Getenv("ENV"))

	if err != nil {
		logger.Error("failed loading config")
		return nil, err
	}

	//base services
	redisClient, err := redis.New(logger)
	if err != nil {
		return nil, err
	}
	config, err = utils.LoadConfig("elastic", os.Getenv("ENV"))
	es, err := elastic.New(logger, config.Databases["elastic"])
	if err != nil {
		return nil, err
	}

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
	go authTransport.ListenAndServe(config.ServicePort)

	//itemming
	itemConfig, err := utils.LoadConfig("itemming", os.Getenv("ENV"))
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
	transport := userring.NewTransport(usersService, userRouter)
	go transport.ListenAndServe("7000")

	//publishing
	publiserr, err := publishing.New(logger)

	if err != nil {
		return nil, err
	}

	//subscribing
	subscriberr, err := subscribing.New(logger)
	if err != nil {
		return nil, err
	}

	//biddering
	bidderConfig, err := utils.LoadConfig("biddering", os.Getenv("ENV"))
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
	go bidderTransport.ListenAndServe(bidderConfig.ServicePort)

	//recovering

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

	return &Server{Auth: authService, Redis: redisClient, ItemService: itemService, UserService: usersService, Logger: logger, Echo: echoServer, AWSClient: awsClient, Elastic: es, Sub: subscriberr, Pub: publiserr}, nil
}
