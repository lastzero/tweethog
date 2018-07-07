package tweethog

import (
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"log"
	"os"
	"os/signal"
	"syscall"
)

type Stream struct {
	config *Config
	client *twitter.Client
}

func NewStream(config *Config) *Stream {
	oauth1Config := oauth1.NewConfig(config.ConsumerKey, config.ConsumerSecret)
	oauth1Token := oauth1.NewToken(config.AccessToken, config.AccessSecret)
	httpClient := oauth1Config.Client(oauth1.NoContext, oauth1Token)

	result := &Stream{
		config: config,
		client: twitter.NewClient(httpClient),
	}

	return result
}

func (stream *Stream) GetConfig() *Config {
	return stream.config
}

func (stream *Stream) Start(action func(status *Status)) error {
	demux := twitter.NewSwitchDemux()

	demux.Tweet = func(tweet *twitter.Tweet) {
		status := NewStatus(tweet, stream)

		action(status)
	}

	log.Println("Starting Twitter stream...")

	// FILTER
	filterParams := &twitter.StreamFilterParams{
		Track:         stream.config.Filter.Topics,
		StallWarnings: twitter.Bool(true),
		Language:      stream.config.Filter.Languages,
	}

	filterStream, err := stream.client.Streams.Filter(filterParams)

	if err != nil {
		log.Fatal(err)
	}

	// Receive messages until stopped or stream quits
	go demux.HandleChan(filterStream.Messages)

	// Create channel for termination signal
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)

	// Wait for SIGINT and SIGTERM (HIT CTRL-C)
	log.Println(<-ch)

	log.Println("Stopping Twitter stream...")

	filterStream.Stop()

	return nil
}
