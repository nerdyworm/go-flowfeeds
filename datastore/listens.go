package datastore

import "bitbucket.org/nerdyworm/go-flowfeeds/models"

type ListensStore interface {
	Create(*models.User, int64) (*models.Listen, error)
}

type listensStore struct{ *Datastore }

func (s *listensStore) Create(user *models.User, id int64) (*models.Listen, error) {
	listen := &models.Listen{User: user.Id, Episode: id}
	row := s.db.QueryRow(`insert into listen ("user", episode) values($1, $2) returning id`, listen.User, listen.Episode)
	err := row.Scan(&listen.Id)
	return listen, err
}
