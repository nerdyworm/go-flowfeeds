package datastore

import (
	"github.com/jmoiron/sqlx"
	"github.com/lann/squirrel"
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

func (s *Datastore) QueryBuilder() squirrel.StatementBuilderType {
	return squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
}

func NewDatastore() *Datastore {
	d := &Datastore{}
	d.db = DB
	d.Episodes = &episodesStore{d}
	d.Feeds = &feedsStore{d}
	return d
}
