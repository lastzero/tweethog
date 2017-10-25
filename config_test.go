package tweethog

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestNewConfig(t *testing.T) {
	config := NewConfig()

	assert.IsType(t, new(Config), config)
}

func TestConfig_LoadFromFile(t *testing.T) {
	config := NewConfig()

	config.LoadFromFile("config.example.yml")

	assert.Equal(t, "YOUR_TOKEN_SECRET", config.AccessSecret)
	assert.Equal(t, "YOUR_TOKEN", config.AccessToken)
	assert.Equal(t, "YOUR_KEY", config.ConsumerKey)
	assert.Equal(t, "YOUR_SECRET", config.ConsumerSecret)
}