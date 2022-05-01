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

func newGitHubClient(token string) *ghapi.Client {
	tokenSource := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	oauthClient := oauth2.NewClient(context.Background(), tokenSource)
	return ghapi.NewClient(oauthClient)
}

// NewInstallationGitHubClient creates a new client for the installation.
func NewInstallationGitHubClient(appID int, privateKey []byte, repoOwner RepositoryOwner) (*ghapi.Client, error) {
	jwt, err := generateJWT(appID, privateKey)
	if err != nil {
		return nil, fmt.Errorf("error generating JWT: %w", err)
	}

	client := newGitHubClient(jwt)
	installations, _, err := client.Apps.ListInstallations(context.Background(), nil)
	if err != nil {
		return nil, fmt.Errorf("error getting app installations: %w", err)
	}

	var installationID int64 = -1
	for _, installation := range installations {
		if installation.GetAccount().GetLogin() == repoOwner.Login {
			installationID = installation.GetID()
		}
	}

	if installationID < 0 {
		return nil, fmt.Errorf("unable to find app installation")
	}

	installationToken, _, err := client.Apps.CreateInstallationToken(
		context.Background(),
		installationID,
		nil,
	)

	return newGitHubClient(installationToken.GetToken()), nil
}
