package tweethog

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewClient(t *testing.T) {
	config := NewConfig()
	client := NewStream(config)

	assert.IsType(t, &Stream{}, client)
}
