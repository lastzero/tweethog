package tweethog

import (
	"github.com/kylelemons/go-gypsy/yaml"
	"github.com/urfave/cli"
	"strings"
)

type Config struct {
	// API credentials
	ConsumerKey    string
	ConsumerSecret string
	AccessToken    string
	AccessSecret   string

	// YAML config file name
	ConfigFile string

	// Filters
	Topics       []string
	Languages    []string
	MinFollowers int
	MaxFollowers int
	MinFollowing int
	MaxFollowing int
	MaxTags      int
	MaxMentions  int
	Retweets     bool
	Replies      bool
	Via          bool
	URLs         bool

	// Actions
	Like      bool
	SmartLike bool
}

func NewConfig() *Config {
	return &Config{}
}

func (config *Config) SetValuesFromFile(fileName string) error {
	yamlConfig, err := yaml.ReadFile(fileName)

	if err != nil {
		return err
	}

	config.ConsumerKey, _ = yamlConfig.Get("consumer-key")
	config.ConsumerSecret, _ = yamlConfig.Get("consumer-secret")
	config.AccessToken, _ = yamlConfig.Get("access-token")
	config.AccessSecret, _ = yamlConfig.Get("access-secret")

	if topics, err := yamlConfig.Get("topic"); err == nil {
		config.Topics = strings.Split(topics, ",")
	}

	if languages, err := yamlConfig.Get("lang"); err == nil {
		config.Languages = strings.Split(languages, ",")
	}

	if minFollowers, err := yamlConfig.GetInt("min-followers"); err == nil {
		config.MinFollowers = int(minFollowers)
	}

	if maxFollowers, err := yamlConfig.GetInt("max-followers"); err == nil {
		config.MaxFollowers = int(maxFollowers)
	}

	if minFollowing, err := yamlConfig.GetInt("min-following"); err == nil {
		config.MinFollowing = int(minFollowing)
	}

	if maxFollowing, err := yamlConfig.GetInt("max-following"); err == nil {
		config.MaxFollowing = int(maxFollowing)
	}

	if maxTags, err := yamlConfig.GetInt("max-tags"); err == nil {
		config.MaxTags = int(maxTags)
	}

	if maxMentions, err := yamlConfig.GetInt("max-mentions"); err == nil {
		config.MaxMentions = int(maxMentions)
	}

	if retweets, err := yamlConfig.GetBool("retweets"); err == nil {
		config.Retweets = retweets
	}

	if replies, err := yamlConfig.GetBool("replies"); err == nil {
		config.Replies = replies
	}

	if via, err := yamlConfig.GetBool("via"); err == nil {
		config.Via = via
	}

	if urls, err := yamlConfig.GetBool("urls"); err == nil {
		config.URLs = urls
	}

	return nil
}

func (config *Config) SetValuesFromCliContext(c *cli.Context) error {
	if c.GlobalIsSet("consumer-key") {
		config.ConsumerKey = c.GlobalString("consumer-key")
	}

	if c.GlobalIsSet("consumer-secret") {
		config.ConsumerSecret = c.GlobalString("consumer-secret")
	}

	if c.GlobalIsSet("access-token") {
		config.AccessToken = c.GlobalString("access-token")
	}

	if c.GlobalIsSet("access-secret") {
		config.AccessSecret = c.GlobalString("access-secret")
	}

	if c.GlobalIsSet("config-file") || c.GlobalIsSet("c") {
		config.ConfigFile = c.GlobalString("config-file")
	}

	if c.GlobalIsSet("topic") || c.GlobalIsSet("t") {
		config.Topics = c.GlobalStringSlice("topic")
	}

	if c.GlobalIsSet("lang") || c.GlobalIsSet("l") {
		config.Languages = c.GlobalStringSlice("lang")
	}

	if c.GlobalIsSet("min-followers") {
		config.MinFollowers = c.GlobalInt("min-followers")
	}

	if c.GlobalIsSet("max-followers") {
		config.MaxFollowers = c.GlobalInt("max-followers")
	}

	if c.GlobalIsSet("min-following") {
		config.MinFollowing = c.GlobalInt("min-following")
	}

	if c.GlobalIsSet("max-following") {
		config.MaxFollowing = c.GlobalInt("max-following")
	}

	if c.GlobalIsSet("max-tags") {
		config.MaxTags = c.GlobalInt("max-tags")
	}

	if c.GlobalIsSet("max-mentions") {
		config.MaxMentions = c.GlobalInt("max-mentions")
	}

	if c.GlobalIsSet("retweets") {
		config.Retweets = c.GlobalBool("retweets")
	}

	if c.GlobalIsSet("replies") {
		config.Replies = c.GlobalBool("replies")
	}

	if c.GlobalIsSet("via") {
		config.Via = c.GlobalBool("via")
	}

	if c.GlobalIsSet("urls") {
		config.URLs = c.GlobalBool("urls")
	}

	if c.GlobalIsSet("like") {
		config.Like = c.GlobalBool("like")
	}

	if c.GlobalIsSet("smart-like") {
		config.SmartLike = c.GlobalBool("smart-like")
	}

	return nil
}
