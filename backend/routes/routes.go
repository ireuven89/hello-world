package routes

import (
	"github.com/labstack/echo"
	"net/http"
)

func AssignRoutes(e *echo.Echo) {

	e.GET("/health", HealthHandler)

	group := e.Group("api/v1")
	group.Add(http.MethodGet, "/db", GetDbbHandler)
}

func GetDbbHandler(c echo.Context) error {

	return nil
}

func HealthHandler(c echo.Context) error {
	var health HealthResponse
	health.Status = "Ok"
	health.Version = 1.0

	return c.JSON(200, health)
}
