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
	JsonLog    string
	SaveImages string

	// Tweet filter
	Filter *Filters
}

type Filters struct {
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
	ImagesOnly   bool
}

func NewConfig() *Config {
	return &Config{
		Filter: &Filters{},
	}
}

func (config *Config) SetValuesFromFile(fileName string) error {
	yamlConfig, err := yaml.ReadFile(fileName)

	if err != nil {
		return err
	}

	config.ConfigFile = fileName

	config.ConsumerKey, _ = yamlConfig.Get("consumer-key")
	config.ConsumerSecret, _ = yamlConfig.Get("consumer-secret")
	config.AccessToken, _ = yamlConfig.Get("access-token")
	config.AccessSecret, _ = yamlConfig.Get("access-secret")

	if topics, err := yamlConfig.Get("topic"); err == nil {
		config.Filter.Topics = strings.Split(topics, ",")
	}

	if languages, err := yamlConfig.Get("lang"); err == nil {
		config.Filter.Languages = strings.Split(languages, ",")
	}

	if minFollowers, err := yamlConfig.GetInt("min-followers"); err == nil {
		config.Filter.MinFollowers = int(minFollowers)
	}

	if maxFollowers, err := yamlConfig.GetInt("max-followers"); err == nil {
		config.Filter.MaxFollowers = int(maxFollowers)
	}

	if minFollowing, err := yamlConfig.GetInt("min-following"); err == nil {
		config.Filter.MinFollowing = int(minFollowing)
	}

	if maxFollowing, err := yamlConfig.GetInt("max-following"); err == nil {
		config.Filter.MaxFollowing = int(maxFollowing)
	}

	if maxTags, err := yamlConfig.GetInt("max-tags"); err == nil {
		config.Filter.MaxTags = int(maxTags)
	}

	if maxMentions, err := yamlConfig.GetInt("max-mentions"); err == nil {
		config.Filter.MaxMentions = int(maxMentions)
	}

	if retweets, err := yamlConfig.GetBool("retweets"); err == nil {
		config.Filter.Retweets = retweets
	}

	if replies, err := yamlConfig.GetBool("replies"); err == nil {
		config.Filter.Replies = replies
	}

	if via, err := yamlConfig.GetBool("via"); err == nil {
		config.Filter.Via = via
	}

	if urls, err := yamlConfig.GetBool("urls"); err == nil {
		config.Filter.URLs = urls
	}

	if imagesOnly, err := yamlConfig.GetBool("images-only"); err == nil {
		config.Filter.ImagesOnly = imagesOnly
	}

	if saveImages, err := yamlConfig.Get("save-images"); err == nil {
		config.SaveImages = GetExpandedFilename(saveImages)
	}

	if jsonlog, err := yamlConfig.Get("json-log"); err == nil {
		config.JsonLog = GetExpandedFilename(jsonlog)
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

	if c.IsSet("topic") || c.IsSet("t") {
		config.Filter.Topics = c.StringSlice("topic")
	}

	if c.IsSet("lang") || c.IsSet("l") {
		config.Filter.Languages = c.StringSlice("lang")
	}

	if c.IsSet("min-followers") {
		config.Filter.MinFollowers = c.Int("min-followers")
	}

	if c.IsSet("max-followers") {
		config.Filter.MaxFollowers = c.Int("max-followers")
	}

	if c.IsSet("min-following") {
		config.Filter.MinFollowing = c.Int("min-following")
	}

	if c.IsSet("max-following") {
		config.Filter.MaxFollowing = c.Int("max-following")
	}

	if c.IsSet("max-tags") {
		config.Filter.MaxTags = c.Int("max-tags")
	}

	if c.IsSet("max-mentions") {
		config.Filter.MaxMentions = c.Int("max-mentions")
	}

	if c.IsSet("retweets") {
		config.Filter.Retweets = c.Bool("retweets")
	}

	if c.IsSet("replies") {
		config.Filter.Replies = c.Bool("replies")
	}

	if c.IsSet("via") {
		config.Filter.Via = c.Bool("via")
	}

	if c.IsSet("urls") {
		config.Filter.URLs = c.Bool("urls")
	}

	if c.IsSet("images-only") {
		config.Filter.ImagesOnly = c.Bool("images-only")
	}

	if c.IsSet("save-images") {
		config.SaveImages = GetExpandedFilename(c.String("save-images"))
	}

	if c.IsSet("json-log") {
		config.JsonLog = GetExpandedFilename(c.String("json-log"))
	}

	return nil
}
