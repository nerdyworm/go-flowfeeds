package datastore

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var (
	DB *sqlx.DB
)

func Connect(config string) error {
	var err error
	DB, err = sqlx.Connect("postgres", config)
	return err
}

type Datastore struct {
	Episodes EpisodesStore
	Feeds    FeedsStore
	db       *sqlx.DB
}

func NewDatastore() *Datastore {
	d := &Datastore{}
	d.db = DB
	d.Episodes = &episodesStore{d}
	d.Feeds = &feedsStore{d}
	return d
}
