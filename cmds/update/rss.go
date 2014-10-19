package update

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"

	"bitbucket.org/nerdyworm/go-flowfeeds/config"
	"bitbucket.org/nerdyworm/go-flowfeeds/datastore"
	"bitbucket.org/nerdyworm/go-flowfeeds/models"
	"bitbucket.org/nerdyworm/go-flowfeeds/rss"
)

type Fetcher interface {
	Fetch() error
	Feed() models.Feed
	Episodes() []models.Episode
}

type Updater interface {
	Fetchers() []Fetcher
}

type Collection struct {
	Name string
	Urls []string
}

func Rss() {
	start := time.Now()

	var wg sync.WaitGroup
	for _, c := range readCollections() {
		for _, url := range c.Urls {
			wg.Add(1)
			go func(url string) {
				defer wg.Done()
				s := time.Now()
				rssFeed := RssFeed{Url: url}
				ensureContentsOfFetcher(&rssFeed)
				log.Printf("%s in %v\n", rssFeed.Feed().Url, time.Since(s))
			}(url)
		}
	}

	wg.Wait()
	log.Printf("finished update in %v\n", time.Since(start))
}

func ensureContentsOfFetcher(fetcher Fetcher) error {
	err := fetcher.Fetch()
	if err != nil {
		log.Printf("error updating: %v\n", err)
		return err
	}

	store := datastore.NewDatastore()
	feed := fetcher.Feed()

	err = store.Feeds.Ensure(&feed)
	if err != nil {
		log.Fatal(err)
	}

	for _, e := range fetcher.Episodes() {
		e.Feed = feed.Id
		e.Image = feed.Image
		err = store.Episodes.Ensure(&e)
		if err != nil {
			return err
		}
	}

	return nil
}

func readCollections() []Collection {
	collections := make([]Collection, 0)

	data, err := ioutil.ReadFile(config.COLLECTIONS)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	err = json.Unmarshal(data, &collections)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	return collections
}

type RssFeed struct {
	Url      string
	feed     models.Feed
	episodes []models.Episode
}

func (r *RssFeed) Fetch() error {
	resp, err := http.Get(r.Url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	rssFeed, err := rss.Read(resp.Body)
	if err != nil {
		return err
	}

	r.feed = models.Feed{
		Url:         r.Url,
		Title:       rssFeed.Title(),
		Description: rssFeed.Description(),
		Image:       rssFeed.Image(),
		Updated:     time.Now(),
	}

	r.episodes = []models.Episode{}
	for _, item := range rssFeed.Items() {
		episode := models.Episode{
			Guid:        item.Guid,
			Title:       item.Title,
			Description: item.Description,
			Url:         item.Enclosure.URL,
			Published:   item.Published(),
		}

		r.episodes = append(r.episodes, episode)
	}

	return nil
}

func (r *RssFeed) Feed() models.Feed {
	return r.feed
}

func (r *RssFeed) Episodes() []models.Episode {
	return r.episodes
}
