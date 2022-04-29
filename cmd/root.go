/*
Package cmd is the entry point for the main command.

Copyright © 2022 Thomas Stringer <thomas@trstringer.com>
*/
package cmd

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

var (
	port            int
	languageKeyFile string
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
		if languageKeyFile == "" {
			fmt.Println("Required parameter --language-key not supplied")
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
	rootCmd.Flags().StringVarP(&languageKeyFile, "language-key", "k", "", "cognitive services language key")
}

func handleSentimentRequest(resp http.ResponseWriter, req *http.Request) {
	fmt.Println("Received request to handle sentiment")

	resp.Write([]byte("hello galaxy"))
}

func startServer(port int) {
	fmt.Printf("Starting server on port %d\n", port)

	http.HandleFunc("/", handleSentimentRequest)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
