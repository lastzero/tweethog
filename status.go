package tweethog

import (
	"github.com/dghubble/go-twitter/twitter"
	"strings"
	"sync"
	"time"
	"log"
)

type Status struct {
	tweet      *twitter.Tweet
	stream     *Stream
	config     *Config
	client     *twitter.Client
	lastAction *LastAction
}

const (
	CompactTime = "2006/01/02 15:04:05"
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

func (status *Status) MatchesFilter(filter *Filters) bool {
	if !filter.Retweets && status.IsRetweet() {
		return false
	}

	if !filter.Replies && status.IsReply() {
		return false
	}

	if !filter.Via && strings.Contains(status.GetText(), "via @") {
		return false
	}

	if !filter.URLs && status.ContainsUrl() {
		return false
	}

	if (status.GetFollowersCount() > filter.MaxFollowers && filter.MaxFollowers > 0) ||
		status.GetFollowersCount() < filter.MinFollowers {
		return false
	}

	if (status.GetFriendsCount() > filter.MaxFollowing && filter.MaxFollowing > 0) ||
		status.GetFriendsCount() < filter.MinFollowing {
		return false
	}

	if strings.Count(status.GetText(), "#") > filter.MaxTags {
		return false
	}

	if strings.Count(status.GetText(), "@") > filter.MaxMentions {
		return false
	}

	return true
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

func (status *Status) GetName() string {
	return status.tweet.User.Name
}

func (status *Status) Like() {
	createParams := &twitter.FavoriteCreateParams{
		ID: status.GetID(),
	}

	status.stream.client.Favorites.Create(createParams)
}

func (status *Status) SmartLike() {
	now := time.Now()

	if lastUserLikeTime, ok := status.lastAction.userNames[status.GetScreenName()]; ok && now.Sub(lastUserLikeTime) < time.Duration(48*time.Hour) {
		log.Println("Skipped like because of user rate limit 🐷")
		return
	}

	if now.Sub(status.lastAction.lastLike) < time.Duration(120*time.Second) {
		log.Println("Skipped like because of global rate limit ⏳")
		return
	}

	status.lastAction.Lock()
	status.lastAction.lastLike = now
	status.lastAction.userNames[status.GetScreenName()] = now
	status.lastAction.Unlock()

	randomSeconds := time.Duration(GetRandomInt(45, 300))

	log.Printf("Going to like status %d after %d seconds ⏰\n", status.GetID(), randomSeconds)

	time.Sleep(time.Second * randomSeconds)

	createParams := &twitter.FavoriteCreateParams{
		ID: status.GetID(),
	}

	status.stream.client.Favorites.Create(createParams)

	log.Printf("\nLiked status %d ❤️\n", status.GetID())
}
