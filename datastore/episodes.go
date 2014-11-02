package datastore

import (
	"log"

	"github.com/lann/squirrel"

	"bitbucket.org/nerdyworm/go-flowfeeds/models"
)

type EpisodesStore interface {
	Get(id int64) (*models.Episode, error)
	GetForUser(*models.User, int64) (*models.Episode, error)
	Related(id int64) ([]*models.Episode, []*models.Feed, error)
	ListFor(*models.User, EpisodeListOptions) (Episodes, []*models.Feed, error)
	Listens(id int64) ([]*models.Listen, []*models.User, error)
	Favorites(id int64) ([]*models.Favorite, []*models.User, error)
	ToggleFavoriteForUser(*models.User, int64) error
	Ensure(*models.Episode) error
}

type EpisodeListOptions struct {
	ListOptions
	Feed int64
}

func (f EpisodeListOptions) OrderOrDefault() string {
	return "published desc"
}

type Episodes struct {
	Episodes []*models.Episode
	Total    int
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

	s.db.Get(&episode.Listened, `select exists(select 1 from listen where "user"=$1 and "episode"=$2)`, user.Id, episode.Id)
	s.db.Get(&episode.Favorited, `select exists(select 1 from favorite where "user"=$1 and "episode"=$2)`, user.Id, episode.Id)

	return episode, err
}

func (s *episodesStore) ToggleFavoriteForUser(user *models.User, id int64) error {
	tx, err := s.db.Beginx()
	if err != nil {
		return err
	}

	exists := false
	err = tx.Get(&exists, `select exists(select 1 from favorite where "user"=$1 and episode=$2)`, user.Id, id)
	if err != nil {
		return err
	}

	if exists {
		tx.Exec(`delete from favorite where "user"=$1 and "episode"=$2`, user.Id, id)
	} else {
		tx.Exec(`insert into favorite("user", episode) values($1, $2)`, user.Id, id)
	}

	return tx.Commit()
}

func (s *episodesStore) Related(id int64) ([]*models.Episode, []*models.Feed, error) {
	episodes := []*models.Episode{}
	feeds := []*models.Feed{}

	err := s.db.Select(&episodes, "select * from episode where id <> $1 order by random() limit 12", id)
	if err != nil {
		return episodes, feeds, err
	}

	ids := []int64{}
	for i := range episodes {
		ids = append(ids, episodes[i].Feed)
	}

	feeds, err = s.Feeds.GetIds(ids)
	return episodes, feeds, err
}

func (s *episodesStore) ListFor(user *models.User, options EpisodeListOptions) (Episodes, []*models.Feed, error) {
	episodes := Episodes{}
	episodes.Episodes = []*models.Episode{}
	feeds := []*models.Feed{}

	countQuery := s.QueryBuilder().Select("count(*)").From("episode")
	q := s.QueryBuilder().Select("*").From("episode")

	if options.Feed != 0 {
		q = q.Where("feed = ?", options.Feed)
		countQuery = countQuery.Where("feed = ?", options.Feed)
	}

	query, args, _ := countQuery.ToSql()
	err := s.db.Get(&episodes.Total, query, args...)

	q = q.Limit(uint64(options.PerPageOrDefault())).Offset(uint64(options.Offset())).OrderBy(options.OrderOrDefault())

	query, args, _ = q.ToSql()
	err = s.db.Select(&episodes.Episodes, query, args...)

	episodeIds := []int64{}
	feedIds := []int64{}

	for _, episode := range episodes.Episodes {
		feedIds = append(feedIds, episode.Feed)
		episodeIds = append(episodeIds, episode.Id)
	}

	feeds, err = s.Feeds.GetIds(feedIds)
	if err != nil {
		return episodes, feeds, err
	}

	err = s.setEpisodeStateFor(user, episodes.Episodes)
	return episodes, feeds, err
}

func (s *episodesStore) setEpisodeStateFor(user *models.User, episodes []*models.Episode) error {
	if len(episodes) == 0 {
		return nil
	}

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
		Where(squirrel.Eq{"listen.episode": ids, "listen.user": user.Id})

	query, args, err := listensQuery.ToSql()
	if err != nil {
		return err
	}

	err = s.db.Select(&listens, query, args...)
	if err != nil {
		log.Printf(`Select: "%s" err: "%v"\n`, query, err)
		return err
	}

	favoritesQuery := builder.Select("*").From("favorite").
		Where(squirrel.Eq{"favorite.episode": ids, "favorite.user": user.Id})

	query, args, err = favoritesQuery.ToSql()
	if err != nil {
		return err
	}

	err = s.db.Select(&favorites, query, args...)
	if err != nil {
		log.Printf(`Select: "%s" err: "%v"\n`, query, err)
		return err
	}

	for _, listen := range listens {
		if _, ok := episodesToListens[listen.Episode]; !ok {
			episodesToListens[listen.Episode] = true
		}
	}

	for _, favorite := range favorites {
		if _, ok := episodesToFavorites[favorite.Episode]; !ok {
			episodesToFavorites[favorite.Episode] = true
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

func (s *episodesStore) Favorites(id int64) ([]*models.Favorite, []*models.User, error) {
	favorites := []*models.Favorite{}
	users := []*models.User{}

	err := s.db.Select(&favorites, "select * from favorite where episode=$1", id)
	if err != nil {
		return favorites, users, err
	}

	ids := []int64{}
	for i := range favorites {
		ids = append(ids, favorites[i].User)
	}

	users, err = s.Users.GetIds(ids)

	return favorites, users, err
}

func (s *episodesStore) Listens(id int64) ([]*models.Listen, []*models.User, error) {
	listens := []*models.Listen{}
	users := []*models.User{}

	query := "select * from listen where episode=$1"
	err := s.db.Select(&listens, query, id)
	if err != nil {
		log.Printf(`Select: "%s" err: "%v"\n`, query, err)
		return listens, users, err
	}

	ids := []int64{}
	for i := range listens {
		ids = append(ids, listens[i].User)
	}

	users, err = s.Users.GetIds(ids)

	return listens, users, err
}

func (s *episodesStore) Ensure(episode *models.Episode) error {
	insert := s.QueryBuilder().Insert("episode").
		Columns("feed", "guid", "title", "description", "url", "image", "published").
		Values(episode.Feed, episode.Guid, episode.Title, episode.Description, episode.Url, episode.Image, episode.Published)

	query, args, _ := insert.ToSql()
	_, err := s.db.Exec(query, args...)
	if isDupeErrorOf(err, "episodes_guid_unique") {
		err = nil
	}

	if err != nil {
		log.Printf(`Exec: "%s" err: "%v"\n`, query, err)
		return err
	}

	episode, err = s.findByGuid(episode.Guid)
	if err != nil {
		return err
	}

	return err
}

func (s *episodesStore) findByGuid(guid string) (*models.Episode, error) {
	episode := &models.Episode{}
	err := s.db.Get(episode, "select * from episode where guid = $1", guid)
	return episode, err
}
