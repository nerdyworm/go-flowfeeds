package favorites

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"bitbucket.org/nerdyworm/go-flowfeeds/models"
	"bitbucket.org/nerdyworm/go-flowfeeds/modules/web/api/v1/serializers"
	"bitbucket.org/nerdyworm/go-flowfeeds/modules/web/ctx"
)

type CreateFavoriteRequest struct {
	Favorite struct {
		Episode string
		Teaser  string
	}
}

func Create(ctx ctx.Context, w http.ResponseWriter, r *http.Request) error {
	if ctx.User.Id == 0 {
		return nil
	}

	request := CreateFavoriteRequest{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		log.Println("listens.Create json.Decode", err)
		return err
	}

	id, err := strconv.Atoi(request.Favorite.Episode)
	if err != nil {
		id, err = strconv.Atoi(request.Favorite.Teaser)
		if err != nil {
			return err
		}
	}

	listen, err := models.CreateFavorite(ctx.User, int64(id))
	if err != nil {
		log.Println("listens.Create models.CreateFavorite", err)
		return err
	}

	episode, err := models.FindEpisodeById(listen.EpisodeId)
	if err != nil {
		log.Println("listens.Create models.FindEpisodeById", err)
		return err
	}

	return serializers.JSON(w, serializers.NewShowFavorite(listen, episode))
}

func Delete(ctx ctx.Context, w http.ResponseWriter, r *http.Request) error {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		return err
	}

	favorite, err := models.FindFavoriteById(int64(id))
	if err != nil {
		return err
	}

	err = models.DeleteFavorite(ctx.User, favorite.Id)
	if err != nil {
		return err
	}

	episode, err := models.FindEpisodeById(favorite.EpisodeId)
	if err != nil {
		return err
	}

	return serializers.JSON(w, serializers.NewShowFavorite(favorite, episode))
}
