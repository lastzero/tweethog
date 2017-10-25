package tweethog

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	config := NewConfig()
	client := NewClient(config)

	assert.IsType(t, new(Client), client)
}