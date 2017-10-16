package main

import (
	"fmt"
	"github.com/urfave/cli"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"log"
	"os"
	"os/signal"
	"syscall"
	"github.com/kylelemons/go-gypsy/yaml"
	"strings"
	"time"
)

func main() {
	app := cli.NewApp()
	app.Usage = "Stream, filter and like Twitter status updates"
	app.Version = "0.2.0"
	app.Copyright = "Michael Mayer <michael@liquidbytes.net>"

	app.Flags = cliFlags

	app.Action = streamTweets

	app.Run(os.Args)
}

type TweethogConfig struct {
	consumerKey    string
	consumerSecret string
	accessToken    string
	accessSecret   string
}

var cliFlags = []cli.Flag{
	cli.StringSliceFlag{
		Name: "topic, t",
		Usage: "Stream filter topic (cat, dog, fish, ...)",
	},
	cli.StringSliceFlag{
		Name: "lang, l",
		Usage: "Stream filter language (en, de, fr, ...)",
	},
	cli.IntFlag{
		Name:  "max-followers",
		Value: 5000,
		Usage: "User max follower count (0 for unlimited)",
	},
	cli.IntFlag{
		Name:  "min-followers",
		Value: 5,
		Usage: "User min follower count",
	},
	cli.IntFlag{
		Name:  "max-following",
		Value: 5000,
		Usage: "User max following count (0 for unlimited)",
	},
	cli.IntFlag{
		Name:  "min-following",
		Value: 5,
		Usage: "User min following count",
	},
	cli.BoolFlag{
		Name:  "no-retweets",
		Usage: "Exclude tweets starting with RT or @",
	},
	cli.BoolFlag{
		Name:  "no-urls",
		Usage: "Exclude tweets containing URLs",
	},
	cli.BoolFlag{
		Name:  "like",
		Usage: "Like tweets",
	},
	cli.StringFlag{
		Name:  "config, c",
		Usage: "Config file name",
		Value: "config.yml",
	},
}

func getConfig(fileName string) (TweethogConfig, error) {
	var result TweethogConfig

	config, configErr := yaml.ReadFile(fileName)

	if configErr != nil {
		return result, cli.NewExitError(configErr, 1)
	}

	if consumerKey, err := config.Get("consumer_key"); err != nil {
		return result, cli.NewExitError(err, 1)
	} else {
		result.consumerKey = consumerKey
	}

	if consumerSecret, err := config.Get("consumer_secret"); err != nil {
		return result, cli.NewExitError(err, 1)
	} else {
		result.consumerSecret = consumerSecret
	}

	if accessToken, err := config.Get("access_token"); err != nil {
		return result, cli.NewExitError(err, 1)
	} else {
		result.accessToken = accessToken
	}

	if accessSecret, err := config.Get("access_secret"); err != nil {
		return result, cli.NewExitError(err, 1)
	} else {
		result.accessSecret = accessSecret
	}

	return result, nil
}

func newTwitterClient(tweethogConfig *TweethogConfig) *twitter.Client {
	oauth1Config := oauth1.NewConfig(tweethogConfig.consumerKey, tweethogConfig.consumerSecret)
	oauth1Token := oauth1.NewToken(tweethogConfig.accessToken, tweethogConfig.accessSecret)
	httpClient := oauth1Config.Client(oauth1.NoContext, oauth1Token)

	// Twitter Client
	result := twitter.NewClient(httpClient)

	return result
}

func streamTweets(c *cli.Context) error {
	if len(c.GlobalStringSlice("topic")) == 0 {
		cli.ShowAppHelp(c)
		return nil
	}

	tweethogConfig, err := getConfig(c.GlobalString("config"))

	if err != nil {
		return err
	}

	client := newTwitterClient(&tweethogConfig)

	demux := twitter.NewSwitchDemux()

	demux.Tweet = func(tweet *twitter.Tweet) {
		handleTweet(tweet, c, client)
	}

	// demux.All = func(message interface{}) {
	//	fmt.Printf("%#v", message)
	// }

	fmt.Printf("Started streaming Twitter status updates on %s...\n", time.Now().Format(time.RFC1123))

	if c.GlobalString("config") != "config.yml" {
		fmt.Printf("Config      : %s\n", c.GlobalString("config"))
	}

	fmt.Printf("Topics      : %s\n", strings.Join(c.GlobalStringSlice("topic"), ", "))
	fmt.Printf("Languages   : %s\n", strings.Join(c.GlobalStringSlice("lang"), ", "))
	fmt.Printf("URLs        : %t\n", !c.GlobalBool("no-urls"))
	fmt.Printf("Retweets    : %t\n", !c.GlobalBool("no-retweets"))
	fmt.Printf("Like tweets : %t\n", c.GlobalBool("like"))

	// FILTER
	filterParams := &twitter.StreamFilterParams{
		Track:         c.GlobalStringSlice("topic"),
		StallWarnings: twitter.Bool(false),
		Language:      c.GlobalStringSlice("lang"),
	}

	stream, err := client.Streams.Filter(filterParams)

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

func handleTweet(tweet *twitter.Tweet, c *cli.Context, client *twitter.Client) {
	if c.GlobalBool("no-retweets") && (tweet.Retweeted ||
		strings.Contains(tweet.Text, "via @") ||
		strings.HasPrefix(tweet.Text, "RT") ||
		strings.HasPrefix(tweet.Text, "@")) {
		fmt.Print(".")
		return
	}

	if c.GlobalBool("no-urls") && strings.Contains(tweet.Text, "://") {
		fmt.Print(".")
		return
	}

	if (tweet.User.FollowersCount > c.GlobalInt("max-followers") && c.GlobalInt("max-followers") > 0) ||
		tweet.User.FollowersCount < c.GlobalInt("min-followers") {
		fmt.Print(".")
		return
	}

	if (tweet.User.FriendsCount > c.GlobalInt("max-following") && c.GlobalInt("max-following") > 0) ||
		tweet.User.FriendsCount < c.GlobalInt("min-following") {
		fmt.Print(".")
		return
	}

	fmt.Printf("\nID: %d  Date: %s  User: @%s  Following: %d  Followers: %d  Likes: %d\n>>> %s\n",
		tweet.ID,
		tweet.CreatedAt,
		tweet.User.ScreenName,
		tweet.User.FriendsCount,
		tweet.User.FollowersCount,
		tweet.User.FavouritesCount,
		tweet.Text)

	if c.GlobalBool("like") {
		createParams := &twitter.FavoriteCreateParams{
			ID: tweet.ID,
		}

		client.Favorites.Create(createParams)

		fmt.Println("Liked ❤️")
	}
}