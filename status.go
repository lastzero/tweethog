package tweethog

import (
	"encoding/json"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/subosito/shorturl"
	"github.com/mvdan/xurls"
	"github.com/ddliu/go-httpclient"
	"time"
	"sync"
	"strings"
	"log"
	"os"
	"errors"
	"crypto/md5"
	"encoding/hex"
)

type Status struct {
	tweet           *twitter.Tweet
	stream          *Stream
	config          *Config
	client          *twitter.Client
	imageUrl        string
	imageUrlChecked bool
	lastAction      *LastAction
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

	if filter.ImagesOnly {
		return status.ContainsImage()
	}

	if !filter.URLs && status.ContainsUrl() {
		return false
	}

	return true
}

func (status *Status) GetAllUrls() []string {
	urls := xurls.Strict.FindAllString(status.GetText(), -1)

	for index, url := range urls {
		expandedUrl, _ := shorturl.Expand(url)
		urls[index] = string(expandedUrl)
	}

	return urls
}

func (status *Status) ContainsImage() bool {
	return status.GetImageUrl() != ""
}

func (status *Status) GetImageUrl() string {
	if status.imageUrlChecked {
		return status.imageUrl
	}

	status.imageUrlChecked = true

	urls := status.GetAllUrls()

	for _, url := range urls {
		resp, err := httpclient.Get(url, map[string]string{})

		if err != nil {
			log.Println(err)
			continue
		}

		if resp.StatusCode != 200 {
			continue
		}

		body, err := resp.ToString()

		if err != nil {
			log.Println(err)
		}

		bodyUrls := xurls.Strict.FindAllString(body, -1)

		for _, bodyUrl := range bodyUrls {
			if strings.Contains(bodyUrl, "https://pbs.twimg.com/media/") && strings.Contains(bodyUrl, ":large") {
				status.imageUrl = bodyUrl
				return bodyUrl
			}
		}
	}

	return ""
}

func (status *Status) SaveImageToFile(path string) (string, error) {
	if !status.ContainsImage() {
		return "", errors.New("contains no image")
	}

	imageUrl := status.GetImageUrl()
	resp, err := httpclient.Get(imageUrl, map[string]string{})

	if err != nil {
		return "", err
	} else if resp.StatusCode != 200 {
		return "", errors.New("could not download image")
	}

	imageBytes, err := resp.ReadAll()

	if err != nil {
		return "", err
	}

	hasher := md5.New()
	hasher.Write(imageBytes)
	hash := hex.EncodeToString(hasher.Sum(nil))

	filename := path + "/" + hash + ".jpg"

	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, os.ModePerm)

	if err != nil {
		return "", err
	}

	defer f.Close()

	_, err = f.Write(imageBytes)

	if err != nil {
		return "", err
	}

	return filename, nil
}

func (status *Status) ContainsUrl() bool {
	return strings.Contains(status.tweet.Text, "://")
}

func (status *Status) GetID() int64 {
	return status.tweet.ID
}

func (status *Status) GetIDString() string {
	return status.tweet.IDStr
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

func (status *Status) GetAsJson() (string, error) {
	if encoded, err := json.Marshal(status.tweet); err != nil {
		return "", err
	} else {
		return string(encoded), nil
	}
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

	log.Printf("Liked status %d ❤️\n", status.GetID())
}
