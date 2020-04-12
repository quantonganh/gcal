package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/spf13/viper"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
)

func getClient(config *oauth2.Config) *http.Client {
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}

	return config.Client(context.Background(), tok)
}

func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, fmt.Errorf("Unable to open file %s: %v", file, err)
	}
	defer f.Close()

	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func newCalendar() (*calendar.Service, error) {
	b, err := ioutil.ReadFile(viper.ConfigFileUsed())
	if err != nil {
		return nil, fmt.Errorf("unable to read client secret file: %v", err)
	}

	config, err := google.ConfigFromJSON(b, calendar.CalendarReadonlyScope)
	if err != nil {
		return nil, fmt.Errorf("unable to parse client secret file from config: %v", err)
	}
	client := getClient(config)

	srv, err := calendar.New(client)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve calendar client: %v", err)
	}

	return srv, nil
}

func getCalendarID(srv *calendar.Service, calendarSummary string) (string, error) {
	if calendarSummary == "" {
		return defaultCalendar, nil
	}

	var calendarID string
	for {
		calendarList, err := srv.CalendarList.List().Do()
		if err != nil {
			return "", fmt.Errorf("Cannot get calendars list: %v", err)
		}
		for _, item := range calendarList.Items {
			if item.Summary == calendarSummary {
				calendarID = item.Id
				break
			}
		}
		pageToken := calendarList.NextPageToken
		if pageToken == "" {
			break
		}
	}
	return calendarID, nil
}
