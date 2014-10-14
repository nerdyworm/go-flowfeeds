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
	Users    UsersStore
	Listens  ListensStore
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
	d.Users = &usersStore{d}
	d.Listens = &listensStore{d}
	return d
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
