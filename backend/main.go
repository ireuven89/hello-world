package main

import (
	"github.com/ireuven89/hello-world/backend/server"
	"github.com/labstack/gommon/log"
)

func main() {

	mainServer, err := server.New()

	/*	tenant := os.Getenv("TENANT")
		configPath, err := filepath.abs(fmt.Sprintf("./users/config/%s", tenant))
		config, err := utils.GetConfiguration(configPath)*/

	/*loggerConfig := zap.NewDevelopmentConfig()
	loggerConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	logger, err := loggerConfig.Build()*/

	//redisClient, err := redis.New(logger)

	if err != nil {
		panic(err)
	}

	//userTransport := users.NewTransport(httprouter.New(), usersService)
	//users.RegisterRoutes(httprouter.New(), usersService)

	//http.ListenAndServe(config.TenantEndpoint)

	log.Fatal("failed to initiate server", mainServer.Echo.Start(":7000"))
}
