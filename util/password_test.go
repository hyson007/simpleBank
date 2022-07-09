package util

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPassword(t *testing.T) {
	password := RandomString(6)
	hashPwd, err := HashPassword(password)
	require.NoError(t, err)
	require.NoError(t, CheckPassword(password, hashPwd))

	wrongPwd := RandomString(6)
	wrongHash, err := HashPassword(wrongPwd)
	require.NoError(t, err)
	require.Error(t, CheckPassword(password, wrongHash))
}
