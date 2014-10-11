package ctx

import (
	"bitbucket.org/nerdyworm/go-flowfeeds/datastore"
	"bitbucket.org/nerdyworm/go-flowfeeds/models"
)

type Context struct {
	User  models.User
	Store *datastore.Datastore
}
