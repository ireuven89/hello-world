package model

type HealthResponse struct {
	Status  string `json:"status"`
	Version int    `json:"version"`
}
