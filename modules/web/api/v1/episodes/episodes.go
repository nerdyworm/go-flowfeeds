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

	episode, err := models.FindEpisodeByIdForUser(int64(id), ctx.User)
	if err != nil {
		return err
	}

	feed, err := models.FindFeedById(episode.FeedId)
	if err != nil {
		return err
	}

	serializers.JSON(w, serializers.NewShowEpisode(episode, feed))
	return nil
}

func Listens(ctx ctx.Context, w http.ResponseWriter, r *http.Request) error {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		return err
	}

	listens, users, err := models.FindListensForEpisode(int64(id))
	if err != nil {
		return err
	}

	serializers.JSON(w, serializers.NewListens(listens, users))
	return nil
}

func Favorites(ctx ctx.Context, w http.ResponseWriter, r *http.Request) error {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		return err
	}

	favorites, users, err := models.FindFavoritesForEpisode(int64(id))
	if err != nil {
		return err
	}

	serializers.JSON(w, serializers.NewFavorites(favorites, users))
	return nil
}

func Related(ctx ctx.Context, w http.ResponseWriter, r *http.Request) error {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		return err
	}

	related, err := models.FindRelatedEpisodes(int64(id))
	if err != nil {
		return err
	}

	serializers.JSON(w, serializers.NewEpisodes(related))
	return nil
}

func ToggleFavorite(ctx ctx.Context, w http.ResponseWriter, r *http.Request) error {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		return err
	}

	err = models.ToggleFavorite(ctx.User, int64(id))
	if err != nil {
		return err
	}

	episode, err := models.FindEpisodeByIdForUser(int64(id), ctx.User)
	if err != nil {
		return err
	}

	feed, err := models.FindFeedById(episode.FeedId)
	if err != nil {
		return err
	}

	return serializers.JSON(w, serializers.NewShowEpisode(episode, feed))
}
