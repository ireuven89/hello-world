package utils

import (
	"net/http"
)

func NewHttpClient() http.Client {
	client := http.Client{}

	return client
}
