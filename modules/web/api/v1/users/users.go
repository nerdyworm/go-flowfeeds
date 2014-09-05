package users

import (
	"encoding/json"
	"log"
	"net/http"

	"bitbucket.org/nerdyworm/go-flowfeeds/models"
	"bitbucket.org/nerdyworm/go-flowfeeds/modules/web/api/v1/serializers"
	"bitbucket.org/nerdyworm/go-flowfeeds/modules/web/ctx"
)

type CreateUserRequest struct {
	User struct {
		Email    string
		Password string
	}
}

func (r CreateUserRequest) Validate() (models.ValidationErrors, error) {
	validationErrors := models.NewValidationErrors()
	params := r.User

	if params.Email == "" {
		validationErrors.Add("Email", "can't be blank")
	} else {
		exists, err := models.UserExistsWithEmail(params.Email)
		if err != nil {
			log.Println("users.Create models.UserExistsWithEmail", err)
			return validationErrors, err
		}

		if exists {
			validationErrors.Add("Email", "already taken")
		}
	}

	if params.Password == "" {
		validationErrors.Add("Password", "can't be blank")
	}

	return validationErrors, nil
}

func Create(ctx ctx.Context, w http.ResponseWriter, r *http.Request) error {
	createUserRequest := CreateUserRequest{}
	err := json.NewDecoder(r.Body).Decode(&createUserRequest)
	if err != nil {
		log.Println("users.Create", err)
		return err
	}

	errors, err := createUserRequest.Validate()
	if err != nil {
		return err
	}

	if errors.Any() {
		return errors
	}

	params := createUserRequest.User

	user, err := models.CreateUser(params.Email, params.Password)
	if err != nil {
		log.Println("users.Create models.UserCreate", err)
		return err
	}

	serializer := serializers.ShowUser{
		serializers.User{
			Id:    user.Id,
			Email: user.Email,
		},
	}

	w.WriteHeader(http.StatusCreated)
	serializers.JSON(w, serializer)
	return nil
}
