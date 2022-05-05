/*
Package cmd is the entry point for the main command.

Copyright © 2022 Thomas Stringer <thomas@trstringer.com>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	gh "github.com/trstringer/comment-sentiment/pkg/github"
	"github.com/trstringer/comment-sentiment/pkg/sentimentanalyzer/azure"
	"github.com/trstringer/comment-sentiment/pkg/version"
)

var (
	port             int
	languageKeyFile  string
	languageKey      string
	languageEndpoint string
	appID            int
	appKeyFile       string
	appKey           []byte
	showVersion      bool
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
		if languageEndpoint == "" {
			fmt.Println("Required parameter --language-endpoint not supplied")
			os.Exit(1)
		}

		filePath, err := filepath.Abs(languageKeyFile)
		if err != nil {
			fmt.Printf("Error getting file path: %v\n", err)
			os.Exit(1)
		}

		languageKeyBytes, err := ioutil.ReadFile(filePath)
		if err != nil {
			fmt.Printf("Error reading key file: %v\n", err)
			os.Exit(1)
		}
		languageKey = string(languageKeyBytes)

		if appID <= 0 {
			fmt.Println("Incorrect GitHub App ID")
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
	rootCmd.Flags().IntVar(&appID, "app-id", 0, "GitHub App ID")
	rootCmd.Flags().StringVarP(&appKeyFile, "app-keyfile", "a", "", "GitHub App key file path")
	rootCmd.Flags().BoolVarP(&showVersion, "version", "v", false, "list the version")
}

func handleSentimentRequest(resp http.ResponseWriter, req *http.Request) {
	fmt.Println("Received request to handle sentiment")

	if req.Method != http.MethodPost {
		resp.WriteHeader(http.StatusBadRequest)
		resp.Write([]byte("Only POST supported"))
		return
	}

	body := req.Body
	if body == nil {
		resp.WriteHeader(http.StatusBadRequest)
		resp.Write([]byte("Missing request body"))
		return
	}
	defer req.Body.Close()

	payloadRaw, err := ioutil.ReadAll(body)
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte("Error reading body of request"))
		return
	}
	commentPayload := gh.CommentPayload{}
	if err = json.Unmarshal(payloadRaw, &commentPayload); err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte("Error unmarshalling payload"))
		return
	}

	client, err := gh.NewInstallationGitHubClient(appID, appKey, commentPayload.Repository.Owner)
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte("Error creating GitHub client"))
		return
	}
	sentimentSvc := azure.NewSentimentService(languageEndpoint, languageKey)
	analysis, err := sentimentSvc.AnalyzeSentiment(commentPayload.Comment.Body)
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte("Error getting sentiment"))
		return
	}
	updatedComment, err := gh.UpdateCommentWithSentiment(commentPayload.Comment.Body, *analysis)
	if err := commentPayload.UpdateComment(client, updatedComment); err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte(fmt.Sprintf("%v", err)))
		return
	}

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
	fmt.Printf("Starting server on port %d\n", port)

	http.HandleFunc("/", handleSentimentRequest)
	http.HandleFunc("/manual", handleManualSentimentRequest)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
