package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

var InputFolder = "./img/input"
var OutputFolder = "./img/output"

func configure() *twitter.Client {
	config := oauth1.NewConfig(consumerKey, consumerSecret)
	token := oauth1.NewToken(accessToken, accessSecret)
	// OAuth1 http.Client will automatically authorize Requests
	httpClient := config.Client(oauth1.NoContext, token)
	// Twitter client
	client := twitter.NewClient(httpClient)
	return client
}

func retweet(client *twitter.Client) {
	// Search tweets to retweet
	searchParams := &twitter.SearchTweetParams{
		Query:      "#golang",
		Count:      5,
		ResultType: "recent",
		Lang:       "en",
	}

	searchResult, _, _ := client.Search.Tweets(searchParams)

	// Retweet the search Results
	for _, tweet := range searchResult.Statuses {
		tweetID := tweet.ID
		client.Statuses.Retweet(tweetID, &twitter.StatusRetweetParams{})

		fmt.Printf("RETWEETED: %+v\n", tweet.Text)
	}
}

func downloadPhoto(url string) {
	response, e := http.Get(url)
	if e != nil {
		log.Fatal(e)
	}
	defer response.Body.Close()

	//open file for writing
	file, err := os.Create()
}

func main() {
	client := configure()
	retweet(client)
}
