package routes

import (
	"github.com/labstack/echo/v4"
	"strconv"
)

type HealthResponse struct {
	Status  string `json:"status"`
	Version int    `json:"version"`
}

type Pagination struct {
	Page int `json:"Page" default:"0"`
	Size int `json:"Size" default:"20"`
}

func getPagination(c echo.Context) Pagination {
	var page Pagination

	pageRequest, err := strconv.Atoi(c.QueryParam("page"))

	if err != nil {
		c.Logger().Warn("invalid page request, will default page")
	} else {
		page.Page = pageRequest
	}

	size, err := strconv.Atoi(c.QueryParam("size"))

	if err != nil {
		c.Logger().Warn("invalid page request, will default size")
	} else {
		page.Size = size
	}
	return page
}
