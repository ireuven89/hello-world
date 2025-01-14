package tests

import (
	"github.com/ireuven89/hello-world/backend/utils"
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"testing"
)

func TestConfig(t *testing.T) {
	path, err := filepath.Abs(".")
	config, err := utils.GetConfiguration(path)

	assert.Nil(t, err)
	assert.NotEmpty(t, config)

}

func TestLevinstein(t *testing.T) {
	distance := utils.Levinstein("ba", "b")

	assert.Equal(t, float32(1), distance)
}
