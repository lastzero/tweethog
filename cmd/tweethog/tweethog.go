package main

import (
	"github.com/lastzero/tweethog"
	"github.com/urfave/cli"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Usage = "Stream, filter and react to Twitter status updates"
	app.Version = "0.5.0"
	app.Copyright = "Michael Mayer <michael@liquidbytes.net>"

	app.Flags = cliFlags

	app.Action = func(c *cli.Context) error {
		if len(c.GlobalStringSlice("topic")) == 0 {
			cli.ShowAppHelp(c)
			return nil
		}

		config := tweethog.NewConfig()

		err := config.SetValuesFromFile(c.GlobalString("config-file"))

		if err != nil {
			return cli.NewExitError(err, 1)
		}

		config.SetValuesFromCliContext(c)

		stream := tweethog.NewStream(config)

		stream.Start()

		return nil
	}

	app.Run(os.Args)
}

var cliFlags = []cli.Flag{
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
		Name:  "like",
		Usage: "Like all matching tweets",
	},
	cli.BoolFlag{
		Name:  "smart-like",
		Usage: "Likes tweets with GetRandomInt delay and rate limit",
	},
	cli.StringFlag{
		Name:  "config-file, c",
		Usage: "YAML config filename",
		Value: "config.yml",
	},
}
