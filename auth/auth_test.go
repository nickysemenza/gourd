package auth

import (
	"testing"

	"github.com/stretchr/testify/require"
	gauth "google.golang.org/api/oauth2/v2"
)

func TestAuth_GetJWT(t *testing.T) {
	user := gauth.Userinfo{Email: "foo@bar.com"}
	a1, _ := New("key1")
	a2, _ := New("key2")

	jwt, err := a1.GetJWT(&user)
	require.NoError(t, err)
	require.NotEmpty(t, jwt)

	claims, err := a1.ParseJWT(jwt)
	require.NoError(t, err)
	require.EqualValues(t, user.Email, claims.User.Email)

	_, err = a2.ParseJWT(jwt)
	require.Error(t, err)
}
