package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"

	pkce "github.com/frankmoreno/oauth2-pkce"
	"github.com/spf13/cobra"
	"github.com/zmb3/spotify"
	"golang.org/x/oauth2"
)

const (
	redirectURI = "http://localhost:8080/callback"
)

var (
	auth            = spotify.NewAuthenticator(redirectURI, spotify.ScopeUserModifyPlaybackState, spotify.ScopeUserReadPlaybackState, spotify.ScopeUserModifyPlaybackState)
	codeVerifier, _ = pkce.GenerateCodeVerifier(64, true)
	codeChallenge   = pkce.GenerateCodeChallenge(codeVerifier, true)
	state           = "abc123"
	ch              = make(chan *spotify.Client)
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Grant temp credentials to CLI",
	Run: func(cmd *cobra.Command, args []string) {
		login()
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
}

func login() {
	// first start an HTTP server
	http.HandleFunc("/callback", completeAuth)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Got request for:", r.URL.String())
	})
	go http.ListenAndServe(":8080", nil)

	url := auth.AuthURLWithOpts(state,
		oauth2.SetAuthURLParam("code_challenge_method", "S256"),
		oauth2.SetAuthURLParam("code_challenge", codeChallenge),
	)

	openbrowser(url)

	// wait for auth to complete
	client := <-ch

	// use the client to make calls that require authorization
	user, err := client.CurrentUser()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("You are logged in as:", user.ID)
}

func completeAuth(w http.ResponseWriter, r *http.Request) {
	tok, err := auth.TokenWithOpts(state, r,
		oauth2.SetAuthURLParam("code_verifier", codeVerifier))
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusForbidden)
		log.Fatal(err)
	}
	if st := r.FormValue("state"); st != state {
		http.NotFound(w, r)
		log.Fatalf("State mismatch: %s != %s\n", st, state)
	}
	// use the token to get an authenticated client
	client := auth.NewClient(tok)
	fmt.Fprintf(w, "Login Completed!")
	ch <- &client

	writeTokenToFile(tok)
}

func openbrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Fatal(err)
	}

}

func writeTokenToFile(token *oauth2.Token) {
	homedir, _ := os.UserHomeDir()
	filePath := fmt.Sprintf("%s/.spotify/credentials", homedir)

	data, err := json.MarshalIndent(token, "", "")
	if err != nil {
		log.Println(err.Error())
		return
	}

	err = ioutil.WriteFile(filePath, data, 0644)
	if err != nil {
		log.Println("CAN'T WRITE")
		log.Println(err.Error())
	}
}
