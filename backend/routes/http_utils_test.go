package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"net/http"
	"testing"
)

func TestGetPagination(t *testing.T) {
	c := echo.New()
	req, _ := http.NewRequest(http.MethodGet, "dummy/url?page=0&size=100", nil)
	writer := http.ResponseWriter(nil)
	actual := getPagination(c.NewContext(req, writer))

	assert.Equal(t, actual.Page, 0)
	assert.Equal(t, actual.Size, 100)
}

func TestGetPaginationDefaultTestCase(t *testing.T) {
	c := echo.New()
	req, _ := http.NewRequest(http.MethodGet, "dummy/url", nil)
	writer := http.ResponseWriter(nil)
	actual := getPagination(c.NewContext(req, writer))

	assert.Equal(t, 0, actual.Page)
	assert.Equal(t, 20, actual.Size)

}
