package datastore

import (
	"bitbucket.org/nerdyworm/go-flowfeeds/models"
	"github.com/lann/squirrel"
)

type FeedsStore interface {
	Get(id int64) (*models.Feed, error)
	List() ([]*models.Feed, error)
	GetIds(ids []int64) ([]*models.Feed, error)
}

type feedsStore struct{ *Datastore }

func (s *feedsStore) Get(id int64) (*models.Feed, error) {
	feed := &models.Feed{}
	err := s.db.Get(feed, "select * from feed where id = $1", id)
	return feed, err
}

func (s *feedsStore) GetIds(ids []int64) ([]*models.Feed, error) {
	feeds := []*models.Feed{}

	builder := s.QueryBuilder().
		Select("*").From("feed").
		Where(squirrel.Eq{"id": ids})

	query, args, err := builder.ToSql()
	if err != nil {
		return feeds, err
	}

	err = s.db.Select(&feeds, query, args...)
	return feeds, err
}

func (s *feedsStore) List() ([]*models.Feed, error) {
	feeds := []*models.Feed{}
	err := s.db.Select(&feeds, "select * from feed order by lower(title) asc")
	return feeds, err
}
