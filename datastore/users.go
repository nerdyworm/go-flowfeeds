package datastore

import (
	"bitbucket.org/nerdyworm/go-flowfeeds/models"
	"github.com/lann/squirrel"
)

type UsersStore interface {
	GetIds(ids []int64) ([]*models.User, error)
}

type usersStore struct{ *Datastore }

func (s *usersStore) GetIds(ids []int64) ([]*models.User, error) {
	users := []*models.User{}
	if len(ids) == 0 {
		return users, nil
	}

	query, args := s.usersByIdsQuery(ids)
	err := s.db.Select(&users, query, args...)
	if err != nil {
		return users, err
	}

	return users, nil
}

func (s *usersStore) usersByIdsQuery(ids []int64) (string, []interface{}) {
	usersQuery := s.QueryBuilder().
		Select("id, email").
		From("users").
		Where(squirrel.Eq{"id": ids})

	query, args, err := usersQuery.ToSql()
	if err != nil {
		panic(err)
	}

	return query, args
}
