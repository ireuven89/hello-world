package routes

import (
	"github.com/labstack/echo/v4"
	"strconv"
)

const defaultPageSize = 20
const defaultPage = 0

type HealthResponse struct {
	Status  string         `json:"status"`
	Version int            `json:"version"`
	UsersDB DatabaseHealth `json:"usersDB"`
	ItemsDB DatabaseHealth `json:"itemsDB"`
	Kafka   string         `json:"kafka"`
}

type DatabaseHealth struct {
	Connected bool        `json:"connected"`
	Stats     interface{} `json:"stats"`
}

type Pagination struct {
	Page int `json:"Page"`
	Size int `json:"Size"`
}

func getPagination(c echo.Context) Pagination {
	var page Pagination

	pageRequest, err := strconv.Atoi(c.QueryParam("page"))

	if err != nil {
		c.Logger().Warn("invalid page request, will default page")
		page.Page = defaultPage
	} else {
		page.Page = pageRequest
	}

	size, err := strconv.Atoi(c.QueryParam("size"))

	if err != nil {
		c.Logger().Warn("invalid page request, will default size")
		page.Size = defaultPageSize
	} else {
		page.Size = size
	}
	return page
}
