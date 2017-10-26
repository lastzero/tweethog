package main

import (
	"github.com/lastzero/tweethog"
	"github.com/urfave/cli"
	"os"
	"fmt"
	"strings"
)

func main() {
	config := tweethog.NewConfig()

	app := cli.NewApp()
	app.Usage = "Stream, filter and react to Twitter status updates"
	app.Version = "0.6.0"
	app.Copyright = "Michael Mayer <michael@liquidbytes.net>"

	app.Flags = globalCliFlags

	app.Commands = []cli.Command{
		{
			Name:  "config",
			Usage: "Displays all configuration values",
			Flags: cliFlags,
			Action: func(c *cli.Context) {
				config.SetValuesFromFile(c.GlobalString("config-file"))

				config.SetValuesFromCliContext(c)

				fmt.Printf("Topics        : %s\n", strings.Join(config.Filter.Topics, ", "))
				fmt.Printf("Languages     : %s\n", strings.Join(config.Filter.Languages, ", "))
				fmt.Printf("Min followers : %d\n", config.Filter.MinFollowers)
				fmt.Printf("Max followers : %d\n", config.Filter.MaxFollowers)
				fmt.Printf("Min following : %d\n", config.Filter.MinFollowing)
				fmt.Printf("Max following : %d\n", config.Filter.MaxFollowing)
				fmt.Printf("Max mentions  : %d\n", config.Filter.MaxMentions)
				fmt.Printf("Max tags      : %d\n", config.Filter.MaxTags)
				fmt.Printf("URLs          : %t\n", config.Filter.URLs)
				fmt.Printf("Retweets      : %t\n", config.Filter.Retweets)
				fmt.Printf("Replies       : %t\n", config.Filter.Replies)
				fmt.Printf("Via           : %t\n", config.Filter.Via)
			},
		},
		{
			Name:  "filter",
			Usage: "Shows all matching tweets without performing any action",
			Flags: cliFlags,
			Action: func(c *cli.Context) {
				startStream(
					c,
					config,
					func(status *tweethog.Status) {
						if status.MatchesFilter(config.Filter) {
							printTweet(status)
						}
					},
				)
			},
		},
		{
			Name:  "like",
			Usage: "Automatically likes all matching tweets",
			Flags: cliFlags,
			Action: func(c *cli.Context) {
				startStream(
					c,
					config,
					func(status *tweethog.Status) {
						if status.MatchesFilter(config.Filter) {
							printTweet(status)

							status.Like()
						}
					},
				)
			},
		},
		{
			Name:  "smartlike",
			Usage: "Likes tweets with random delay and rate limit",
			Flags: cliFlags,
			Action: func(c *cli.Context) {
				startStream(
					c,
					config,
					func(status *tweethog.Status) {
						if status.MatchesFilter(config.Filter) {
							printTweet(status)

							go status.SmartLike()
						}
					},
				)
			},
		},
	}

	app.Run(os.Args)
}

func printTweet(status *tweethog.Status) {
	fmt.Printf("\n%s @%s (Following: %d, Followers: %d, Likes: %d)\n%s\n",
		status.GetName(),
		status.GetScreenName(),
		status.GetFriendsCount(),
		status.GetFollowersCount(),
		status.GetFavouritesCount(),
		status.GetText(),
	)
}

func startStream(c *cli.Context, config *tweethog.Config, action func(status *tweethog.Status)) error {
	err := config.SetValuesFromFile(c.GlobalString("config-file"))

	if err != nil {
		return cli.NewExitError(err, 1)
	}

	config.SetValuesFromCliContext(c)

	if len(config.Filter.Topics) == 0 {
		cli.ShowAppHelp(c)
		return nil
	}

	stream := tweethog.NewStream(config)

	stream.Start(action)

	return nil
}

var globalCliFlags = []cli.Flag{
	cli.StringFlag{
		Name:  "config-file, c",
		Usage: "YAML config filename",
		Value: "config.yml",
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
		Usage: "Twitter API access token",
	},
	cli.StringFlag{
		Name:  "access-secret",
		Usage: "Twitter API access token secret",
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
}