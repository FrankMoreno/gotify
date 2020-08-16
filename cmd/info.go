package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"golang.org/x/oauth2"

	"github.com/spf13/cobra"
)

// var auth = spotify.NewAuthenticator(redirectURI, spotify.ScopeUserModifyPlaybackState, spotify.ScopeUserReadPlaybackState)

// infoCmd represents the info command
var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		token := readTokenFromFile()

		client := auth.NewClient(&token)
		user, err := client.CurrentUser()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("You are logged in as:", user.ID)
		currPlaying, _ := client.PlayerCurrentlyPlaying()
		log.Println(currPlaying.Item)
		log.Printf("%v", currPlaying)
		client.Pause()
		log.Println(currPlaying.Item)
	},
}

func init() {
	rootCmd.AddCommand(infoCmd)
}

func readTokenFromFile() oauth2.Token {
	homedir, _ := os.UserHomeDir()
	filePath := fmt.Sprintf("%s/.spotify/credentials", homedir)

	var token oauth2.Token

	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Println(err.Error())
	}

	json.Unmarshal(data, &token)
	return token
}
