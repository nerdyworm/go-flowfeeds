package users

import (
	"encoding/json"
	"log"
	"net/http"

	"bitbucket.org/nerdyworm/go-flowfeeds/models"
	"bitbucket.org/nerdyworm/go-flowfeeds/modules/web/api/v1/serializers"
)

type CreateUserRequest struct {
	User struct {
		Email    string
		Password string
	}
}

func Create(w http.ResponseWriter, r *http.Request) {
	log.Println("Hello create")

	createUserRequest := CreateUserRequest{}
	err := json.NewDecoder(r.Body).Decode(&createUserRequest)
	if err != nil {
		log.Println("users.Create", err)
		return
	}

	validationErrors := models.NewValidationErrors()
	params := createUserRequest.User

	if params.Email == "" {
		validationErrors.Add("Email", "can't be blank")
	} else {
		exists, err := models.UserExistsWithEmail(params.Email)
		if err != nil {
			log.Println("users.Create models.UserExistsWithEmail", err)
			return
		}

		if exists {
			validationErrors.Add("Email", "already taken")
		}
	}

	if params.Password == "" {
		validationErrors.Add("Password", "can't be blank")
	}

	if validationErrors.Any() {
		w.WriteHeader(422)
		serializers.JSON(w, validationErrors)
		return
	}

	user, err := models.CreateUser(params.Email, params.Password)
	if err != nil {
		log.Println("users.Create models.UserCreate", err)
		return
	}

	serializer := serializers.ShowUser{
		serializers.User{
			Id:    user.Id,
			Email: user.Email,
		},
	}

	w.WriteHeader(http.StatusCreated)
	serializers.JSON(w, serializer)
}
