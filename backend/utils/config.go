package utils

import (
	"encoding/json"
	"math"
	"os"
	"strings"
)

type Config struct {
	TenantEndpoint string               `json:"endpoint"`
	Databases      []DataBaseConnection `json:"databaseConnections"`
}

type DataBaseConnection struct {
	Name     string `json:"name"`
	Host     string `json:"host"`
	User     string `json:"user"`
	Password string `json:"password"`
}

func GetConfiguration(dir string) (Config, error) {
	var config Config
	file, err := os.Open(dir + "/config.json")

	if err != nil {
		return Config{}, err
	}

	if err = json.NewDecoder(file).Decode(&config); err != nil {
		return Config{}, err
	}

	return config, nil
}

// Levenshtein - returns the distance between a and b
func Levinstein(a, b string) float32 {
	aChars := strings.Split(a, "")
	bChars := strings.Split(b, "")
	absDistance := math.Abs(float64(float32(len(bChars) - len(aChars))))

	for index := range aChars {
		if index >= len(bChars) {
			break
		}
		if aChars[index] != bChars[index] {
			absDistance++
		}
	}

	return float32(absDistance / float64(len(bChars)))
}
