package datastore

import (
	"github.com/lann/squirrel"

	"bitbucket.org/nerdyworm/go-flowfeeds/models"
)

type EpisodesStore interface {
	Get(id int64) (*models.Episode, error)
	GetForUser(*models.User, int64) (*models.Episode, error)
	ToggleFavoriteForUser(*models.User, int64) error
	Related(id int64) ([]*models.Episode, error)
	ListFor(*models.User, models.ListOptions) ([]*models.Episode, []*models.Feed, error)
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

func (s *episodesStore) setEpisodeStateFor(user *models.User, episodes []*models.Episode) error {
	builder := s.QueryBuilder()

	listens := []*models.Listen{}
	episodesToListens := make(map[int64]bool)

	favorites := []*models.Favorite{}
	episodesToFavorites := make(map[int64]bool)

	ids := []int64{}
	for _, episode := range episodes {
		ids = append(ids, episode.Id)
	}

	listensQuery := builder.Select("*").From("listen").
		Where(squirrel.Eq{"episode_id": ids, "user_id": user.Id})

	query, args, err := listensQuery.ToSql()
	if err != nil {
		return err
	}

	err = s.db.Select(&listens, query, args...)
	if err != nil {
		return err
	}

	favoritesQuery := builder.Select("*").From("favorite").
		Where(squirrel.Eq{"episode_id": ids, "user_id": user.Id})

	query, args, err = favoritesQuery.ToSql()
	if err != nil {
		return err
	}

	err = s.db.Select(&favorites, query, args...)
	if err != nil {
		return err
	}

	for _, listen := range listens {
		if _, ok := episodesToListens[listen.EpisodeId]; !ok {
			episodesToListens[listen.EpisodeId] = true
		}
	}

	for _, favorite := range favorites {
		if _, ok := episodesToFavorites[favorite.EpisodeId]; !ok {
			episodesToFavorites[favorite.EpisodeId] = true
		}
	}

	for i, e := range episodes {
		if listened, ok := episodesToListens[e.Id]; ok {
			episodes[i].Listened = listened
		}

		if favorited, ok := episodesToFavorites[e.Id]; ok {
			episodes[i].Favorited = favorited
		}
	}

	return nil
}

func (s *episodesStore) ListFor(user *models.User, options models.ListOptions) ([]*models.Episode, []*models.Feed, error) {
	episodes := []*models.Episode{}
	feeds := []*models.Feed{}

	err := s.db.Select(&episodes, "select * from episode order by published desc limit $1 offset $2", options.PerPageOrDefault(), options.Offset())

	episodeIds := []int64{}
	feedIds := []int64{}

	for _, episode := range episodes {
		feedIds = append(feedIds, episode.FeedId)
		episodeIds = append(episodeIds, episode.Id)
	}

	if len(feedIds) > 0 {
		feeds, err = s.Feeds.GetIds(feedIds)
		if err != nil {
			return episodes, feeds, err
		}
	}

	if len(episodeIds) > 0 {
		err = s.setEpisodeStateFor(user, episodes)
	}

	return episodes, feeds, err
}
