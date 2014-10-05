package featured

import (
	"log"
	"net/http"

	"bitbucket.org/nerdyworm/go-flowfeeds/models"
	"bitbucket.org/nerdyworm/go-flowfeeds/modules/web/api/v1/serializers"
	"bitbucket.org/nerdyworm/go-flowfeeds/modules/web/ctx"
)

func Index(ctx ctx.Context, w http.ResponseWriter, r *http.Request) error {
	serializer := serializers.FeaturedsSerializer{}
	serializer.Featureds = make([]models.Featured, 0)
	serializer.Episodes = make([]serializers.Episode, 0)
	serializer.Feeds = make([]serializers.Feed, 0)

	teasers, feeds, listens, favorites, err := models.FeaturedEpisodes(ctx.User)
	if err != nil {
		log.Println("featued.Index models.FeaturedEpisodeEpisodes", err)
		return err
	}

	for i, teaser := range teasers {
		serializer.Featureds = append(serializer.Featureds, models.Featured{Rank: i, Episode: teaser.Id})
		serializer.Episodes = append(serializer.Episodes, serializers.NewEpisode(teaser))
	}

	for _, feed := range feeds {
		serializer.Feeds = append(serializer.Feeds, serializers.NewFeed(feed))
	}

	// XXX - maybe move this into the model layer...
	episodesToListens := make(map[int64]bool)
	for _, listen := range listens {
		if _, ok := episodesToListens[listen.EpisodeId]; !ok {
			episodesToListens[listen.EpisodeId] = true
		}
	}

	episodesToFavorites := make(map[int64]bool)
	for _, favorite := range favorites {
		if _, ok := episodesToFavorites[favorite.EpisodeId]; !ok {
			episodesToFavorites[favorite.EpisodeId] = true
		}
	}

	for i, t := range serializer.Episodes {
		if listened, ok := episodesToListens[t.Id]; ok {
			serializer.Episodes[i].Listened = listened
		}

		if favorited, ok := episodesToFavorites[t.Id]; ok {
			serializer.Episodes[i].Favorited = favorited
		}
	}

	serializers.JSON(w, serializer)
	return nil
}
