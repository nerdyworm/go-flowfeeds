package episodes

import (
	"net/http"
	"strconv"

	"bitbucket.org/nerdyworm/go-flowfeeds/models"
	"bitbucket.org/nerdyworm/go-flowfeeds/modules/web/api/v1/serializers"
	"bitbucket.org/nerdyworm/go-flowfeeds/modules/web/ctx"
	"github.com/gorilla/mux"
)

func Index(ctx ctx.Context, w http.ResponseWriter, r *http.Request) error {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	options := models.ListOptions{
		PerPage: 24,
		Page:    page,
	}

	episodes, feeds, err := ctx.Store.Episodes.ListFor(&ctx.User, options)
	if err != nil {
		return err
	}

	serializers.JSON(w, serializers.NewEpisodes(episodes, feeds))
	return nil
}

func Show(ctx ctx.Context, w http.ResponseWriter, r *http.Request) error {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		return err
	}

	episode, err := ctx.Store.Episodes.GetForUser(&ctx.User, int64(id))
	if err != nil {
		return err
	}

	feed, err := ctx.Store.Feeds.Get(episode.FeedId)
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

	related, err := ctx.Store.Episodes.Related(int64(id))
	if err != nil {
		return err
	}

	serializers.JSON(w, serializers.NewEpisodes(related, []*models.Feed{}))
	return nil
}

func ToggleFavorite(ctx ctx.Context, w http.ResponseWriter, r *http.Request) error {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		return err
	}

	err = ctx.Store.Episodes.ToggleFavoriteForUser(&ctx.User, int64(id))
	if err != nil {
		return err
	}

	episode, err := ctx.Store.Episodes.GetForUser(&ctx.User, int64(id))
	if err != nil {
		return err
	}

	feed, err := ctx.Store.Feeds.Get(episode.FeedId)
	if err != nil {
		return err
	}

	return serializers.JSON(w, serializers.NewShowEpisode(episode, feed))
}
