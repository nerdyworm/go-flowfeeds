package episodes

import (
	"fmt"
	"net/http"
	"strconv"

	"bitbucket.org/nerdyworm/go-flowfeeds/models"
	"bitbucket.org/nerdyworm/go-flowfeeds/modules/web/api/v1/serializers"
	"bitbucket.org/nerdyworm/go-flowfeeds/modules/web/ctx"
	"github.com/gorilla/mux"
)

func Index(ctx ctx.Context, w http.ResponseWriter, r *http.Request) error {
	p := serializers.Episodes{}
	p.Episodes = make([]serializers.Episode, 0)

	for i := 0; i < 50; i++ {
		p.Episodes = append(p.Episodes, serializers.Episode{
			Id:          int64(i + 1),
			Title:       fmt.Sprintf("Title %d", i+1),
			Description: "Description",
			Url:         "http://example.com/id.mp3",
		})
	}

	serializers.JSON(w, p)
	return nil
}

func Show(ctx ctx.Context, w http.ResponseWriter, r *http.Request) error {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		return err
	}

	episode, err := models.FindEpisodeById(int64(id))
	if err != nil {
		return err
	}

	serializers.JSON(w, serializers.NewShowEpisode(episode))
	return nil
}

func Listens(ctx ctx.Context, w http.ResponseWriter, r *http.Request) error {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		return err
	}

	listens, err := models.FindListensForEpisode(int64(id))
	if err != nil {
		return err
	}

	serializers.JSON(w, serializers.NewListens(listens))
	return nil
}

func Favorites(ctx ctx.Context, w http.ResponseWriter, r *http.Request) error {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		return err
	}

	listens, err := models.FindFavoritesForEpisode(int64(id))
	if err != nil {
		return err
	}

	serializers.JSON(w, serializers.NewFavorites(listens))
	return nil
}

func Related(ctx ctx.Context, w http.ResponseWriter, r *http.Request) error {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		return err
	}

	related, err := models.FindRelatedTeasers(int64(id))
	if err != nil {
		return err
	}

	serializers.JSON(w, serializers.NewTeasers(related))
	return nil
}
