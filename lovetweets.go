package main

import (
	"fmt"
	"github.com/urfave/cli"
	"github.com/ChimeraCoder/anaconda"
	"os"
)

func main()  {
	app := cli.NewApp()
	app.Usage = "Find tweets you love"
	app.Version = "0.0.1"

	app.Action = func (c *cli.Context) error {
		anaconda.SetConsumerKey("your-consumer-key")
		anaconda.SetConsumerSecret("your-consumer-secret")
		api := anaconda.NewTwitterApi("your-access-token", "your-access-token-secret")
		searchResult, err := api.GetSearch("golang", nil)

		if err != nil {
			return cli.NewExitError(err, 1)
		}

		for _, tweet := range searchResult.Statuses {
			fmt.Println(tweet.Text)
		}

		return nil
	}

	app.Run(os.Args)
}
