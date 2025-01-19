package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	TenantEndpoint string                        `json:"endpoint"`
	Databases      map[string]DataBaseConnection `json:"databaseConnections"`
	ServicePort    string                        `json:"servicePort"`
}

type DataBaseConnection struct {
	Host     string `json:"host"`
	UserName string `json:"user"`
	Password string `json:"password"`
}

const configPath = "%s/config/%s.json"

func LoadConfig(service, env string) (Config, error) {
	var config Config
	dir, err := filepath.Abs(service)
	file, err := os.Open(fmt.Sprintf(configPath, dir, env))

	if err != nil {
		return Config{}, err
	}

	if err = json.NewDecoder(file).Decode(&config); err != nil {
		return Config{}, err
	}

	return config, nil
}
