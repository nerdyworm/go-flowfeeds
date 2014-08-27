package episodes

import (
	"fmt"
	"net/http"
	"strconv"

	"bitbucket.org/nerdyworm/go-flowfeeds/models"
	"bitbucket.org/nerdyworm/go-flowfeeds/modules/web/api/v1/serializers"
	"github.com/gorilla/mux"
)

func Index(w http.ResponseWriter, r *http.Request) {
	p := serializers.EpisodesSerializer{}
	p.Episodes = make([]models.Episode, 0)

	for i := 0; i < 50; i++ {
		p.Episodes = append(p.Episodes, models.Episode{
			Id:          int64(i + 1),
			Title:       fmt.Sprintf("Title %d", i+1),
			Description: "Description",
			Url:         "http://example.com/id.mp3",
		})
	}

	serializers.JSON(w, p)
}

func Show(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		panic(err)
	}

	episode, err := models.FindEpisodeById(int64(id))
	if err != nil {
		panic(err)
	}

	p := serializers.EpisodeSerializer{}
	p.Episode.Id = episode.Id
	p.Episode.Url = episode.Url
	p.Episode.Title = episode.Title
	p.Episode.Description = episode.Description
	p.Episode.Comments = make([]int64, 10)
	p.Comments = make([]models.Comment, 10)

	for i := 0; i < 10; i++ {
		id := int64(i + 1)
		p.Episode.Comments[i] = id
		p.Comments[i] = models.Comment{
			Id:   id,
			Body: "this just a test",
		}
	}

	serializers.JSON(w, p)
}
