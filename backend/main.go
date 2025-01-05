package main

import (
	"github.com/ireuven89/hello-world/routes"
	"github.com/labstack/echo"
)

func main() {
	e := echo.New()

	routes.AssignRoutes(e)

	e.Logger.Fatal("failed to initiate server", e.Start(":7000"))
}
