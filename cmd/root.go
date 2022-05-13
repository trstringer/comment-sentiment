/*
Package cmd is the entry point for the main command.

Copyright © 2022 Thomas Stringer <thomas@trstringer.com>
*/
package cmd

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	gh "github.com/trstringer/comment-sentiment/pkg/github"
	"github.com/trstringer/comment-sentiment/pkg/sentimentanalyzer/azure"
	"github.com/trstringer/comment-sentiment/pkg/version"
)

var (
	port              int
	languageKeyFile   string
	languageKey       string
	webhookSecret     []byte
	webhookSecretFile string
	languageEndpoint  string
	appID             int
	appKeyFile        string
	appKey            []byte
	showVersion       bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "comment-sentiment",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		if showVersion {
			fmt.Println(version.Version)
			os.Exit(0)
		}

		if languageKeyFile == "" {
			fmt.Println("Required parameter --language-key not supplied")
			os.Exit(1)
		}
		if webhookSecretFile == "" {
			fmt.Println("Required parameter --webhook-secretfile not supplied")
			os.Exit(1)
		}
		if languageEndpoint == "" {
			fmt.Println("Required parameter --language-endpoint not supplied")
			os.Exit(1)
		}
		if appKeyFile == "" {
			fmt.Println("Required parameter --app-keyfile not supplied")
			os.Exit(1)
		}

		languageKeyFilePath, err := filepath.Abs(languageKeyFile)
		if err != nil {
			fmt.Printf("Error getting file path for language key: %v\n", err)
			os.Exit(1)
		}
		languageKeyBytes, err := ioutil.ReadFile(languageKeyFilePath)
		if err != nil {
			fmt.Printf("Error reading language key file: %v\n", err)
			os.Exit(1)
		}
		languageKey = string(languageKeyBytes)

		webhookSecretFilePath, err := filepath.Abs(webhookSecretFile)
		if err != nil {
			fmt.Printf("Error getting file path for webhook secret: %v\n", err)
			os.Exit(1)
		}
		webhookSecret, err = ioutil.ReadFile(webhookSecretFilePath)
		if err != nil {
			fmt.Printf("Error reading webhook secret file: %v\n", err)
			os.Exit(1)
		}

		if appID <= 0 {
			fmt.Println("Required parameter --app-id not supplied or incorrect value")
			os.Exit(1)
		}

		appKey, err = ioutil.ReadFile(appKeyFile)
		if err != nil {
			fmt.Printf("Error reading app key file: %v\n", err)
			os.Exit(1)
		}

		startServer(port)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().IntVarP(&port, "port", "p", 8080, "port that the server should be listening on")
	rootCmd.Flags().StringVarP(&languageKeyFile, "language-keyfile", "l", "", "cognitive services language key file path")
	rootCmd.Flags().StringVarP(&languageEndpoint, "language-endpoint", "e", "", "cognitive services language endpoint")
	rootCmd.Flags().StringVarP(&webhookSecretFile, "webhook-secretfile", "w", "", "file storing the webhook secret")
	rootCmd.Flags().IntVar(&appID, "app-id", 0, "GitHub App ID")
	rootCmd.Flags().StringVarP(&appKeyFile, "app-keyfile", "a", "", "GitHub App key file path")
	rootCmd.Flags().BoolVarP(&showVersion, "version", "v", false, "list the version")
}

func isRequestValid(signature string, body []byte, secretKey []byte) (bool, string) {
	hash := hmac.New(sha256.New, secretKey)
	hash.Write(body)
	calculatedHashBytes := hash.Sum(nil)
	calculatedHash := fmt.Sprintf("sha256=%x", string(calculatedHashBytes))
	return calculatedHash == signature, calculatedHash
}

func handleSentimentRequest(resp http.ResponseWriter, req *http.Request) {
	log.Info().Msg("Received sentiment handle request")

	if req.Method != http.MethodPost {
		resp.WriteHeader(http.StatusBadRequest)
		resp.Write([]byte("Only POST supported"))
		log.Warn().Msgf("Received non-POST request %s for sentiment handler", req.Method)
		return
	}

	body := req.Body
	if body == nil {
		resp.WriteHeader(http.StatusBadRequest)
		resp.Write([]byte("Missing request body"))
		log.Warn().Msg("Request body for sentiment handler is missing")
		return
	}
	defer req.Body.Close()

	log.Debug().Msg("Reading body of sentiment analysis request")
	payloadRaw, err := ioutil.ReadAll(body)
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte("Error reading body of request"))
		log.Error().Err(err).Msg("Error reading body")
		return
	}

	githubSignature := req.Header.Get("X-Hub-Signature-256")
	requestIsValid, computedHash := isRequestValid(githubSignature, payloadRaw, webhookSecret)
	if !requestIsValid {
		resp.WriteHeader(http.StatusUnauthorized)
		resp.Write([]byte("Unauthorized access denied"))
		log.Warn().Msgf(
			"Mismatched signature from request (%s) to computed (%s)",
			githubSignature,
			computedHash,
		)
		return
	}

	commentPayload := gh.CommentPayload{}
	log.Debug().Msg("Unmarshalling payload")
	if err = json.Unmarshal(payloadRaw, &commentPayload); err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte("Error unmarshalling payload"))
		log.Error().Err(err).Msg("Error unmarshalling paylog")
		return
	}

	if commentPayload.Comment.CommentUser.Login != commentPayload.Sender.Login {
		log.Debug().Msgf(
			"Sender %s is not the comment login %s",
			commentPayload.Sender.Login,
			commentPayload.Comment.CommentUser.Login,
		)
		_, err := resp.Write(nil)
		if err != nil {
			log.Error().Err(err).Msg("Error responding with nil body")
		}
		return
	}

	log.Debug().Msgf("Creating new GitHub client for repo owner %s", commentPayload.Repository.Owner.Login)
	client, err := gh.NewInstallationGitHubClient(appID, appKey, commentPayload.Repository.Owner)
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte("Error creating GitHub client"))
		log.Error().Err(err).Msg("Error creating github client")
		return
	}

	log.Debug().Msg("Creating new sentiment service and analyzing")
	sentimentSvc := azure.NewSentimentService(languageEndpoint, languageKey)
	analysis, err := sentimentSvc.AnalyzeSentiment(commentPayload.Comment.Body)
	log.Debug().Msgf("Analysis result: %s", analysis.Sentiment.String())
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte("Error getting sentiment"))
		log.Error().Err(err).Msg("Error getting sentiment analysis")
		return
	}

	log.Debug().Msg("Updating comment")
	updatedComment, err := gh.UpdateCommentWithSentiment(commentPayload.Comment.Body, *analysis)
	if err := commentPayload.UpdateComment(client, updatedComment); err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte(fmt.Sprintf("%v", err)))
		log.Error().Err(err).Msg("Error updating comment")
		return
	}

	log.Info().Msg("Successfully processed request")
	resp.Write([]byte("success"))
}

func handleManualSentimentRequest(resp http.ResponseWriter, req *http.Request) {
	fmt.Println("Received manual request to handle sentiment")

	body := req.Body
	if body == nil {
		resp.WriteHeader(http.StatusBadRequest)
		resp.Write([]byte("Missing request body"))
		return
	}
	defer req.Body.Close()

	commentDataRaw, err := ioutil.ReadAll(body)
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte("Error reading body of request"))
		return
	}
	commentData := string(commentDataRaw)

	sentimentSvc := azure.NewSentimentService(languageEndpoint, languageKey)
	analysis, err := sentimentSvc.AnalyzeSentiment(commentData)

	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte("Error getting sentiment"))
		return
	}

	if analysis == nil {
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte("Unexpectedly no analysis returned"))
		return
	}

	resp.Write([]byte(fmt.Sprintf(
		"Analysis: %s - Confidence: %.2f",
		analysis.Sentiment,
		analysis.Confidence,
	)))
}

func startServer(port int) {
	log.Info().Msgf("Starting server on port %d", port)

	http.HandleFunc("/", handleSentimentRequest)
	http.HandleFunc("/manual", handleManualSentimentRequest)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		log.Fatal().Msgf("Error creating server: %v", err)
	}
}
