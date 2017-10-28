package tweethog

import (
	"github.com/dghubble/go-twitter/twitter"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewStatus(t *testing.T) {
	tweet := &twitter.Tweet{}
	config := NewConfig()
	config.Filter.MaxTags = 1234
	stream := NewStream(config)
	status := NewStatus(tweet, stream)
	assert.Equal(t, config.Filter.MaxTags, status.config.Filter.MaxTags)
	assert.IsType(t, &Status{}, status)
}

func TestStatus_GetId(t *testing.T) {
	tweet := &twitter.Tweet{
		ID: int64(123),
		IDStr: "356",
		Text: "I love https://twitter.com/i/web/status/923880846393790467",
	}

	config := NewConfig()
	stream := NewStream(config)
	status := NewStatus(tweet, stream)

	ID := status.GetID()
	IDStr := status.GetIDString()

	assert.Equal(t, int64(123), ID)
	assert.Equal(t, "356", IDStr)
}