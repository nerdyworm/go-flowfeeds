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
	Id             int64
	FeedId         int64 `db:"feed_id"`
	Guid           string
	Title          string
	Description    string
	Url            string
	Image          string
	Published      time.Time
	ListensCount   int  `db:"listens_count"`
	FavoritesCount int  `db:"favorites_count"`
	Favorited      bool `xorm:"-",db:"-"`
	Listened       bool `xorm:"-",db:"-"`
}

type Listen struct {
	Id        int64
	UserId    int64 `db:"user_id"`
	EpisodeId int64 `db:"episode_id"`
}

type Favorite struct {
	Id        int64
	UserId    int64 `db:"user_id"`
	EpisodeId int64 `db:"episode_id"`
}

type Featured struct {
	Rank    int `json:"Id"`
	Episode int64
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

const DefaultPerPage = 10

type ListOptions struct {
	PerPage int
	Page    int
}

func (o ListOptions) PageOrDefault() int {
	if o.Page <= 0 {
		return 1
	}

	return o.Page

}

func (o ListOptions) Offset() int {
	return (o.PageOrDefault() - 1) * o.PerPageOrDefault()
}

func (o ListOptions) PerPageOrDefault() int {
	if o.PerPage <= 0 {
		return DefaultPerPage
	}
	return o.PerPage
}

func FindFeedByURL(url string) (Feed, error) {
	feed := Feed{}
	_, err := x.Where("url=?", url).Get(&feed)
	return feed, err
}

func CreateListen(user User, episodeId int64) (Listen, error) {
	listen := Listen{UserId: user.Id, EpisodeId: episodeId}
	_, err := x.Insert(&listen)
	return listen, err
}

func FindListensForEpisode(id int64) ([]Listen, []User, error) {
	listens := []Listen{}
	users := []User{}

	err := x.Where("episode_id = ?", id).OrderBy("id desc").Limit(20).Find(&listens)
	if err != nil || len(listens) == 0 {
		return listens, users, err
	}

	ids := []int64{}
	for i := range listens {
		ids = append(ids, listens[i].UserId)
	}

	err = x.Table("users").In("id", ids).Find(&users)
	return listens, users, err
}

func FindFavoritesForEpisode(id int64) ([]Favorite, []User, error) {
	favorites := []Favorite{}
	users := []User{}

	err := x.Where("episode_id = ?", id).OrderBy("id desc").Limit(20).Find(&favorites)
	if err != nil || len(favorites) == 0 {
		return favorites, users, err
	}

	ids := []int64{}
	for i := range favorites {
		ids = append(ids, favorites[i].UserId)
	}

	err = x.Table("users").In("id", ids).Find(&users)
	return favorites, users, err
}
