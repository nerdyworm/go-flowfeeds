package models

import (
	"time"

	"code.google.com/p/go.crypto/bcrypt"

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
