package serializers

import (
	"fmt"

	"bitbucket.org/nerdyworm/go-flowfeeds/config"
	"bitbucket.org/nerdyworm/go-flowfeeds/datastore"
	"bitbucket.org/nerdyworm/go-flowfeeds/models"
	"bitbucket.org/nerdyworm/go-flowfeeds/web/helpers"
)

type omit *struct{}

type Episode struct {
	*models.Episode
	Cover string
	Links EpisodeLinks `json:"links"`
}

type EpisodeLinks struct {
	Favorites string `json:"favorites"`
	Listens   string `json:"listens"`
	Related   string `json:"related"`
}

func NewEpisode(episode *models.Episode) Episode {
	return Episode{
		Episode: episode,
		Cover:   fmt.Sprintf("%s/feeds/%d/cover.jpg", config.ASSETS, episode.Feed),
		Links: EpisodeLinks{
			Favorites: fmt.Sprintf("/api/v1/episodes/%d/favorites", episode.Id),
			Listens:   fmt.Sprintf("/api/v1/episodes/%d/listens", episode.Id),
			Related:   fmt.Sprintf("/api/v1/episodes/%d/related", episode.Id),
		},
	}
}

type ShowEpisode struct {
	Episode Episode
	Feeds   []Feed
}

func NewShowEpisode(episode *models.Episode, feed *models.Feed) ShowEpisode {
	return ShowEpisode{
		Episode: NewEpisode(episode),
		Feeds:   []Feed{NewFeed(feed)},
	}
}

type Pagination struct {
	Pages int
	Page  int
	Limit int
	Total int
}

type Episodes struct {
	Episodes []Episode
	Feeds    []Feed
	Meta     struct {
		Pagination Pagination
	}
}

func NewEpisodes(episodes datastore.Episodes, feeds []*models.Feed) Episodes {
	s := Episodes{}
	s.Episodes = make([]Episode, len(episodes.Episodes))
	s.Feeds = make([]Feed, len(feeds))
	s.Meta.Pagination.Total = episodes.Total

	for i, episode := range episodes.Episodes {
		s.Episodes[i] = NewEpisode(episode)
	}

	for i, feed := range feeds {
		s.Feeds[i] = NewFeed(feed)
	}

	return s
}

type Feed struct {
	*models.Feed
	Cover string
}

func NewFeed(feed *models.Feed) Feed {
	return Feed{
		Feed:  feed,
		Cover: fmt.Sprintf("%s/feeds/%d/cover.jpg", config.ASSETS, feed.Id),
	}
}

type ShowFeed struct {
	Feed Feed
}

type Feeds struct {
	Feeds []Feed
}

func NewFeeds(feeds []*models.Feed) Feeds {
	s := Feeds{}
	s.Feeds = make([]Feed, len(feeds))
	for i, feed := range feeds {
		s.Feeds[i] = NewFeed(feed)
	}
	return s
}

func NewShowFeed(feed *models.Feed) ShowFeed {
	return ShowFeed{
		Feed: NewFeed(feed),
	}
}

type User struct {
	Id     int64
	Avatar string
}

func NewUser(user *models.User) User {
	return User{
		Id:     user.Id,
		Avatar: helpers.Gravatar(user.Email),
	}
}

type ShowUser struct {
	User User
}

func NewShowUser(user *models.User) ShowUser {
	return ShowUser{NewUser(user)}
}

type Listen struct {
	*models.Listen
}

type Listens struct {
	Listens []Listen
	Users   []User
}

type ShowListen struct {
	Listen   Listen
	Episodes []Episode
}

func NewListen(listen *models.Listen) Listen {
	return Listen{listen}
}

func NewShowListen(listen *models.Listen, episode *models.Episode) ShowListen {
	s := ShowListen{}
	s.Listen = NewListen(listen)

	s.Episodes = make([]Episode, 1)
	s.Episodes[0] = NewEpisode(episode)

	return s
}

func NewListens(listens []*models.Listen, users []*models.User) Listens {
	s := Listens{}
	s.Listens = make([]Listen, len(listens))
	s.Users = make([]User, len(users))

	for i, listen := range listens {
		s.Listens[i] = NewListen(listen)
	}

	for i, user := range users {
		s.Users[i] = NewUser(user)
	}

	return s
}

type Favorite struct {
	*models.Favorite
}

type Favorites struct {
	Favorites []Favorite
	Users     []User
}

func NewFavorite(favorite *models.Favorite) Favorite {
	return Favorite{favorite}
}

func NewFavorites(favorites []*models.Favorite, users []*models.User) Favorites {
	s := Favorites{}
	s.Favorites = make([]Favorite, len(favorites))
	s.Users = make([]User, len(favorites))

	for i, favorite := range favorites {
		s.Favorites[i] = NewFavorite(favorite)
	}

	for i, user := range users {
		s.Users[i] = NewUser(user)
	}

	return s
}
