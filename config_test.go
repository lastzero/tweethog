package tweethog

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewConfig(t *testing.T) {
	config := NewConfig()

	assert.IsType(t, &Config{}, config)
}

func TestConfig_LoadFromFile(t *testing.T) {
	config := NewConfig()

	config.SetValuesFromFile("config.example.yml")

	assert.Equal(t, "YOUR_TOKEN_SECRET", config.AccessSecret)
	assert.Equal(t, "YOUR_TOKEN", config.AccessToken)
	assert.Equal(t, "YOUR_KEY", config.ConsumerKey)
	assert.Equal(t, "YOUR_SECRET", config.ConsumerSecret)
}
