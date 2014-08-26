package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"bitbucket.org/nerdyworm/go-flowfeeds/models"

	"os"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

var (
	RootEmberAppPath   = os.Getenv("ROOT_EMBER_APP_PATH")
	rootEmberIndexHtml = ""
)

func main() {
	err := models.Connect("dbname=flowfeeds2 sslmode=disable")
	if err != nil {
		log.Fatalln(err)
	}
	defer models.Close()

	log.Printf("Starting server...\n")
	log.Printf("RootEmberAppPath: %s\n", RootEmberAppPath)

	index, err := ioutil.ReadFile(RootEmberAppPath + "/index.html")
	if err != nil {
		log.Printf("could not find the ember app's index.html in `%s`", RootEmberAppPath)
		panic(err)
	}

	rootEmberIndexHtml = string(index)

	r := mux.NewRouter()
	r.HandleFunc("/", HomeHandler)
	r.NotFoundHandler = http.HandlerFunc(HomeHandler)

	r.PathPrefix("/assets").Handler(http.FileServer(http.Dir(RootEmberAppPath)))

	apiRouter := r.PathPrefix("/api/v1").Subrouter()
	apiRouter.HandleFunc("/episodes", EpisodesIndexHandler).Methods("GET")
	apiRouter.HandleFunc("/episodes/{id}", EpisodesShowHandler).Methods("GET")
	apiRouter.HandleFunc("/featureds", FeaturedsIndexHandler).Methods("GET")
	apiRouter.HandleFunc("/feeds", FeedsIndexHandler).Methods("GET")
	apiRouter.HandleFunc("/feeds/{id}", FeedShowHandler).Methods("GET")
	http.Handle("/", r)

	if err := http.ListenAndServe(":3000", nil); err != nil {
		panic(err)
	}
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html")
	fmt.Fprint(w, rootEmberIndexHtml)
}

type EpisodesSerializer struct {
	Episodes []models.Episode
}

type EpisodeSerializer struct {
	Episode struct {
		models.Episode
		Comments []int64
	}
	Comments []models.Comment
}

func EpisodesIndexHandler(w http.ResponseWriter, r *http.Request) {
	p := EpisodesSerializer{}
	p.Episodes = make([]models.Episode, 0)

	for i := 0; i < 50; i++ {
		p.Episodes = append(p.Episodes, models.Episode{
			Id:          int64(i + 1),
			Title:       fmt.Sprintf("Title %d", i+1),
			Description: "Description",
			Url:         "http://example.com/id.mp3",
		})
	}

	writeJSON(w, p)
}

func EpisodesShowHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		panic(err)
	}

	episode, err := models.FindEpisodeById(int64(id))
	if err != nil {
		panic(err)
	}

	p := EpisodeSerializer{}
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

	writeJSON(w, p)
}

func writeJSON(w http.ResponseWriter, v interface{}) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		panic(err)
	}

	w.Header().Set("content-type", "application/json; charset=utf-8")
	_, err = w.Write(data)
	return err
}

type FeaturedsSerializer struct {
	Featureds []models.Featured
	Teasers   []models.Teaser
}

func FeaturedsIndexHandler(w http.ResponseWriter, r *http.Request) {
	p := FeaturedsSerializer{}
	p.Featureds = make([]models.Featured, 0)
	p.Teasers = make([]models.Teaser, 0)

	teasers, err := models.FeaturedEpisodeTeasers()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	for i, teaser := range teasers {
		p.Featureds = append(p.Featureds, models.Featured{Rank: i, Teaser: teaser.Id})
		p.Teasers = append(p.Teasers, teaser)
	}

	writeJSON(w, p)
}

type FeedsSerializer struct {
	Feeds []FeedSerializer
}

type FeedSerializer struct {
	Id          int64
	Title       string
	Description string
	Url         string
	Thumb       string
	Cover       string
	Updated     time.Time
}

func FeedsIndexHandler(w http.ResponseWriter, r *http.Request) {
	p := FeedsSerializer{}
	p.Feeds = make([]FeedSerializer, 0)

	feeds, err := models.Feeds()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	for _, feed := range feeds {
		p.Feeds = append(p.Feeds, FeedSerializer{
			Id:          feed.Id,
			Title:       feed.Title,
			Description: feed.Description,
			Url:         feed.Url,
			Thumb:       fmt.Sprintf("http://s3.amazonaws.com/flowfeeds2/feeds/%d/thumb.jpg", feed.Id),
			Cover:       fmt.Sprintf("http://s3.amazonaws.com/flowfeeds2/feeds/%d/cover.jpg", feed.Id),
		})
	}

	writeJSON(w, p)
}

type FeedShowSerializer struct {
	Feed models.Feed
}

func FeedShowHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		panic(err)
	}

	feed, err := models.FindFeedById(int64(id))
	if err != nil {
		panic(err)
	}

	serializer := FeedShowSerializer{}
	serializer.Feed = feed
	writeJSON(w, serializer)
}
