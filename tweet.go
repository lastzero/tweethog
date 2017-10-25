package tweethog

import (
	"time"
	"sync"
	"strings"
	"fmt"
	"github.com/urfave/cli"
	"github.com/dghubble/go-twitter/twitter"
	"math/rand"
)

type Tweet struct {
	twitterTweet *twitter.Tweet
	client       *Client
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

func NewTweet(twitterTweet *twitter.Tweet, client *Client) *Tweet {
	result := new(Tweet)

	result.twitterTweet = twitterTweet
	result.client = client

	return result
}

func (tweet *Tweet) handleTweet(c *cli.Context) {
	if !c.GlobalBool("retweets") && (tweet.twitterTweet.Retweeted || strings.HasPrefix(tweet.twitterTweet.Text, "RT")) {
		fmt.Print(".")
		return
	}

	if !c.GlobalBool("replies") && strings.HasPrefix(tweet.twitterTweet.Text, "@") {
		fmt.Print(".")
		return
	}

	if !c.GlobalBool("via") && strings.Contains(tweet.twitterTweet.Text, "via @") {
		fmt.Print(".")
		return
	}

	if !c.GlobalBool("urls") && tweetContainsUrl(tweet.twitterTweet) {
		fmt.Print(".")
		return
	}

	if (tweet.twitterTweet.User.FollowersCount > c.GlobalInt("max-followers") && c.GlobalInt("max-followers") > 0) ||
		tweet.twitterTweet.User.FollowersCount < c.GlobalInt("min-followers") {
		fmt.Print(".")
		return
	}

	if (tweet.twitterTweet.User.FriendsCount > c.GlobalInt("max-following") && c.GlobalInt("max-following") > 0) ||
		tweet.twitterTweet.User.FriendsCount < c.GlobalInt("min-following") {
		fmt.Print(".")
		return
	}

	if strings.Count(tweet.twitterTweet.Text, "#") > c.GlobalInt("max-tags") {
		fmt.Print(".")
		return
	}

	if strings.Count(tweet.twitterTweet.Text, "@") > c.GlobalInt("max-mentions") {
		fmt.Print(".")
		return
	}

	fmt.Printf("\nID: %d  Date: %s  User: @%s  Following: %d  Followers: %d  Likes: %d\n>>> %s\n",
		tweet.twitterTweet.ID,
		tweet.twitterTweet.CreatedAt,
		tweet.twitterTweet.User.ScreenName,
		tweet.twitterTweet.User.FriendsCount,
		tweet.twitterTweet.User.FollowersCount,
		tweet.twitterTweet.User.FavouritesCount,
		tweet.twitterTweet.Text)

	if c.GlobalBool("smart-like") {
		go tweet.smartLikeTweet()
	} else if c.GlobalBool("like") {
		tweet.Like()
	}
}

func tweetContainsUrl(tweet *twitter.Tweet) bool {
	return strings.Contains(tweet.Text, "://")
}

func (tweet *Tweet) GetID() int64 {
	return tweet.twitterTweet.ID
}

func (tweet *Tweet) GetScreenName() string {
	return tweet.twitterTweet.User.ScreenName
}

func (tweet *Tweet) Like() {
	createParams := &twitter.FavoriteCreateParams{
		ID: tweet.GetID(),
	}

	tweet.client.twitterClient.Favorites.Create(createParams)

	fmt.Printf("Liked tweet %d ❤️\n", tweet.GetID())
}

func GetRandomInt(min, max int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(max-min) + min
}

func (tweet *Tweet) smartLikeTweet() {
	now := time.Now()

	if lastUserLikeTime, ok := likeUserNames.userNames[tweet.GetScreenName()]; ok && now.Sub(lastUserLikeTime) < time.Duration(48*time.Hour) {
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

	likeUserNames.userNames[tweet.GetScreenName()] = now
	likeUserNames.Unlock()

	randomSeconds := time.Duration(GetRandomInt(45, 300))

	fmt.Printf("Going to like tweet %d after %d seconds ⏰\n", tweet.GetID(), randomSeconds)

	time.Sleep(time.Second * randomSeconds)

	createParams := &twitter.FavoriteCreateParams{
		ID: tweet.GetID(),
	}

	tweet.client.twitterClient.Favorites.Create(createParams)

	fmt.Printf("\nLiked tweet %d ❤️\n", tweet.GetID())
}
