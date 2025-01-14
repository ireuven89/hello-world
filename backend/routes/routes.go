package routes

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// @title Swagger Example API
// @version 1.0
// @description This is a sample server Petstore server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host petstore.swagger.io
// @BasePath /v2
func AssignRoutes(e *echo.Echo) {

	e.GET("/health", HealthHandler)

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	group := e.Group("api/v1")

	//	e.Use(middlewareFunc)

	//handlers

	//user
	group.Add(http.MethodGet, "/users", GetUsersHandler)
	group.Add(http.MethodGet, "/users/:userUuid", GetUsersHandler)
	group.Add(http.MethodPost, "/users", UpsertUserHandler)
	group.Add(http.MethodPut, "/users/:userUuid", PutUserHandler)
	group.Add(http.MethodDelete, "/users/:userUuid", DeleteUserHandler)

	//auction
	group.Add(http.MethodGet, "/auctions", GetAuctionsHandler)
	group.Add(http.MethodGet, "/auctions/:userUuid/:auctionUuid", GetAuctionHandler)
	group.Add(http.MethodPost, "/auctions/:userUuid", PostAuctionHandler)
	group.Add(http.MethodPut, "/auctions/:userUuid", PutAuctionHandler)
	group.Add(http.MethodDelete, "/auctions/:userUuid", DeleteAuctionHandler)

	//item
	group.Add(http.MethodGet, "/items", GetItemsHandler)
	group.Add(http.MethodGet, "/items/:itemUuid", GetItemHandler)
	group.Add(http.MethodPost, "/items/:userUuid", PostItemHandler)
	group.Add(http.MethodPut, "/items/:userUuid/:itemUuid", PutItemHandler)
	group.Add(http.MethodDelete, "/items/:userUuid/:itemUuid", DeleteItemHandler)

	//elastic
	group.Add(http.MethodGet, "/elastic/:index/:id", GetIndexElasticHandler)
	group.Add(http.MethodGet, "/elastic/:index/_doc/:doc_id", SearchDocElasticHandler)

}

func GetUsersHandler(c echo.Context) error {
	//var input model.UserFetchInput
	/*if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, "input is invalid")
	}
	repos := c.Get("repositories").(*server.Repositories)

	result, err := repos.UserRepo.List(input)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, "failed to list users")
	}
	*/
	return c.JSON(http.StatusOK, "")
}

func UpsertUserHandler(c echo.Context) error {
	//var req model.UserUpsertInput

	/*if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, "bad request")
	}*/

	//	res, err := user.New().Upsert(req)
	/*
		if err != nil {
			return c.JSON(http.StatusInternalServerError, fmt.Sprintf("failed to get user %v", err))
		}

		if res.IsEmpty() {
			return c.JSON(http.StatusNotFound, "user not found")
		}*/

	return c.JSON(http.StatusAccepted, "accepted")
}

func PostUsersHandler(c echo.Context) error {
	var jsonMap map[string]interface{}

	err := json.NewDecoder(c.Request().Body).Decode(&jsonMap)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, "failed to decode request")
	}

	return nil
}

func PutUserHandler(c echo.Context) error {
	var jsonMap map[string]interface{}

	err := json.NewDecoder(c.Request().Body).Decode(&jsonMap)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, "failed to decode request")
	}

	return nil
}

func DeleteUserHandler(c echo.Context) error {
	//uuid := c.Param("uuid")

	/*		if err := user.New().Delete(model.DeleteUserInput{Uuid: uuid}); err != nil {
			return c.JSON(http.StatusInternalServerError, fmt.Sprintf("failed to delete user %v", err))
		}*/

	return c.JSON(http.StatusNoContent, "deleted")

	return nil
}

func UploadToS3(file http.File) {
}

func GetUserHandler(c echo.Context) error {

	return nil
}

func GetBiderHandler(c echo.Context) error {

	return nil
}

func HealthHandler(c echo.Context) error {
	var health HealthResponse
	health.Status = "Ok"
	health.Version = 1.0

	return c.JSON(200, health)
}

func healthFunc(usersDB, itemsDB *sql.DB) HealthResponse {
	var health HealthResponse
	health.Status = "Ok"
	health.Version = 1.0

	return HealthResponse{UsersDB: DatabaseHealth{Connected: true,
		Stats: usersDB.Stats()},
		ItemsDB: DatabaseHealth{Connected: true,
			Stats: itemsDB.Stats()}}
}

func PostAuctionHandler(c echo.Context) error {

	return c.JSON(200, "response")
}

func PutAuctionHandler(c echo.Context) error {

	return c.JSON(200, "response")
}

func DeleteAuctionHandler(c echo.Context) error {

	return c.JSON(200, "response")
}

func GetAuctionHandler(c echo.Context) error {

	return c.JSON(200, "response")
}
func GetAuctionsHandler(c echo.Context) error {
	return c.JSON(200, "response")
}

func PostItemHandler(c echo.Context) error {

	return c.JSON(200, "response")
}

func PutItemHandler(c echo.Context) error {

	return c.JSON(200, "response")
}
func DeleteItemHandler(c echo.Context) error {

	return c.JSON(200, "response")
}
func GetItemHandler(c echo.Context) error {

	return c.JSON(200, "response")
}
func GetItemsHandler(c echo.Context) error {

	return c.JSON(200, "response")
}

func SetUserHandler(c echo.Context) error {

	return c.JSON(200, "response")
}

func GetIndexElasticHandler(c echo.Context) error {
	index := c.Get("index")
	id := c.Get("id")

	if index == "" || id == "" {
		return c.JSON(400, "invalid or missing params")
	}

	return c.JSON(200, "response")
}

func SearchDocElasticHandler(c echo.Context) error {
	index := c.Get("index")
	id := c.Get("id")

	if index == "" || id == "" {
		return c.JSON(400, "invalid or missing params")
	}

	return c.JSON(200, "response")
}
