package datastore

import "bitbucket.org/nerdyworm/go-flowfeeds/models"

type EpisodesStore interface {
	Get(id int64) (*models.Episode, error)
	GetForUser(*models.User, int64) (*models.Episode, error)
	ToggleFavoriteForUser(*models.User, int64) error
	Related(id int64) ([]*models.Episode, error)
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

func (s *episodesStore) ToggleFavoriteForUser(user *models.User, id int64) error {
	tx, err := s.db.Beginx()
	if err != nil {
		return err
	}

	exists := false
	err = tx.Get(&exists, "select exists(select 1 from favorite where user_id=$1 and episode_id=$2)", user.Id, id)
	if err != nil {
		return err
	}

	if !exists {
		tx.Exec("insert into favorite(user_id, episode_id) values($1, $2)", user.Id, id)
	} else {
		tx.Exec("delete from favorite where user_id=$1 and episode_id=$2", user.Id, id)
	}

	return tx.Commit()
}

func (s *episodesStore) Related(id int64) ([]*models.Episode, error) {
	related := []*models.Episode{}
	err := s.db.Select(&related, "select * from episode where id <> $1 order by random() limit 12", id)
	return related, err
}
