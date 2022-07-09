package token

import (
	"testing"
	"time"

	"github.com/hyson007/simpleBank/util"
	"github.com/stretchr/testify/require"
)

func TestPasetoToken(t *testing.T) {

	//Generate token
	maker, err := NewPasetoMaker(util.RandomString(32))
	require.NoError(t, err)
	username := util.RandomOwner()
	duration := time.Minute
	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)
	token, err := maker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	//Verify token
	payload, err := maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.NotZero(t, payload.ID)
	require.Equal(t, username, payload.Username)
	require.WithinDuration(t, issuedAt, payload.IssueAt, time.Second)
	require.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)
}

func TestExpiredPasetoToken(t *testing.T) {
	//Generate an expired token
	maker, err := NewPasetoMaker(util.RandomString(32))
	require.NoError(t, err)

	token, err := maker.CreateToken(util.RandomOwner(), -time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrExpiredToken.Error())
	require.Nil(t, payload)
}

// func TestInValidPasetoTokenAlgNone(t *testing.T) {
// 	// this token implement a different signing Method than
// 	// the one we define in our implementation which is SigningMethodHS256
// 	payload, err := NewPayload(util.RandomOwner(), time.Minute)
// 	require.NoError(t, err)

// 	jwtToken := jwt.NewWithClaims(jwt.SigningMethodNone, payload)

// 	// this libary only allows using this constant to sign a none siginging method
// 	token, err := jwtToken.SignedString(jwt.UnsafeAllowNoneSignatureType)
// 	require.NoError(t, err)

// 	maker, err := NewJWTMaker(util.RandomString(32))
// 	require.NoError(t, err)

// 	payload, err = maker.VerifyToken(token)
// 	// t.Log(token)
// 	// t.Log(payload)
// 	// t.Log(err)
// 	require.Error(t, err)
// 	require.EqualError(t, err, ErrInvalidToken.Error())
// 	require.Nil(t, payload)
// }
