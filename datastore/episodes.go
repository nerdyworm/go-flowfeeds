package datastore

import "bitbucket.org/nerdyworm/go-flowfeeds/models"

import _ "github.com/lib/pq"

type EpisodesService interface {
	Get(id int64) (*models.Episode, error)
	GetForUser(*models.User, int64) (*models.Episode, error)
}

type episodesStore struct{ *Datastore }

func (s *episodesStore) Get(id int64) (*models.Episode, error) {
	episode := &models.Episode{}
	err := s.db.Get(episode, "select * from episode where id = $1", id)
	return episode, err
}

func (s *episodesStore) GetForUser(user *models.User, id int64) (*models.Episode, error) {
	episode, err := s.Get(id)
	if err != nil {
		return episode, err
	}

	s.db.Get(&episode.Listened, "select exists(select 1 from listen where user_id=$1 and episode_id=$2)", user.Id, episode.Id)
	s.db.Get(&episode.Favorited, "select exists(select 1 from favorite where user_id=$1 and episode_id=$2)", user.Id, episode.Id)

	return episode, err
}
