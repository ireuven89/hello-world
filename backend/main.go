package main

import (
	"github.com/ireuven89/hello-world/backend/server"

	"log"
)

func main() {

	mainServer, err := server.New()

	if err != nil {
		panic(err)
	}

	log.Fatal("failed to initiate server", mainServer.Echo.Start(":7000"))
}
