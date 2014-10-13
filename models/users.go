package models

import (
	"database/sql"

	"code.google.com/p/go.crypto/bcrypt"
)

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

const USER_EMAIL_EXISTS_QUERY = "select exists (select true from users where lower(email) = lower($1));"

func UserExistsWithEmail(email string) (exists bool, err error) {
	row := x.DB().QueryRow(USER_EMAIL_EXISTS_QUERY, email)
	err = row.Scan(&exists)
	return
}

func NewUser(email, password string) User {
	user := User{Email: email}
	user.SetPassword(password)
	return user
}

func FindUserForSignin(email string) (User, error) {
	user := User{}

	row := x.DB().QueryRow("select id, email, encrypted_password from users where lower(email) = lower($1)", email)

	err := row.Scan(
		&user.Id,
		&user.Email,
		&user.EncryptedPassword,
	)

	if err == sql.ErrNoRows {
		return user, ErrNotFound
	}

	return user, err
}
