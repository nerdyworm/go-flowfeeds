package models

import (
	"fmt"
	"time"

	"code.google.com/p/go.crypto/bcrypt"

	"strings"

	"github.com/go-xorm/xorm"
	_ "github.com/lib/pq"
)

var x *xorm.Engine

var (
	PasswordCost = bcrypt.DefaultCost
)

func Connect(config string) error {
	var err error
	x, err = xorm.NewEngine("postgres", config)
	//x.ShowSQL = true
	return err
}

func Close() {
	x.Close()
}

type Episode struct {
	Id          int64
	FeedId      int64
	Guid        string
	Title       string
	Description string
	Url         string
	Image       string
	Published   time.Time
}

type Comment struct {
	Id       int64
	Body     string
	AuthorId int64
}

type Top100 struct {
	Rank   int `json:"Id"`
	Teaser int64
}

type Teaser struct {
	Id          int64
	Episode     int64
	Title       string
	Description string
	Url         string
	Image       string
	FeedId      int64
	Published   time.Time
}

func (e Episode) Teaser() Teaser {
	return Teaser{
		Id:          e.Id,
		Episode:     e.Id,
		Title:       e.Title,
		Description: e.Description,
		Url:         e.Url,
		Image:       e.Image,
		FeedId:      e.FeedId,
		Published:   e.Published,
	}
}

type Featured struct {
	Rank   int `json:"Id"`
	Teaser int64
}

type Feed struct {
	Id          int64
	Title       string
	Description string
	Url         string
	Image       string
	Updated     time.Time
}

func EnsureEpisode(episode *Episode) error {
	_, err := x.Insert(episode)
	if isDupeErrorOf(err, "episodes_guid_unique") {
		return nil
	}

	return err
}

func EnsureFeed(feed *Feed) error {
	_, err := x.Insert(feed)
	if isDupeErrorOf(err, "feeds_url_unique") {
		return nil
	}

	return err
}

func isDupeErrorOf(err error, indexName string) bool {
	if err == nil {
		return false
	}

	message := fmt.Sprintf(`duplicate key value violates unique constraint "%s"`, indexName)
	if strings.Contains(err.Error(), message) {
		return true
	}

	return false
}

func FeaturedEpisodeTeasers() ([]Teaser, []Feed, error) {
	episodes := []Episode{}
	feeds := []Feed{}

	err := x.OrderBy("published desc").Limit(20, 0).Find(&episodes)
	if err != nil {
		return nil, nil, err
	}

	feedIds := []int64{}
	teasers := make([]Teaser, len(episodes))
	for i, episode := range episodes {
		teasers[i] = episode.Teaser()
		feedIds = append(feedIds, episode.FeedId)
	}

	if len(feedIds) > 0 {
		err = x.In("id", feedIds).Find(&feeds)
	}

	return teasers, feeds, err
}

func FindEpisodeById(id int64) (Episode, error) {
	episode := Episode{}
	_, err := x.Id(id).Get(&episode)
	return episode, err
}

func Feeds() ([]Feed, error) {
	feeds := []Feed{}
	err := x.OrderBy("lower(title) asc").Find(&feeds)
	return feeds, err
}

func FindFeedById(id int64) (Feed, error) {
	feed := Feed{}
	_, err := x.Id(id).Get(&feed)
	return feed, err
}

func FindFeedByURL(url string) (Feed, error) {
	feed := Feed{}
	_, err := x.Where("url=?", url).Get(&feed)
	return feed, err
}
