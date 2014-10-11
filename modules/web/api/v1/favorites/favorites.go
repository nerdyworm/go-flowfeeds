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
		return err
	}

	err = models.ToggleFavorite(ctx.User, int64(id))
	if err != nil {
		log.Println("listens.Create models.ToggleFavorite", err)
		return err
	}

	episode, err := ctx.Store.Episodes.GetForUser(&ctx.User, int64(id))
	if err != nil {
		log.Println("listens.Create models.FindEpisodeById", err)
		return err
	}

	return serializers.JSON(w, serializers.NewShowFavorite(models.Favorite{}, *episode))
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

	episode, err := ctx.Store.Episodes.GetForUser(&ctx.User, int64(id))
	if err != nil {
		return err
	}

	return serializers.JSON(w, serializers.NewDeleteFavorite(favorite, *episode))
}
