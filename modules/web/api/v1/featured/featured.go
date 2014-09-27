package featured

import (
	"fmt"
	"log"
	"net/http"

	"bitbucket.org/nerdyworm/go-flowfeeds/models"
	"bitbucket.org/nerdyworm/go-flowfeeds/modules/web/api/v1/serializers"
	"bitbucket.org/nerdyworm/go-flowfeeds/modules/web/ctx"
)

func Index(ctx ctx.Context, w http.ResponseWriter, r *http.Request) error {
	serializer := serializers.FeaturedsSerializer{}
	serializer.Featureds = make([]models.Featured, 0)
	serializer.Teasers = make([]serializers.Teaser, 0)
	serializer.Feeds = make([]serializers.Feed, 0)
	serializer.Listens = make([]serializers.Listen, 0)
	serializer.Favorites = make([]serializers.Favorite, 0)

	teasers, feeds, listens, favorites, err := models.FeaturedEpisodeTeasers(ctx.User)
	if err != nil {
		log.Println("featued.Index models.FeaturedEpisodeTeasers", err)
		return err
	}

	for i, teaser := range teasers {
		serializer.Featureds = append(serializer.Featureds, models.Featured{Rank: i, Teaser: teaser.Id})
		serializer.Teasers = append(serializer.Teasers, serializers.NewTeaser(teaser))
	}

	for _, feed := range feeds {
		serializer.Feeds = append(serializer.Feeds, serializers.Feed{
			Id:          feed.Id,
			Title:       feed.Title,
			Description: feed.Description,
			Url:         feed.Url,
			Thumb:       fmt.Sprintf("http://s3.amazonaws.com/flowfeeds2/feeds/%d/thumb.jpg", feed.Id),
			Cover:       fmt.Sprintf("http://s3.amazonaws.com/flowfeeds2/feeds/%d/cover.jpg", feed.Id),
		})
	}

	episodesToListens := make(map[int64][]models.Listen)
	for _, listen := range listens {
		serializer.Listens = append(serializer.Listens, serializers.Listen{
			Id:      listen.Id,
			User:    listen.UserId,
			Episode: listen.EpisodeId,
		})

		if _, ok := episodesToListens[listen.EpisodeId]; !ok {
			episodesToListens[listen.EpisodeId] = []models.Listen{}
		}

		episodesToListens[listen.EpisodeId] = append(episodesToListens[listen.EpisodeId], listen)
	}

	episodesToFavorites := make(map[int64][]models.Favorite)
	for _, favorite := range favorites {
		serializer.Favorites = append(serializer.Favorites, serializers.Favorite{
			Id:      favorite.Id,
			User:    favorite.UserId,
			Episode: favorite.EpisodeId,
		})

		if _, ok := episodesToFavorites[favorite.EpisodeId]; !ok {
			episodesToFavorites[favorite.EpisodeId] = []models.Favorite{}
		}

		episodesToFavorites[favorite.EpisodeId] = append(episodesToFavorites[favorite.EpisodeId], favorite)
	}

	for i, t := range serializer.Teasers {
		ids := make([]int64, 0)
		if listens, ok := episodesToListens[t.Id]; ok {
			for _, l := range listens {
				ids = append(ids, l.Id)
			}
		}
		serializer.Teasers[i].Listens = ids

		favIds := make([]int64, 0)
		if favorites, ok := episodesToFavorites[t.Id]; ok {
			for _, l := range favorites {
				favIds = append(favIds, l.Id)
			}
		}
		serializer.Teasers[i].Favorites = favIds
	}

	serializers.JSON(w, serializer)
	return nil
}
