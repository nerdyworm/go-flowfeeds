package favorites

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

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
		log.Println("listens.Create strconv.Atoi", err)
		return err
	}

	listen, err := models.CreateFavorite(ctx.User, int64(id))
	if err != nil {
		log.Println("listens.Create models.CreateFavorite", err)
		return err
	}

	return serializers.JSON(w, serializers.NewShowFavorite(listen))
}
