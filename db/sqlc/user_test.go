package db

import (
	"context"
	"testing"
	"time"

	"github.com/hyson007/simpleBank/util"
	"github.com/stretchr/testify/require"
)

//in order to test CreateAccount, we have to setup the database connection
//first, the right place to do that is in main_test.go

func createRandomUser(t *testing.T) User {
	hashPwd, err := util.HashPassword(util.RandomString(6))
	require.NoError(t, err)

	arg := CreateUserParams{
		Username:     util.RandomOwner(),
		HashPassword: hashPwd,
		FullName:     util.RandomOwner(),
		Email:        util.RandomEmail(),
	}

	user, err := testQuery.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.HashPassword, user.HashPassword)
	require.Equal(t, arg.FullName, user.FullName)
	require.NotZero(t, user.CreatedAt)

	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGeUser(t *testing.T) {
	user1 := createRandomUser(t)
	user2, err := testQuery.GetUser(context.Background(), user1.Username)
	require.NoError(t, err)
	require.NotEmpty(t, user2)
	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.HashPassword, user2.HashPassword)
	require.Equal(t, user1.FullName, user2.FullName)
	require.Equal(t, user1.Email, user2.Email)
	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)
}
