package datastore

import "bitbucket.org/nerdyworm/go-flowfeeds/models"

type ListensStore interface {
	Create(*models.User, int64) (*models.Listen, error)
}

type listensStore struct{ *Datastore }

func (s *listensStore) Create(user *models.User, id int64) (*models.Listen, error) {
	listen := &models.Listen{UserId: user.Id, EpisodeId: id}
	row := s.db.QueryRow("insert into listen (user_id, episode_id) values($1, $2) returning id", listen.UserId, listen.EpisodeId)
	err := row.Scan(&listen.Id)
	return listen, err
}
