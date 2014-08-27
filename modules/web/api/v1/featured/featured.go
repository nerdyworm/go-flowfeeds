package featured

import (
	"log"
	"net/http"

	"bitbucket.org/nerdyworm/go-flowfeeds/models"
	"bitbucket.org/nerdyworm/go-flowfeeds/modules/web/api/v1/serializers"
)

func Index(w http.ResponseWriter, r *http.Request) {
	serializer := serializers.FeaturedsSerializer{}
	serializer.Featureds = make([]models.Featured, 0)
	serializer.Teasers = make([]models.Teaser, 0)

	teasers, err := models.FeaturedEpisodeTeasers()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	for i, teaser := range teasers {
		serializer.Featureds = append(serializer.Featureds, models.Featured{Rank: i, Teaser: teaser.Id})
		serializer.Teasers = append(serializer.Teasers, teaser)
	}

	serializers.JSON(w, serializer)
}
