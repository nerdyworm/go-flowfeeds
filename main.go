package main

import (
	"log"
	"os"

	"bitbucket.org/nerdyworm/go-flowfeeds/config"
	"bitbucket.org/nerdyworm/go-flowfeeds/models"
	"bitbucket.org/nerdyworm/go-flowfeeds/modules/updates"
	"bitbucket.org/nerdyworm/go-flowfeeds/modules/web"
	"github.com/codegangsta/cli"
)

func main() {
	pgConfig := os.Getenv("DATABASE_CONFIG")

	if pgConfig == "" {
		pgConfig = config.PgConfig
	}

	err := models.Connect(pgConfig)
	if err != nil {
		log.Fatalln(err)
	}
	defer models.Close()

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
			Name:        "server",
			Description: "This command will start the http server",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "ember",
					Usage:  "tells the http server where it can find the index.html file",
					EnvVar: "EMBER_APP_PATH",
				},
			},
			Action: func(c *cli.Context) {
				web.Run(web.ServerOptions{
					c.String("ember"),
				})
			},
		},
	}
	app.Run(os.Args)
}
