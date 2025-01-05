package utils

import (
	"github.com/labstack/echo/v4"
	"strconv"
)

type DefaultRequest struct {
	Page Pagination `json:"page"`
}
type Pagination struct {
	Page   int `json:"Page"`
	Size   int `json:"Size"`
	Offset int `json:"offset"`
}

const defaultSize = 20

func getPagination(c echo.Context) Pagination {
	var page Pagination

	pageRequest, err := strconv.Atoi(c.QueryParam("page"))

	if err != nil {
		c.Logger().Warn("invalid page request, will default page")
	}

	page.Page = pageRequest

	size, err := strconv.Atoi(c.QueryParam("size"))

	if err != nil {
		page.Size = defaultSize
		c.Logger().Warn("invalid page request, will default size")
	} else {
		page.Size = size
	}

	offset, err := strconv.Atoi(c.QueryParam("offset"))

	if err != nil {
		c.Logger().Warn("invalid offset, will default")
	}

	page.Offset = offset

	return page
}
