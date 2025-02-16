package utils

type ServiceHealthCheck struct {
	ServiceStatus string     `json:"status"`
	DBStatus      []DbStatus `json:"DBStatus"`
}

type DbStatus struct {
	Name      string `json:"name"`
	Connected bool   `json:"connected"`
}
