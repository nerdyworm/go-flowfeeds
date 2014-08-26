package main

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"sync"
	"time"

	"bitbucket.org/nerdyworm/go-flowfeeds/models"
	"bitbucket.org/nerdyworm/go-flowfeeds/modules/rss"

	"github.com/mitchellh/goamz/aws"
	"github.com/mitchellh/goamz/s3"
	"gopkg.in/yaml.v1"
)

type Collection struct {
	Name string
	Urls []string
}

func imageWorker(id int, feeds <-chan models.Feed, results chan<- int) {
	for feed := range feeds {
		log.Printf("worker[%d] making images %s", id, feed.Title)
		if err := makeImages(feed); err != nil {
			log.Printf("worker[%d] error for %s\n\t%v", id, feed.Title, err)
		}
		results <- 1
	}
}

func main() {
	err := models.Connect("dbname=flowfeeds2 sslmode=disable")
	if err != nil {
		log.Fatalln(err)
	}
	defer models.Close()

	startUpdateFromCollectionsYml()
	makeImagesFromFeeds()
}

func makeImagesFromFeeds() {
	feeds := make(chan models.Feed)
	results := make(chan int)

	for w := 1; w <= 10; w++ {
		go imageWorker(w, feeds, results)
	}

	f, _ := models.Feeds()
	total := len(f)

	go func() {
		for _, feed := range f {
			feeds <- feed
		}
		close(feeds)
	}()

	for a := 0; a < total; a++ {
		<-results
	}
}

func makeImages(feed models.Feed) error {
	imagePath := fmt.Sprintf("/tmp/feed-%d", feed.Id)
	thumbnailPath := imagePath + ".thumbnail.jpg"
	coverPath := imagePath + ".cover.jpg"

	if feed.Image == "" {
		return errors.New(fmt.Sprintf("No image for %s", feed.Url))
	}

	auth := aws.Auth{os.Getenv("S3_KEY"), os.Getenv("S3_SEC"), ""}
	sss := s3.New(auth, aws.USEast)
	bucket := sss.Bucket("flowfeeds2")

	log.Printf("getting %s", feed.Image)
	resp, err := http.Get(feed.Image)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	output, err := os.Create(imagePath)
	if err != nil {
		return err
	}
	defer output.Close()

	_, err = io.Copy(output, resp.Body)
	if err != nil {
		return err
	}
	log.Printf("got %s", feed.Image)

	/// THUMBNAILS
	out, err := exec.Command("convert", imagePath, "-thumbnail", "300x300^", "-gravity", "center", "-extent", "300x300", thumbnailPath).Output()
	if err != nil {
		log.Println(out)
		return err
	}

	thumbnail, err := os.Open(thumbnailPath)
	if err != nil {
		return err
	}
	defer thumbnail.Close()

	stat, err := thumbnail.Stat()
	if err != nil {
		return err
	}

	key := fmt.Sprintf("/feeds/%d/thumbnail.jpg", feed.Id)
	err = bucket.PutReader(key, thumbnail, stat.Size(), "image/jpeg", s3.PublicRead)
	if err != nil {
		return err
	}

	// COVER
	out, err = exec.Command("convert", imagePath, "-thumbnail", "360x270^", "-gravity", "center", "-extent", "360x270", coverPath).Output()
	if err != nil {
		log.Println(out)
		return err
	}

	cover, err := os.Open(coverPath)
	if err != nil {
		return err
	}
	defer cover.Close()

	stat, err = cover.Stat()
	if err != nil {
		return err
	}

	key = fmt.Sprintf("/feeds/%d/cover.jpg", feed.Id)
	err = bucket.PutReader(key, cover, stat.Size(), "image/jpeg", s3.PublicRead)
	if err != nil {
		return err
	}

	return nil
}

func worker(id int, urls <-chan string) {
	log.Printf("Starting worker %d\n", id)
	for url := range urls {
		log.Println("worker", id, "processing job", url)
		if err := updateUrl(url); err != nil {
			log.Printf("error updating: %s\n\t%v", url, err)
		}
	}
}

func startUpdateFromCollectionsYml() {
	var wg sync.WaitGroup

	for _, c := range readCollectionsYml() {
		for _, url := range c.Urls {
			wg.Add(1)
			go func(url string) {
				defer wg.Done()
				if err := updateUrl(url); err != nil {
					log.Printf("error updating: %s\n\t%v", url, err)
				}
			}(url)
		}
	}

	wg.Wait()
}

func readCollectionsYml() []Collection {
	collections := make([]Collection, 0)

	data, err := ioutil.ReadFile("db/collections.yml")
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	err = yaml.Unmarshal(data, &collections)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	return collections
}

func updateUrl(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	rssFeed, err := rss.Read(resp.Body)
	if err != nil {
		return err
	}

	feed := models.Feed{
		Url:         url,
		Title:       rssFeed.Title(),
		Description: rssFeed.Description(),
		Image:       rssFeed.Image(),
		Updated:     time.Now(),
	}

	err = models.EnsureFeed(&feed)
	if err != nil {
		log.Fatal(err)
	}

	for _, item := range rssFeed.Items() {
		episode := models.Episode{
			Guid:        item.Guid,
			Title:       item.Title,
			Description: item.Description,
			Url:         item.Enclosure.URL,
			Image:       feed.Image,
			Published:   item.Published(),
		}

		err = models.EnsureEpisode(&episode)
		if err != nil {
			log.Fatal(err)
		}
	}

	log.Println("finished", url)

	return nil
}
