package main

import (
	"log"
	"os"

	"github.com/nerdyworm/go-flowfeeds/cmds/fake"
	"github.com/nerdyworm/go-flowfeeds/cmds/update"
	"github.com/nerdyworm/go-flowfeeds/config"
	"github.com/nerdyworm/go-flowfeeds/datastore"
	"github.com/nerdyworm/go-flowfeeds/web"
	"github.com/codegangsta/cli"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

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
				update.Rss()
				update.Image()
			},
		},
		cli.Command{
			Name:        "fake",
			Description: "This command will generate fake users and activity - shhhhhh",
			Action: func(c *cli.Context) {
				fake.Run()
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
