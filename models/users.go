package models

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"code.google.com/p/go.crypto/bcrypt"
)

var (
	ErrNotFound = errors.New("Record not found")
)

type RecordNotFound struct {
	error
	Message string
}

func (r RecordNotFound) Error() string {
	return r.Message
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

func CreateUser(email string, password string) (User, error) {
	user := NewUser(email, password)

	row := x.DB().QueryRow("insert into users (email, encrypted_password) values($1,$2) returning id",
		user.Email, user.EncryptedPassword)

	err := row.Scan(&user.Id)

	return user, err
}

func FindUserById(id int64) (User, error) {
	user := User{}

	row := x.DB().QueryRow("select id, email from users where id = $1", id)

	err := row.Scan(
		&user.Id,
		&user.Email,
	)

	if err == sql.ErrNoRows {
		return user, ErrNotFound
	}

	return user, err
}

func FindUserByIds(ids []int64) ([]User, error) {
	users := []User{}

	if len(ids) == 0 {
		return users, nil
	}

	stringIds := []string{}
	for i := range ids {
		stringIds = append(stringIds, fmt.Sprintf("%d", ids[i]))
	}

	query := fmt.Sprintf("select id, email from users where id in (%s)", strings.Join(stringIds, ","))
	rows, err := x.DB().Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		user := User{}

		err := rows.Scan(
			&user.Id,
			&user.Email,
		)

		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
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
