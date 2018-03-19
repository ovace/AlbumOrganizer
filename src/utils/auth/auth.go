// Auth handles the OAuth2 authentication and token exchange.
// It is not much more than the example code from the oauth2
// package.
//
// https://godoc.org/golang.org/x/oauth2
package auth

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const TokenFile = ".gdriveutil"

// Your credentials should be obtained from the Google Developer Console
// (https://console.developers.google.com).
var conf = &oauth2.Config{
	ClientID:     "YOUR_CLIENT_ID",
	ClientSecret: "YOUR_SECRET",
	RedirectURL:  "urn:ietf:wg:oauth:2.0:oob",
	Scopes: []string{
		"https://www.googleapis.com/auth/drive",
	},
	Endpoint: google.Endpoint,
}

// writeToken stores tok to the file TokenFile, which serves as a token cache.
// If the token exists, then the process will refresh the token
// and a new token will not have to be generated.
func writeToken(tok *oauth2.Token) error {
	b, err := json.Marshal(tok)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(TokenFile, b, 0644)
}

// readToken attempts to materialize the token stored in TokenFile.
func readToken() (*oauth2.Token, error) {
	dat, err := ioutil.ReadFile(TokenFile)
	if err != nil {
		return nil, err
	}
	var tok oauth2.Token
	err = json.Unmarshal(dat, &tok)
	if err != nil {
		return nil, err
	}
	return &tok, nil
}

// DoAuth attempts to retrieve the token and, if that is not possible,
// it will generate a new token, and then return a http.Client with the
// OAuth2 context for the token.
func DoAuth() (*http.Client, error) {
	tok, err := readToken()
	if err != nil {
		// Redirect user to Google's consent page to ask for permission
		// for the scopes specified above.
		url := conf.AuthCodeURL("state")
		fmt.Printf("Visit the URL for the auth dialog: %v\n", url)

		fmt.Printf("Enter verification code: ")
		var code string
		fmt.Scanln(&code)

		// Handle the exchange code to initiate a transport.
		tok, err := conf.Exchange(oauth2.NoContext, code)
		if err != nil {
			return nil, err
		}

		// Save Token
		if err := writeToken(tok); err != nil {
			return nil, err
		}
	}

	return conf.Client(oauth2.NoContext, tok), nil
}
