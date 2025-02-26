package main

import (
	"fmt"

	"github.com/labstack/gommon/log"

	"github.com/ireuven89/hello-world/backend/server"
	"github.com/ireuven89/hello-world/backend/utils"
)

func main() {

	secretes := utils.Secrets{
		Name: "root",
		Secrete: []utils.Secrets{
			{
				Name: "path",
				Secrete: []utils.Secrets{{
					Name:    "secret",
					Secrete: nil,
				}},
			}, {
				Name: "different-path",
				Secrete: []utils.Secrets{{
					Name:    "secret1",
					Secrete: nil,
				}},
			},
		},
	}

	res := utils.ListSecrets(secretes, "")

	for index, s := range res {
		fmt.Printf("path %v is %v\n", index, s)
	}

	mainServer, err := server.New()

	if err != nil {
		panic(err)
	}

	/*	tenant := os.Getenv("TENANT")
		configPath, err := filepath.abs(fmt.Sprintf("./userring/config/%s", tenant))
		config, err := utils.GetConfiguration(configPath)*/

	/*loggerConfig := zap.NewDevelopmentConfig()
	loggerConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	logger, err := loggerConfig.Build()*/

	//redisClient, err := redis.New(logger)

	/*	if err != nil {
			panic(err)
		}
	*/
	//userTransport := userring.NewTransport(httprouter.New(), usersService)
	//userring.RegisterRoutes(httprouter.New(), usersService)

	//http.ListenAndServe(config.TenantEndpoint)

	log.Fatal("failed to initiate server", mainServer.Echo.Start(":7000"))
}
