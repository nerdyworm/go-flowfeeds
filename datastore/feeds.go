package datastore

import "bitbucket.org/nerdyworm/go-flowfeeds/models"

type FeedsStore interface {
	Get(id int64) (*models.Feed, error)
	List() ([]*models.Feed, error)
}

type feedsStore struct{ *Datastore }

func (s *feedsStore) Get(id int64) (*models.Feed, error) {
	feed := &models.Feed{}
	err := s.db.Get(feed, "select * from feed where id = $1", id)
	return feed, err
}

func (s *feedsStore) List() ([]*models.Feed, error) {
	feeds := []*models.Feed{}
	err := s.db.Select(&feeds, "select * from feed order by lower(title) asc")
	return feeds, err
}
