package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/user"
	"path/filepath"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/people/v1"
)

func getOAuthClient() *http.Client {
	c, err := ioutil.ReadFile("client_id.json")
	if err != nil {
		log.Fatalln("Unable to read client secret file:", err)
	}

	config, err := google.ConfigFromJSON(c, people.ContactsScope)
	if err != nil {
		log.Fatalln("Unable to parse client secret file to config:", err)
	}

	cacheFile, err := tokenCacheFile()
	if err != nil {
		log.Fatalln("Unable to get path to cached credential file:", err)
	}

	tok, err := tokenFromFile(cacheFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(cacheFile, tok)
	}

	ctx := context.Background()
	return config.Client(ctx, tok)
}

func tokenCacheFile() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}

	tokenCacheDir := filepath.Join(usr.HomeDir, ".credentials")
	os.MkdirAll(tokenCacheDir, 0700)

	return filepath.Join(tokenCacheDir, "contact-cleaner.json"), err
}

func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	t := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(t)

	return t, err
}

func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to following URL for auth:\n%v\n", authURL)

	var code string
	if _, err := fmt.Scan(&code); err != nil {
		log.Fatalln("Unable to read auth code:", err)
	}

	tok, err := config.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Fatalln("Unable to retrieve token from web:", err)
	}

	return tok
}

func saveToken(file string, token *oauth2.Token) {
	fmt.Println("Saving credentials file to ", file)

	f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalln("Unable to cache OAuth token:", err)
	}
	defer f.Close()

	json.NewEncoder(f).Encode(token)
}
