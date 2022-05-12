package github

import (
	"context"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	ghapi "github.com/google/go-github/v44/github"
	"github.com/rs/zerolog/log"
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

	decodedPem, _ := pem.Decode(privateKey)
	if decodedPem == nil {
		return "", fmt.Errorf("unexpected empty decoded pem")
	}
	log.Info().Msgf("Decoded PEM of type %s", decodedPem.Type)

	rsaPrivateKey, err := x509.ParsePKCS1PrivateKey(decodedPem.Bytes)
	if err != nil {
		return "", fmt.Errorf("error parsing private key: %w", err)
	}

	tokenSigned, err := token.SignedString(rsaPrivateKey)
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
