package datastore

import (
	"errors"
	"fmt"
	"strings"

	"bitbucket.org/nerdyworm/go-flowfeeds/models"
	"github.com/lann/squirrel"
)

type FeedsStore interface {
	Get(id int64) (*models.Feed, error)
	GetIds(ids []int64) ([]*models.Feed, error)
	List() ([]*models.Feed, error)
	FindByURL(string) (*models.Feed, error)
	Ensure(*models.Feed) error
}

type feedsStore struct{ *Datastore }

func (s *feedsStore) Get(id int64) (*models.Feed, error) {
	feed := &models.Feed{}
	err := s.db.Get(feed, "select * from feed where id = $1", id)
	return feed, err
}

func (s *feedsStore) GetIds(ids []int64) ([]*models.Feed, error) {
	feeds := []*models.Feed{}

	if len(ids) == 0 {
		return feeds, nil
	}

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

func (s *feedsStore) FindByURL(url string) (*models.Feed, error) {
	feeds := []*models.Feed{}

	err := s.db.Select(&feeds, "select * from feed where url=$1 limit 1", url)
	if err != nil {
		return nil, err
	}

	if len(feeds) == 0 {
		return nil, errors.New("Feed not found")
	}

	return feeds[0], err
}

func (s *feedsStore) Ensure(feed *models.Feed) error {
	_, err := s.db.Exec("insert into feed (url, title, description, image, updated) values($1, $2, $3, $4, now())", feed.Url, feed.Title, feed.Description, feed.Image)
	if isDupeErrorOf(err, "feeds_url_unique") {
		err = nil
	}

	if err != nil {
		return err
	}

	f, err := s.FindByURL(feed.Url)
	if err != nil {
		return err
	}

	feed.Id = f.Id
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
