package github

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	ghapi "github.com/google/go-github/v44/github"
	"golang.org/x/oauth2"
)

func generateJWT(appID int, privateKey []byte) (string, error) {
	unixNow := time.Now().Unix()

	token := jwt.NewWithClaims(
		jwt.SigningMethodRS256,
		jwt.StandardClaims{
			IssuedAt:  unixNow - 60,
			ExpiresAt: unixNow + 300,
			Issuer:    fmt.Sprintf("%d", appID),
		},
	)

	tokenSigned, err := token.SignedString(privateKey)
	if err != nil {
		return "", fmt.Errorf("error signing JWT token: %w", err)
	}

	return tokenSigned, nil
}

func newGitHubClient(jwt string) *ghapi.Client {
	tokenSource := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: jwt},
	)
	oauthClient := oauth2.NewClient(context.Background(), tokenSource)
	return ghapi.NewClient(oauthClient)
}
