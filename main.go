package main

import (
	"log"
	"os"

	"bitbucket.org/nerdyworm/go-flowfeeds/models"
	"bitbucket.org/nerdyworm/go-flowfeeds/modules/updates"
	"bitbucket.org/nerdyworm/go-flowfeeds/modules/web"
	"github.com/codegangsta/cli"
)

var ()

func main() {
	err := models.Connect("dbname=flowfeeds2 sslmode=disable")
	if err != nil {
		log.Fatalln(err)
	}
	defer models.Close()

	app := cli.NewApp()
	app.Name = "Flowfeeds"
	app.Usage = "Flowfeeds Service"
	app.Version = "0.0.1"
	app.Commands = []cli.Command{
		cli.Command{
			Name:        "update",
			Description: "This command will update all the rss feeds",
			Action: func(c *cli.Context) {
				updates.Run(c.String("file"))
			},
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "file",
					Value: "db/collections.yml",
					Usage: "rss feeds collection file",
				},
			},
		},
		cli.Command{
			Name:        "server",
			Description: "This command will start the http server",
			Action: func(c *cli.Context) {
				web.Run(web.ServerOptions{
					c.String("ember"),
				})
			},
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "ember",
					Usage:  "tells the http server where it can find the index.html file",
					EnvVar: "EMBER_APP_PATH",
				},
			},
		},
	}
	app.Run(os.Args)
}
