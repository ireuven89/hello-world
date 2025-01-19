package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
)

type ConfigurationJson struct {
	ElasticUrlDev     string `json:"elastic_url_dev"`
	ElasticUrlLocal   string `json:"elastic_url_local"`
	ElasticUrlStaging string `json:"elastic_url_staging"`
}

func GetConfigJson() ConfigurationJson {
	var config ConfigurationJson
	_, filePath, _, ok := runtime.Caller(0)
	if !ok {
		return config
	}
	dirPath := filepath.Dir(filePath)
	file, err := os.ReadFile(dirPath + "/dev.json")
	err = json.Unmarshal(file, &config)
	if err != nil {
		return ConfigurationJson{}
	}

	return config
}
