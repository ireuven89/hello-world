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
	GroupId  string `bson:"groupId"`
	UserName string `json:"user"`
	Password string `json:"password"`
}

const configPath = "%s/config/%s.json"

func LoadConfig(service, env string) (Config, error) {
	var config Config
	dir, err := filepath.Abs(service)
	configFile := fmt.Sprintf(configPath, dir, env)
	file, err := os.Open(configFile)

	if err != nil {
		fmt.Printf("failed to load config %v", err)
		return Config{}, err
	}

	fmt.Printf("env is %s", env)
	fmt.Printf("config file is %v\n", dir)

	if err = json.NewDecoder(file).Decode(&config); err != nil {
		return Config{}, err
	}
	fmt.Printf(config.toString())

	return config, nil
}

func (c *Config) toString() string {
	return fmt.Sprintf("dbs: %v, tenanat: %v, servicePort: %v\n", c.Databases, c.TenantEndpoint, c.ServicePort)
}
