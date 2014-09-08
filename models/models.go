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
	x.ShowSQL = true
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

type Listen struct {
	Id        int64
	UserId    int64
	EpisodeId int64
}

type Favorite struct {
	Id        int64
	UserId    int64
	EpisodeId int64
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

	err := x.OrderBy("published desc").Limit(25, 0).Find(&episodes)
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

func FindRelatedTeasers(episodeId int64) ([]Teaser, error) {
	related := []Episode{}
	err := x.Where("id <> ?", episodeId).OrderBy("random()").Limit(5).Find(&related)
	if err != nil {
		return nil, err
	}

	teasers := make([]Teaser, len(related))
	for i := range related {
		teasers[i] = related[i].Teaser()
	}

	return teasers, err
}

func CreateListen(user User, episodeId int64) (Listen, error) {
	listen := Listen{UserId: user.Id, EpisodeId: episodeId}
	_, err := x.Insert(&listen)
	return listen, err
}

func FindListensForEpisode(id int64) ([]Listen, error) {
	listens := []Listen{}
	err := x.Where("episode_id = ?", id).Limit(8).Find(&listens)
	return listens, err
}

func CreateFavorite(user User, episodeId int64) (Favorite, error) {
	favorite := Favorite{UserId: user.Id, EpisodeId: episodeId}
	_, err := x.Insert(&favorite)
	return favorite, err
}

func FindFavoritesForEpisode(id int64) ([]Favorite, error) {
	favorites := []Favorite{}
	err := x.Where("episode_id = ?", id).Limit(8).Find(&favorites)
	return favorites, err
}
