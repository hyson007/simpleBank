package api

import (
	"testing"

	db "github.com/hyson007/simpleBank/db/sqlc"
	"github.com/hyson007/simpleBank/util"
)

func randomUser(t *testing.T) db.User {
	return db.User{
		Username:     util.RandomString(6),
		HashPassword: "secret",
		FullName:     util.RandomString(6),
		Email:        util.RandomEmail(),
	}
}
