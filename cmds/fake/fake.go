package fake

import (
	"fmt"
	"log"
	"math/rand"
	"os/exec"
	"strings"
	"time"

	"github.com/nerdyworm/go-flowfeeds/datastore"
	"github.com/nerdyworm/go-flowfeeds/models"
)

var (
	USERS = 10
	store *datastore.Datastore
)

func Run() {
	store = datastore.NewDatastore()

	page := 0
	users := createNewUsers()
	rand.Seed(time.Now().UnixNano())

	for {
		page += 1
		episodes, _, err := store.Episodes.ListFor(&models.User{}, datastore.EpisodeListOptions{
			ListOptions: datastore.ListOptions{PerPage: 1000, Page: page},
		})

		if err != nil {
			log.Fatal(err)
		}

		if len(episodes.Episodes) == 0 {
			break
		}

		for _, user := range users {
			for _, f := range episodes.Episodes {
				if shouldListen() {
					store.Listens.Create(&user, f.Id)
				}

				if shouldFavorite() {
					store.Episodes.ToggleFavoriteForUser(&user, f.Id)
				}
			}
		}
	}
}

func shouldFavorite() bool {
	numbers := []bool{true, false, false, false}
	return numbers[rand.Intn(len(numbers))]
}

func shouldListen() bool {
	numbers := []bool{true, false}
	return numbers[rand.Intn(len(numbers))]
}

func createNewUsers() []models.User {
	users := []models.User{}

	for i := 0; i < USERS; i++ {
		output, err := exec.Command("uuidgen").Output()
		if err != nil {
			log.Fatal(err)
		}

		uuid := strings.ToLower(strings.TrimSpace(string(output)))
		email := fmt.Sprintf("%s@flowfeeds.com", uuid)

		user := models.User{Email: email, EncryptedPassword: []byte(uuid)}
		err = store.Users.Insert(&user)
		if err != nil {
			log.Fatal(err)
		}

		users = append(users, user)
	}

	return users
}
