package main

import (
	"github.com/urfave/cli"
	"os"
	"github.com/lastzero/tweethog"
)

func main() {
	app := cli.NewApp()
	app.Usage = "Stream, filter and react to Twitter status updates"
	app.Version = "0.5.0"
	app.Copyright = "Michael Mayer <michael@liquidbytes.net>"

	app.Flags = cliFlags

	app.Action = func (c *cli.Context) error {
		if len(c.GlobalStringSlice("topic")) == 0 {
			cli.ShowAppHelp(c)
			return nil
		}

		config := tweethog.NewConfig()

		err := config.LoadFromFile(c.GlobalString("config"))

		if err != nil {
			return cli.NewExitError(err, 1)
		}

		client := tweethog.NewClient(config)

		client.StreamTopic(c)

		return nil
	}

	app.Run(os.Args)
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
		Usage: "Likes tweets with GetRandomInt delay and rate limit",
	},
	cli.StringFlag{
		Name:  "config, c",
		Usage: "Config file name",
		Value: "config.yml",
	},
}