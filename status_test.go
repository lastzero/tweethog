package tweethog

import (
	"github.com/dghubble/go-twitter/twitter"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewStatus(t *testing.T) {
	tweet := &twitter.Tweet{}
	config := NewConfig()
	config.MaxTags = 1234
	stream := NewStream(config)
	status := NewStatus(tweet, stream)
	assert.Equal(t, config.MaxTags, status.config.MaxTags)
	assert.IsType(t, &Status{}, status)
}
