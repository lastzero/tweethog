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
	"math/rand"
	"sync"
)

func main() {
	app := cli.NewApp()
	app.Usage = "Stream, filter and react to Twitter status updates"
	app.Version = "0.4.2"
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

type LastLike struct {
	lastLike time.Time
	sync.Mutex
}

var lastLike LastLike

type LikeUserNames struct {
	userNames map[string]time.Time
	sync.Mutex
}

var likeUserNames LikeUserNames

var cliFlags = []cli.Flag{
	cli.StringSliceFlag{
		Name:  "topic, t",
		Usage: "Stream filter topic e.g. cat, dog, fish",
	},
	cli.StringSliceFlag{
		Name:  "lang, l",
		Usage: "Stream filter language e.g. en, de, fr",
	},
	cli.IntFlag{
		Name:  "max-followers",
		Usage: "User max followers, 0 for unlimited",
	},
	cli.IntFlag{
		Name:  "min-followers",
		Value: 5,
		Usage: "User min followers",
	},
	cli.IntFlag{
		Name:  "max-following",
		Usage: "User max following, 0 for unlimited",
	},
	cli.IntFlag{
		Name:  "min-following",
		Value: 5,
		Usage: "User min following",
	},
	cli.IntFlag{
		Name:  "max-tags",
		Value: 2,
		Usage: "Max number of hash #tags",
	},
	cli.IntFlag{
		Name:  "max-mentions",
		Value: 1,
		Usage: "Max number of @mentions",
	},
	cli.BoolFlag{
		Name:  "retweets",
		Usage: "Include tweets starting with RT",
	},
	cli.BoolFlag{
		Name:  "replies",
		Usage: "Include tweets starting with @",
	},
	cli.BoolFlag{
		Name:  "via",
		Usage: "Include tweets containing via @",
	},
	cli.BoolFlag{
		Name:  "urls",
		Usage: "Include tweets containing URLs",
	},
	cli.BoolFlag{
		Name:  "like",
		Usage: "Like all matching tweets",
	},
	cli.BoolFlag{
		Name:  "smart-like",
		Usage: "Likes tweets with random delay and rate limit",
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
	if !c.GlobalBool("retweets") && (tweet.Retweeted || strings.HasPrefix(tweet.Text, "RT")) {
		fmt.Print(".")
		return
	}

	if !c.GlobalBool("replies") && strings.HasPrefix(tweet.Text, "@") {
		fmt.Print(".")
		return
	}

	if !c.GlobalBool("via") && strings.Contains(tweet.Text, "via @") {
		fmt.Print(".")
		return
	}

	if !c.GlobalBool("urls") && tweetContainsUrl(tweet) {
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

	if strings.Count(tweet.Text, "#") > c.GlobalInt("max-tags") {
		fmt.Print(".")
		return
	}

	if strings.Count(tweet.Text, "@") > c.GlobalInt("max-mentions") {
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

	if c.GlobalBool("smart-like") {
		go smartLikeTweet(tweet, client)
	} else if c.GlobalBool("like") {
		likeTweet(tweet, client)
	}
}

func tweetContainsUrl(tweet *twitter.Tweet) bool {
	return strings.Contains(tweet.Text, "://")
}

func likeTweet(tweet *twitter.Tweet, client *twitter.Client) {
	createParams := &twitter.FavoriteCreateParams{
		ID: tweet.ID,
	}

	client.Favorites.Create(createParams)

	fmt.Printf("Liked tweet %d ❤️\n", tweet.ID)
}

func random(min, max int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(max-min) + min
}

func smartLikeTweet(tweet *twitter.Tweet, client *twitter.Client) {
	now := time.Now()

	if lastUserLikeTime, ok := likeUserNames.userNames[tweet.User.ScreenName]; ok && now.Sub(lastUserLikeTime) < time.Duration(48*time.Hour) {
		fmt.Println("Skipped Like because of user rate limit 🐷")
		return
	}

	if now.Sub(lastLike.lastLike) < time.Duration(120*time.Second) {
		fmt.Println("Skipped Like because of global rate limit ⏳")
		return
	}

	lastLike.Lock()
	lastLike.lastLike = now
	lastLike.Unlock()

	likeUserNames.Lock()

	if likeUserNames.userNames == nil {
		likeUserNames.userNames = make(map[string]time.Time)
	}

	likeUserNames.userNames[tweet.User.ScreenName] = now
	likeUserNames.Unlock()

	randomSeconds := time.Duration(random(45, 300))

	fmt.Printf("Going to like tweet %d after %d seconds ⏰\n", tweet.ID, randomSeconds)

	time.Sleep(time.Second * randomSeconds)

	createParams := &twitter.FavoriteCreateParams{
		ID: tweet.ID,
	}

	client.Favorites.Create(createParams)

	fmt.Printf("\nLiked tweet %d ❤️\n", tweet.ID)
}
