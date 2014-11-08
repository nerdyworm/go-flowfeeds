package models

import (
	"errors"
	"time"

	"code.google.com/p/go.crypto/bcrypt"
)

var (
	PasswordCost = bcrypt.DefaultCost
	ErrNotFound  = errors.New("Record not found")
)

type Feed struct {
	Id          int64
	Title       string
	Description string
	Url         string
	Image       string
	Updated     time.Time
}

type Episode struct {
	Id             int64
	Feed           int64
	Guid           string
	Title          string
	Description    string
	Url            string
	Image          string
	Published      time.Time
	ListensCount   int  `db:"listens_count"`
	FavoritesCount int  `db:"favorites_count"`
	Favorited      bool `db:"-"`
	Listened       bool `db:"-"`
}

type Listen struct {
	Id      int64
	User    int64
	Episode int64
}

type Favorite struct {
	Id      int64
	User    int64
	Episode int64
}

type User struct {
	Id                int64
	Email             string
	EncryptedPassword []byte
}

func (user User) CheckPassword(password string) error {
	return bcrypt.CompareHashAndPassword(user.EncryptedPassword, []byte(password))
}

func (user *User) SetPassword(password string) (err error) {
	user.EncryptedPassword, err = bcrypt.GenerateFromPassword([]byte(password), PasswordCost)
	return err
}

func (user User) IsAnonymous() bool {
	return user.Id == 0
}

func NewUser(email, password string) *User {
	user := User{Email: email}
	user.SetPassword(password)
	return &user
}
