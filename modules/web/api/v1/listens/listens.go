package listens

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"bitbucket.org/nerdyworm/go-flowfeeds/models"
	"bitbucket.org/nerdyworm/go-flowfeeds/modules/web/api/v1/serializers"
	"bitbucket.org/nerdyworm/go-flowfeeds/modules/web/ctx"
)

type CreateListenRequest struct {
	Listen struct {
		Episode string
	}
}

func Create(ctx ctx.Context, w http.ResponseWriter, r *http.Request) error {
	if ctx.User.Id == 0 {
		return nil
	}

	request := CreateListenRequest{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		log.Println("listens.Create json.Decode", err)
		return err
	}

	id, err := strconv.Atoi(request.Listen.Episode)
	if err != nil {
		return err
	}

	listen, err := models.CreateListen(ctx.User, int64(id))
	if err != nil {
		log.Println("listens.Create models.CreateListen", err)
		return err
	}

	episode, err := models.FindEpisodeByIdForUser(listen.EpisodeId, ctx.User)
	if err != nil {
		log.Println("listens.Create models.FindEpisodeById", err)
		return err
	}

	return serializers.JSON(w, serializers.NewShowListen(listen, episode))
}
