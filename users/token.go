package users

import (
	"context"

	"github.com/dghubble/gologin/v2/oauth1"
	"github.com/dghubble/gologin/v2/twitter"
)

func userFromContext(ctx context.Context) *User {
	tu, _ := twitter.UserFromContext(ctx)
	accessToken, accessSecret, _ := oauth1.AccessTokenFromContext(ctx)
	return &User{
		AccessToken:  accessToken,
		AccessSecret: accessSecret,
		User:         tu,
	}
}
