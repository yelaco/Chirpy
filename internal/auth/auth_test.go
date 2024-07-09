package auth

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRefreshToken(t *testing.T) {
	token, err := CreateRefreshToken()
	require.NoError(t, err)
	require.NotEmpty(t, token)

	fmt.Println(token)
}

func TestAccessToken(t *testing.T) {
	userIDs := []int{1}

	secret := "G5rNeKgkNGQWNp90ekGRJq8BQgskiIXUDxuJoBmoho2RInuVlvlfIaHlDa26WCA+0c8QSjmed2PI9nJTNDOYWw=="
	tokenString, err := CreateAccessToken(userIDs[0], 60, secret)
	require.NoError(t, err)
	require.NotEmpty(t, tokenString)

	userID, err := ValidateAccessToken(tokenString, secret)
	require.NoError(t, err)
	require.NotEmpty(t, userID)

	require.Equal(t, userID, fmt.Sprint(userIDs[0]))
}
