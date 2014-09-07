package rss

import (
	"crypto/md5"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"
)

func fileExists(name string) bool {
	_, err := os.Stat(name)
	return !os.IsNotExist(err)
}

func cachedFeed(url string) (io.ReadCloser, error) {
	hash := md5.New()
	io.WriteString(hash, url)
	cache := fmt.Sprintf("/tmp/%x", hash.Sum(nil))

	if !fileExists(cache) {
		r, err := http.Get(url)
		if err != nil {
			return nil, err
		}
		defer r.Body.Close()

		f, err := os.Create(cache)
		if err != nil {
			return nil, err
		}
		defer f.Close()

		io.Copy(f, r.Body)
	}

	return os.Open(cache)
}

func TestBadPocastTitle(t *testing.T) {
	url := "http://feeds.feedburner.com/PhilDahouseCat"

	data, err := cachedFeed(url)
	if err != nil {
		t.Fatal(err)
	}
	defer data.Close()

	feed, err := Read(data)
	if err != nil {
		t.Fatal(err)
	}

	if feed.Title() == "" {
		t.Fatal("Expected there to be a title")
	}
}
