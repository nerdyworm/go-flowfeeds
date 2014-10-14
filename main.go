package main

import (
	"log"
	"os"

	"bitbucket.org/nerdyworm/go-flowfeeds/config"
	"bitbucket.org/nerdyworm/go-flowfeeds/datastore"
	"bitbucket.org/nerdyworm/go-flowfeeds/modules/faker"
	"bitbucket.org/nerdyworm/go-flowfeeds/modules/updates"
	"bitbucket.org/nerdyworm/go-flowfeeds/modules/web"
	"github.com/codegangsta/cli"
)

func main() {
	err := datastore.Connect(config.PG_CONFIG)
	if err != nil {
		log.Fatalln(err)
	}

	app := cli.NewApp()
	app.Name = "Flowfeeds"
	app.Usage = "http server and utils"
	app.Author = "Benjamin Silas Rhodes"
	app.Email = "ben@nerdyworm.com"
	app.Version = "0.0.1"
	app.Commands = []cli.Command{
		cli.Command{
			Name:        "update",
			Description: "This command will update all the rss feeds",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "file",
					Value: "db/collections.json",
					Usage: "rss feeds collection file",
				},
			},
			Action: func(c *cli.Context) {
				updates.Run(c.String("file"))
			},
		},
		cli.Command{
			Name:        "faker",
			Description: "This command will generate fake users and activity - shhhhhh",
			Action: func(c *cli.Context) {
				faker.Run()
			},
		},
		cli.Command{
			Name:        "server",
			Description: "This command will start the http server",
			Action: func(c *cli.Context) {
				web.Run()
			},
		},
	}
	app.Run(os.Args)
}
