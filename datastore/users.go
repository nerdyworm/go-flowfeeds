package datastore

import (
	"database/sql"

	"github.com/nerdyworm/go-flowfeeds/models"
	"github.com/lann/squirrel"
)

type UsersStore interface {
	Get(id int64) (*models.User, error)
	GetIds(ids []int64) ([]*models.User, error)
	Insert(*models.User) error
	Exists(email string) (bool, error)
	FindForSignin(email string) (*models.User, error)
}

type usersStore struct{ *Datastore }

func (s *usersStore) Get(id int64) (*models.User, error) {
	users, err := s.GetIds([]int64{id})
	if err != nil {
		return nil, err
	}

	if len(users) == 0 {
		return nil, models.ErrNotFound
	}

	return users[0], nil
}

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

func (s *usersStore) Insert(user *models.User) error {
	row := s.db.QueryRow("insert into users (email, encrypted_password) values($1,$2) returning id",
		user.Email, user.EncryptedPassword)
	return row.Scan(&user.Id)
}

func (s *usersStore) Exists(email string) (bool, error) {
	exists := false
	row := s.db.QueryRow("select exists(select 1 from users where lower(email) = lower($1))", email)
	return exists, row.Scan(&exists)
}

func (s *usersStore) FindForSignin(email string) (*models.User, error) {
	user := &models.User{}

	row := s.db.QueryRow("select * from users where lower(email) = lower($1) limit 1", email)
	err := row.Scan(
		&user.Id,
		&user.Email,
		&user.EncryptedPassword,
	)

	if err == sql.ErrNoRows {
		return user, models.ErrNotFound
	}

	return user, nil
}
