package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"bitbucket.org/nerdyworm/go-flowfeeds/config"
	"bitbucket.org/nerdyworm/go-flowfeeds/models"
	"bitbucket.org/nerdyworm/go-flowfeeds/modules/web/api/v1/serializers"
)

type HomeController struct {
	ApplicationController
}

type manifest struct {
	CurrentUser int64
	Payload     struct {
		User serializers.User
	}
}

func newManifest(user models.User) manifest {
	m := manifest{CurrentUser: user.Id}
	m.Payload.User = serializers.NewUser(user)
	return m
}

func (c *HomeController) Index() error {
	c.ResponseWriter.Header().Add("Content-Type", "text/html")

	manifest, err := c.manifest()
	if err != nil {
		return err
	}

	html, err := c.html()
	if err != nil {
		return err
	}

	fmt.Fprint(c.ResponseWriter, strings.Replace(html, "</head>", manifest, 1))
	return nil
}

func (c *HomeController) manifest() (string, error) {
	if c.CurrentUser.Id != 0 {
		m := newManifest(c.CurrentUser)
		b, err := json.Marshal(m)
		if err != nil {
			return "", err
		}

		return fmt.Sprintf("<script>window.FlowfeedsManifest = %s;</script>\n</head>", b), nil
	}

	return "", nil
}

func (c *HomeController) html() (string, error) {
	index, err := ioutil.ReadFile(config.EMBER_APP + "/index.html")
	if err != nil {
		return "", err
	}

	return string(index), nil
}
