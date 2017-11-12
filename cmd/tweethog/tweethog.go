package main

import (
	"github.com/lastzero/tweethog"
	"github.com/urfave/cli"
	"os"
	"fmt"
	"strings"
	"log"
)

func main() {
	config := tweethog.NewConfig()

	app := cli.NewApp()
	app.Usage = "Stream, filter and react to Twitter status updates"
	app.Version = "0.8.0"
	app.Copyright = "Michael Mayer <michael@liquidbytes.net>"

	app.Flags = globalCliFlags

	app.Commands = []cli.Command{
		{
			Name:  "config",
			Usage: "Displays all configuration values",
			Flags: cliFlags,
			Action: func(c *cli.Context) {
				config.SetValuesFromFile(tweethog.GetExpandedFilename(c.GlobalString("config-file")))

				config.SetValuesFromCliContext(c)

				fmt.Printf("Name              | Value\n")
				fmt.Printf("------------------|--------------------------------------------------------\n")
				fmt.Printf("config-file       | %s\n", config.ConfigFile)
				fmt.Printf("consumer-key      | %s\n", config.ConsumerKey)
				fmt.Printf("consumer-secret   | %s\n", config.ConsumerSecret)
				fmt.Printf("access-token      | %s\n", config.AccessToken)
				fmt.Printf("access-secret     | %s\n", config.AccessSecret)
				fmt.Printf("topic             | %s\n", strings.Join(config.Filter.Topics, ", "))
				fmt.Printf("lang              | %s\n", strings.Join(config.Filter.Languages, ", "))
				fmt.Printf("min-followers     | %d\n", config.Filter.MinFollowers)
				fmt.Printf("max-followers     | %d\n", config.Filter.MaxFollowers)
				fmt.Printf("min-following     | %d\n", config.Filter.MinFollowing)
				fmt.Printf("max-following     | %d\n", config.Filter.MaxFollowing)
				fmt.Printf("max-tags          | %d\n", config.Filter.MaxTags)
				fmt.Printf("max-mentions      | %d\n", config.Filter.MaxMentions)
				fmt.Printf("retweets          | %t\n", config.Filter.Retweets)
				fmt.Printf("replies           | %t\n", config.Filter.Replies)
				fmt.Printf("via               | %t\n", config.Filter.Via)
				fmt.Printf("urls              | %t\n", config.Filter.URLs)
				fmt.Printf("images-only       | %t\n", config.Filter.ImagesOnly)
				fmt.Printf("save-images       | %s\n", config.SaveImages)
				fmt.Printf("json-log          | %s\n", config.JsonLog)
			},
		},
		{
			Name:  "auth",
			Usage: "Requests a user access token and secret for the Twitter API",
			Flags: []cli.Flag{},
			Action: func(c *cli.Context) {
				config.SetValuesFromFile(tweethog.GetExpandedFilename(c.GlobalString("config-file")))

				config.SetValuesFromCliContext(c)

				tweethog.CliAuth(config.ConsumerKey, config.ConsumerSecret)
			},
		},
		{
			Name:  "filter",
			Usage: "Shows all matching tweets without performing any action",
			Flags: cliFlags,
			Action: func(c *cli.Context) error {
				return startStream(
					c,
					config,
					func(status *tweethog.Status) {
						if status.MatchesFilter(config.Filter) {
							handleTweet(status, config)
						}
					},
				)
			},
		},
		{
			Name:  "like",
			Usage: "Automatically likes all matching tweets",
			Flags: cliFlags,
			Action: func(c *cli.Context) error {
				return startStream(
					c,
					config,
					func(status *tweethog.Status) {
						if status.MatchesFilter(config.Filter) {
							handleTweet(status, config)

							status.Like()

							log.Printf("Liked status ❤️\n")

						}
					},
				)
			},
		},
		{
			Name:  "follow",
			Usage: "Automatically follows all users with matching tweets",
			Flags: cliFlags,
			Action: func(c *cli.Context) error {
				return startStream(
					c,
					config,
					func(status *tweethog.Status) {
						if status.MatchesFilter(config.Filter) {
							handleTweet(status, config)

							status.Follow()

							log.Printf("Followed user @%s 🐷\n", status.GetScreenName())
						}
					},
				)
			},
		},
		{
			Name:  "smartlike",
			Usage: "Likes tweets with random delay and rate limit",
			Flags: cliFlags,
			Action: func(c *cli.Context) error {
				return startStream(
					c,
					config,
					func(status *tweethog.Status) {
						if status.MatchesFilter(config.Filter) {
							handleTweet(status, config)

							go status.SmartLike()
						}
					},
				)
			},
		},
	}

	app.Run(os.Args)
}

func handleTweet(status *tweethog.Status, config *tweethog.Config) {
	fmt.Printf("\n%s %s @%s (Following: %d, Followers: %d, Likes: %d)\n%s\n",
		status.GetCreatedAt().Local().Format(tweethog.CompactTime),
		status.GetName(),
		status.GetScreenName(),
		status.GetFriendsCount(),
		status.GetFollowersCount(),
		status.GetFavouritesCount(),
		status.GetText(),
	)

	if config.SaveImages != "" {
		if filename, err := status.SaveImageToFile(config.SaveImages); err != nil {
			log.Println(err)
		} else {
			log.Println("Saved image to " + filename + " 💾")
		}
	}

	if config.JsonLog != "" {
		encoded, err := status.GetAsJson()

		if err != nil {
			log.Println(err)
			return
		}

		if err := appendLineToLog(config.JsonLog, encoded); err != nil {
			log.Println(err)
			return
		}
	}
}

func appendLineToLog(path, text string) error {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModeAppend)

	if err != nil {
		return err
	}

	defer f.Close()

	os.Chmod(path, 0600)

	_, err = f.WriteString(text + "\n")

	if err != nil {
		return err
	}

	return nil
}

func startStream(c *cli.Context, config *tweethog.Config, action func(status *tweethog.Status)) error {
	err := config.SetValuesFromFile(tweethog.GetExpandedFilename(c.GlobalString("config-file")))

	if err != nil {
		return cli.NewExitError(err, 1)
	}

	config.SetValuesFromCliContext(c)

	if len(config.Filter.Topics) == 0 {
		return cli.NewExitError("At least one topic is required, use -h to show available filters", 1)
	}

	stream := tweethog.NewStream(config)

	stream.Start(action)

	return nil
}

var globalCliFlags = []cli.Flag{
	cli.StringFlag{
		Name:  "config-file, c",
		Usage: "Config filename",
		Value: "~/.tweethog",
	},
	cli.StringFlag{
		Name:  "consumer-key",
		Usage: "Twitter API consumer key",
	},
	cli.StringFlag{
		Name:  "consumer-secret",
		Usage: "Twitter API consumer secret",
	},
	cli.StringFlag{
		Name:  "access-token",
		Usage: "Twitter API user access token",
	},
	cli.StringFlag{
		Name:  "access-secret",
		Usage: "Twitter API user access token secret",
	},
}

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
		Name:  "min-followers",
		Usage: "User min followers",
	},
	cli.IntFlag{
		Name:  "max-followers",
		Usage: "User max followers, 0 for unlimited",
	},
	cli.IntFlag{
		Name:  "min-following",
		Usage: "User min following",
	},
	cli.IntFlag{
		Name:  "max-following",
		Usage: "User max following, 0 for unlimited",
	},
	cli.IntFlag{
		Name:  "max-tags",
		Usage: "Max number of hash #tags",
	},
	cli.IntFlag{
		Name:  "max-mentions",
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
		Name:  "images-only",
		Usage: "Only tweets containing images",
	},
	cli.StringFlag{
		Name:  "save-images",
		Usage: "Save all images in a directory",
	},
	cli.StringFlag{
		Name:  "json-log",
		Usage: "Log matching tweets in a file as newline delimited JSON",
	},
}