package models

import (
	"fmt"
	"time"

	"code.google.com/p/go.crypto/bcrypt"

	"strings"

	"github.com/go-xorm/xorm"
	_ "github.com/lib/pq"
)

var x *xorm.Engine

var (
	PasswordCost = bcrypt.DefaultCost
)

func Connect(config string) error {
	var err error
	x, err = xorm.NewEngine("postgres", config)
	x.ShowSQL = true
	return err
}

func Close() {
	x.Close()
}

type Episode struct {
	Id             int64
	FeedId         int64
	Guid           string
	Title          string
	Description    string
	Url            string
	Image          string
	Published      time.Time
	ListensCount   int
	FavoritesCount int
	Favorited      bool `xorm:"-"`
	Listened       bool `xorm:"-"`
}

type Listen struct {
	Id        int64
	UserId    int64
	EpisodeId int64
}

type Favorite struct {
	Id        int64
	UserId    int64
	EpisodeId int64
}

type Top100 struct {
	Rank   int `json:"Id"`
	Teaser int64
}

type Featured struct {
	Rank    int `json:"Id"`
	Episode int64
}

type Feed struct {
	Id          int64
	Title       string
	Description string
	Url         string
	Image       string
	Updated     time.Time
}

func EnsureEpisode(episode *Episode) error {
	_, err := x.Insert(episode)
	if isDupeErrorOf(err, "episodes_guid_unique") {
		return nil
	}

	return err
}

func EnsureFeed(feed *Feed) error {
	_, err := x.Insert(feed)
	if isDupeErrorOf(err, "feeds_url_unique") {
		return nil
	}

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

func FeaturedEpisodes(user User) ([]Episode, []Feed, []Listen, []Favorite, error) {
	episodes := []Episode{}
	feeds := []Feed{}
	listens := []Listen{}
	favorites := []Favorite{}

	err := x.OrderBy("published desc").Limit(25, 0).Find(&episodes)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	episodeIds := []int64{}
	feedIds := []int64{}

	for _, episode := range episodes {
		feedIds = append(feedIds, episode.FeedId)
		episodeIds = append(episodeIds, episode.Id)
	}

	if len(feedIds) > 0 {
		err = x.In("id", feedIds).Find(&feeds)
	}

	if len(episodeIds) > 0 {
		err = x.Where("user_id = ?", user.Id).In("episode_id", episodeIds).Find(&listens)
		err = x.Where("user_id = ?", user.Id).In("episode_id", episodeIds).Find(&favorites)
	}

	return episodes, feeds, listens, favorites, err
}

func FindEpisodeById(id int64) (Episode, error) {
	episode := Episode{}
	_, err := x.Id(id).Get(&episode)
	return episode, err
}

func FindEpisodeByIdForUser(id int64, user User) (Episode, error) {
	episode, err := FindEpisodeById(id)
	if err != nil {
		return episode, err
	}

	listens, err := x.Where("episode_id = ? AND user_id = ?", id, user.Id).Count(&Listen{})
	if err != nil {
		return episode, err
	}

	episode.Listened = listens > 0

	favs, err := x.Where("episode_id = ? AND user_id = ?", id, user.Id).Count(&Favorite{})
	if err != nil {
		return episode, err
	}

	episode.Favorited = favs > 0

	return episode, err
}

func Feeds() ([]Feed, error) {
	feeds := []Feed{}
	err := x.OrderBy("lower(title) asc").Find(&feeds)
	return feeds, err
}

func FindFeedById(id int64) (Feed, error) {
	feed := Feed{}
	_, err := x.Id(id).Get(&feed)
	return feed, err
}

func FindFeedByIds(ids []int64) ([]Feed, error) {
	feeds := []Feed{}
	err := x.In("id", ids).Find(&feeds)
	return feeds, err
}

func FindFeedByURL(url string) (Feed, error) {
	feed := Feed{}
	_, err := x.Where("url=?", url).Get(&feed)
	return feed, err
}

func FindRelatedEpisodes(episodeId int64) ([]Episode, error) {
	related := []Episode{}
	err := x.Where("id <> ?", episodeId).OrderBy("random()").Limit(5).Find(&related)
	if err != nil {
		return nil, err
	}

	return related, err
}

func CreateListen(user User, episodeId int64) (Listen, error) {
	listen := Listen{UserId: user.Id, EpisodeId: episodeId}
	_, err := x.Insert(&listen)
	return listen, err
}

func FindListensForEpisode(id int64) ([]Listen, []User, error) {
	listens := []Listen{}
	users := []User{}

	err := x.Where("episode_id = ?", id).Limit(8).Find(&listens)
	if err != nil || len(listens) == 0 {
		return listens, users, err
	}

	ids := []int64{}
	for i := range listens {
		ids = append(ids, listens[i].UserId)
	}

	err = x.Table("users").In("id", ids).Find(&users)
	return listens, users, err
}

func ToggleFavorite(user User, episodeId int64) error {
	favorite := Favorite{UserId: user.Id, EpisodeId: episodeId}
	_, err := x.Where("user_id = ? and episode_id = ?", user.Id, episodeId).Get(&favorite)

	if favorite.Id == 0 {
		_, err = x.Insert(&favorite)
	} else {
		_, err = x.Exec("delete from favorite where user_id = ? and episode_id = ?", user.Id, episodeId)
	}

	return err
}

func CreateFavorite(user User, episodeId int64) (Favorite, error) {
	favorite := Favorite{UserId: user.Id, EpisodeId: episodeId}
	_, err := x.Insert(&favorite)
	return favorite, err
}

func DeleteFavorite(user User, id int64) error {
	_, err := x.Where("user_id = ?", user.Id).Delete(Favorite{Id: id})
	return err
}

func FindFavoritesForEpisode(id int64) ([]Favorite, []User, error) {
	favorites := []Favorite{}
	users := []User{}

	err := x.Where("episode_id = ?", id).Limit(8).Find(&favorites)
	if err != nil || len(favorites) == 0 {
		return favorites, users, err
	}

	ids := []int64{}
	for i := range favorites {
		ids = append(ids, favorites[i].UserId)
	}

	err = x.Table("users").In("id", ids).Find(&users)
	return favorites, users, err
}

func FindFavoriteById(id int64) (Favorite, error) {
	favorite := Favorite{}
	_, err := x.Id(id).Get(&favorite)
	return favorite, err
}
