package updates

import (
	"encoding/json"
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

func Run(collection string) {
	startUpdateFromCollections(collection)
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
	if feed.Image == "" {
		return errors.New(fmt.Sprintf("No image for %s", feed.Url))
	}

	image, err := ioutil.TempFile("/tmp", fmt.Sprintf("feed-%d-", feed.Id))
	if err != nil {
		return err
	}
	defer image.Close()

	log.Printf("getting %s", feed.Image)
	resp, err := http.Get(feed.Image)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(image, resp.Body)
	if err != nil {
		return err
	}
	image.Close()
	log.Printf("got %s", feed.Image)

	auth := aws.Auth{os.Getenv("AWS_ACCESS_KEY_ID"), os.Getenv("AWS_SECRET_ACCESS_KEY"), ""}
	sss := s3.New(auth, aws.USEast)
	bucket := sss.Bucket("flowfeeds2")

	/// THUMBNAILS
	key := fmt.Sprintf("/feeds/%d/thumb.jpg", feed.Id)
	err = processImage(image.Name(), "280x280", key, bucket)
	if err != nil {
		os.Remove(image.Name())
		return err
	}

	key = fmt.Sprintf("/feeds/%d/thumb-x2.jpg", feed.Id)
	err = processImage(image.Name(), "560x560", key, bucket)
	if err != nil {
		os.Remove(image.Name())
		return err
	}

	// COVER
	key = fmt.Sprintf("/feeds/%d/cover.jpg", feed.Id)
	err = processImage(image.Name(), "600x600", key, bucket)
	if err != nil {
		os.Remove(image.Name())
		return err
	}

	os.Remove(image.Name())
	return nil
}

func processImage(input string, size string, key string, bucket *s3.Bucket) error {
	output, err := ioutil.TempFile("/tmp", "process-image-")
	if err != nil {
		return err
	}
	output.Close()

	out, err := exec.Command("convert", input, "-thumbnail", size+"^", "-gravity", "center", "-extent", size, output.Name()).Output()
	if err != nil {
		log.Println(out)
		return err
	}

	image, err := os.Open(output.Name())
	if err != nil {
		return err
	}
	defer image.Close()

	stat, err := image.Stat()
	if err != nil {
		return err
	}

	err = bucket.PutReader(key, image, stat.Size(), "image/jpeg", s3.PublicRead)
	if err != nil {
		return err
	}

	os.Remove(output.Name())

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

func startUpdateFromCollections(jsonFile string) {
	var wg sync.WaitGroup

	for _, c := range readCollections(jsonFile) {
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

func readCollections(jsonFile string) []Collection {
	collections := make([]Collection, 0)

	data, err := ioutil.ReadFile(jsonFile)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	err = json.Unmarshal(data, &collections)
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
			FeedId:      feed.Id,
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
