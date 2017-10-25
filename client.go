package tweethog

import (
	"github.com/dghubble/oauth1"
	"fmt"
	"time"
	"strings"
	"os"
	"os/signal"
	"syscall"
	"github.com/urfave/cli"
	"github.com/dghubble/go-twitter/twitter"
	"log"
)

type Client struct {
	config *Config
	twitterClient *twitter.Client
}

func NewClient(config *Config) *Client {
	oauth1Config := oauth1.NewConfig(config.ConsumerKey, config.ConsumerSecret)
	oauth1Token := oauth1.NewToken(config.AccessToken, config.AccessSecret)
	httpClient := oauth1Config.Client(oauth1.NoContext, oauth1Token)

	// Twitter Client
	twitterClient := twitter.NewClient(httpClient)

	result := new(Client)

	result.twitterClient = twitterClient

	return result
}

func (client *Client) StreamTopic(c *cli.Context) error {
	demux := twitter.NewSwitchDemux()

	demux.Tweet = func(twitterTweet *twitter.Tweet) {
		tweet := NewTweet(twitterTweet, client)
		tweet.handleTweet(c)
	}

	fmt.Printf("Started streaming Twitter status updates on %s...\n", time.Now().Format(time.RFC1123))

	if c.GlobalString("config") != "config.yml" {
		fmt.Printf("Config      : %s\n", c.GlobalString("config"))
	}

	fmt.Printf("Topics       : %s\n", strings.Join(c.GlobalStringSlice("topic"), ", "))
	fmt.Printf("Languages    : %s\n", strings.Join(c.GlobalStringSlice("lang"), ", "))
	fmt.Printf("URLs         : %t\n", c.GlobalBool("urls"))
	fmt.Printf("Retweets     : %t\n", c.GlobalBool("retweets"))
	fmt.Printf("Replies      : %t\n", c.GlobalBool("replies"))
	fmt.Printf("Via          : %t\n", c.GlobalBool("via"))
	fmt.Printf("Max mentions : %d\n", c.GlobalInt("max-mentions"))
	fmt.Printf("Max tags     : %d\n", c.GlobalInt("max-tags"))
	fmt.Printf("Like tweets  : %t\n", c.GlobalBool("like") || c.GlobalBool("smart-like"))

	// FILTER
	filterParams := &twitter.StreamFilterParams{
		Track:         c.GlobalStringSlice("topic"),
		StallWarnings: twitter.Bool(false),
		Language:      c.GlobalStringSlice("lang"),
	}

	stream, err := client.twitterClient.Streams.Filter(filterParams)

	if err != nil {
		log.Fatal(err)
	}

	// Receive messages until stopped or stream quits
	go demux.HandleChan(stream.Messages)

	// Wait for SIGINT and SIGTERM (HIT CTRL-C)
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	log.Println(<-ch)

	fmt.Println("Stopping Twitter stream...")
	stream.Stop()

	return nil
}
