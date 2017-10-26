package tweethog

import (
	"fmt"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

type Stream struct {
	config *Config
	client *twitter.Client
}

func NewStream(config *Config) *Stream {
	oauth1Config := oauth1.NewConfig(config.ConsumerKey, config.ConsumerSecret)
	oauth1Token := oauth1.NewToken(config.AccessToken, config.AccessSecret)
	httpClient := oauth1Config.Client(oauth1.NoContext, oauth1Token)

	// Twitter Stream
	twitterClient := twitter.NewClient(httpClient)

	result := &Stream{
		config: config,
		client: twitterClient,
	}

	return result
}

func (stream *Stream) GetConfig() *Config {
	return stream.config
}

func (stream *Stream) Start() error {
	demux := twitter.NewSwitchDemux()

	demux.Tweet = func(tweet *twitter.Tweet) {
		status := NewStatus(tweet, stream)
		status.Handle()
	}

	fmt.Printf("Started streaming Twitter status updates on %s...\n", time.Now().Format(CompactTime))

	fmt.Printf("Topics       : %s\n", strings.Join(stream.config.Topics, ", "))
	fmt.Printf("Languages    : %s\n", strings.Join(stream.config.Languages, ", "))
	fmt.Printf("URLs         : %t\n", stream.config.URLs)
	fmt.Printf("Retweets     : %t\n", stream.config.Retweets)
	fmt.Printf("Replies      : %t\n", stream.config.Replies)
	fmt.Printf("Via          : %t\n", stream.config.Via)
	fmt.Printf("Max mentions : %d\n", stream.config.MaxMentions)
	fmt.Printf("Max tags     : %d\n", stream.config.MaxTags)
	fmt.Printf("Like tweets  : %t\n", stream.config.Like || stream.config.SmartLike)

	// FILTER
	filterParams := &twitter.StreamFilterParams{
		Track:         stream.config.Topics,
		StallWarnings: twitter.Bool(false),
		Language:      stream.config.Languages,
	}

	filterStream, err := stream.client.Streams.Filter(filterParams)

	if err != nil {
		log.Fatal(err)
	}

	// Receive messages until stopped or stream quits
	go demux.HandleChan(filterStream.Messages)

	// Wait for SIGINT and SIGTERM (HIT CTRL-C)
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	log.Println(<-ch)

	fmt.Println("Stopping Twitter stream...")
	filterStream.Stop()

	return nil
}
