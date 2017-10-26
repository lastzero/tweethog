package tweethog

import (
	"fmt"
	"github.com/dghubble/go-twitter/twitter"
	"strings"
	"sync"
	"time"
)

type Status struct {
	tweet      *twitter.Tweet
	stream     *Stream
	config     *Config
	client     *twitter.Client
	lastAction *LastAction
}

const (
	CompactTime = "2006-01-02 15:04:05"
)

type LastAction struct {
	lastLike  time.Time
	userNames map[string]time.Time
	sync.Mutex
}

var lastAction = &LastAction{
	userNames: make(map[string]time.Time),
}

func NewStatus(tweet *twitter.Tweet, stream *Stream) *Status {
	return &Status{
		tweet:      tweet,
		stream:     stream,
		config:     stream.config,
		client:     stream.client,
		lastAction: lastAction,
	}
}

func (status *Status) Handle() {
	if !status.config.Retweets && status.IsRetweet() {
		fmt.Print(".")
		return
	}

	if !status.config.Replies && status.IsReply() {
		fmt.Print(".")
		return
	}

	if !status.config.Via && strings.Contains(status.GetText(), "via @") {
		fmt.Print(".")
		return
	}

	if !status.config.URLs && status.ContainsUrl() {
		fmt.Print(".")
		return
	}

	if (status.GetFollowersCount() > status.config.MaxFollowers && status.config.MaxFollowers > 0) ||
		status.GetFollowersCount() < status.config.MinFollowers {
		fmt.Print(".")
		return
	}

	if (status.GetFriendsCount() > status.config.MaxFollowing && status.config.MaxFollowing > 0) ||
		status.GetFriendsCount() < status.config.MinFollowing {
		fmt.Print(".")
		return
	}

	if strings.Count(status.tweet.Text, "#") > status.config.MaxTags {
		fmt.Print(".")
		return
	}

	if strings.Count(status.tweet.Text, "@") > status.config.MaxMentions {
		fmt.Print(".")
		return
	}

	fmt.Printf("\nID: %d  Date: %s  User: @%s  Following: %d  Followers: %d  Likes: %d\n>>> %s\n",
		status.GetID(),
		status.GetCreatedAt().Format(CompactTime),
		status.GetScreenName(),
		status.GetFriendsCount(),
		status.GetFollowersCount(),
		status.GetFavouritesCount(),
		status.GetText())

	if status.config.SmartLike {
		go status.SmartLike()
	} else if status.config.Like {
		status.Like()
	}
}

func (status *Status) ContainsUrl() bool {
	return strings.Contains(status.tweet.Text, "://")
}

func (status *Status) GetID() int64 {
	return status.tweet.ID
}

func (status *Status) GetCreatedAt() *time.Time {
	// Thu Oct 26 06:01:42 +0000 2017
	result, _ := time.Parse("Mon Jan 02 15:04:05 -0700 2006", status.tweet.CreatedAt)
	return &result
}

func (status *Status) IsRetweet() bool {
	return status.tweet.Retweeted || strings.HasPrefix(status.GetText(), "RT")
}

func (status *Status) IsReply() bool {
	return strings.HasPrefix(status.GetText(), "@")
}

func (status *Status) GetScreenName() string {
	return status.tweet.User.ScreenName
}

func (status *Status) GetFriendsCount() int {
	return status.tweet.User.FriendsCount
}

func (status *Status) GetFollowersCount() int {
	return status.tweet.User.FollowersCount
}

func (status *Status) GetFavouritesCount() int {
	return status.tweet.User.FavouritesCount
}

func (status *Status) GetText() string {
	return status.tweet.Text
}

func (status *Status) Like() {
	createParams := &twitter.FavoriteCreateParams{
		ID: status.GetID(),
	}

	status.stream.client.Favorites.Create(createParams)

	fmt.Printf("Liked status %d ❤️\n", status.GetID())
}

func (status *Status) SmartLike() {
	now := time.Now()

	if lastUserLikeTime, ok := status.lastAction.userNames[status.GetScreenName()]; ok && now.Sub(lastUserLikeTime) < time.Duration(48*time.Hour) {
		fmt.Println("Skipped Like because of user rate limit 🐷")
		return
	}

	if now.Sub(status.lastAction.lastLike) < time.Duration(120*time.Second) {
		fmt.Println("Skipped Like because of global rate limit ⏳")
		return
	}

	status.lastAction.Lock()
	status.lastAction.lastLike = now
	status.lastAction.userNames[status.GetScreenName()] = now
	status.lastAction.Unlock()

	randomSeconds := time.Duration(GetRandomInt(45, 300))

	fmt.Printf("Going to like status %d after %d seconds ⏰\n", status.GetID(), randomSeconds)

	time.Sleep(time.Second * randomSeconds)

	createParams := &twitter.FavoriteCreateParams{
		ID: status.GetID(),
	}

	status.stream.client.Favorites.Create(createParams)

	fmt.Printf("\nLiked status %d ❤️\n", status.GetID())
}
