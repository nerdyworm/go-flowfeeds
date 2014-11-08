package update

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"

	"bitbucket.org/nerdyworm/go-flowfeeds/config"
	"bitbucket.org/nerdyworm/go-flowfeeds/datastore"
	"bitbucket.org/nerdyworm/go-flowfeeds/models"

	"github.com/mitchellh/goamz/aws"
	"github.com/mitchellh/goamz/s3"
)

func Image() {
	makeImagesFromFeeds()
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

func makeImagesFromFeeds() {
	feeds := make(chan models.Feed)
	results := make(chan int)

	for w := 1; w <= 10; w++ {
		go imageWorker(w, feeds, results)
	}

	store := datastore.NewDatastore()

	f, _ := store.Feeds.List()
	total := len(f)

	go func() {
		for _, feed := range f {
			feeds <- *feed
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
	defer os.Remove(image.Name())

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

	log.Printf("got %s", feed.Image)

	auth := aws.Auth{config.AWS_ACCESS_KEY_ID, config.AWS_SECRET_ACCESS_KEY, ""}
	sss := s3.New(auth, aws.USEast)
	bucket := sss.Bucket(config.BUCKET)

	key := fmt.Sprintf("/feeds/%d/cover.jpg", feed.Id)
	err = processImage(image.Name(), "600x600", key, bucket)
	if err != nil {
		return err
	}

	return nil
}

func processImage(input string, size string, key string, bucket *s3.Bucket) error {
	output, err := ioutil.TempFile("/tmp", "process-image-")
	if err != nil {
		return err
	}
	output.Close()
	defer os.Remove(output.Name())

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

	return nil
}
