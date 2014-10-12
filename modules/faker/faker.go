package faker

import (
	"fmt"
	"log"
	"math/rand"
	"os/exec"
	"strings"
	"time"

	"bitbucket.org/nerdyworm/go-flowfeeds/datastore"
	"bitbucket.org/nerdyworm/go-flowfeeds/models"
)

func Run() {
	log.Println("Faking data")

	featured, _, _, _, err := models.FeaturedEpisodes(models.User{}, models.ListOptions{})
	if err != nil {
		log.Fatal(err)
	}

	store := datastore.NewDatastore()

	rand.Seed(time.Now().UnixNano())

	users := createNewUsers()
	for _, user := range users {
		rand.Seed(time.Now().UnixNano())

		for _, f := range featured {
			if shouldListen() {
				models.CreateListen(user, f.Id)
			}

			if shouldFavorite() {
				store.Episodes.ToggleFavoriteForUser(&user, f.Id)
			}

			related, err := store.Episodes.Related(f.Id)
			if err != nil {
				log.Fatal(err)
			}

			for _, e := range related {
				if shouldListen() {
					models.CreateListen(user, e.Id)
				}

				if shouldFavorite() {
					store.Episodes.ToggleFavoriteForUser(&user, f.Id)
				}
			}
		}
	}
}

func shouldFavorite() bool {
	return rand.Intn(20) == 10
}

func shouldListen() bool {
	return rand.Intn(10) == 5
}

func createNewUsers() []models.User {
	users := []models.User{}

	count := rand.Intn(1000)
	for i := 0; i < count; i++ {
		output, err := exec.Command("uuidgen").Output()
		if err != nil {
			log.Fatal(err)
		}

		uuid := strings.ToLower(strings.TrimSpace(string(output)))
		user, err := models.CreateUser(fmt.Sprintf("%s@flowfeeds.com", uuid), uuid)
		if err != nil {
			log.Fatal(err)
		}

		users = append(users, user)
	}

	return users
}