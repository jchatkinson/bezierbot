package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

import mathrand "math/rand"

const (
	_GET          = iota
	_POST         = iota
	_DELETE       = iota
	_PUT          = iota
	UploadBaseUrl = "https://upload.twitter.com/1.1"
	InputFolder   = "./img/input/"
	OutputFolder  = "./img/output/"
	BingURL       = "https://www.bing.com/HPImageArchive.aspx?format=js&idx=0&n=1&mkt=en-US"
)

func configure() *twitter.Client {
	config := oauth1.NewConfig(ConsumerKey, ConsumerSecret)
	token := oauth1.NewToken(AccessToken, AccessSecret)
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

func postNewPhoto() {
	photourl := searchForPhoto()
	uuid, inputFile := downloadPhoto(photourl)
	outputFile := "./img/output/" + uuid + "-out.jpg"
	n := 200 + mathrand.Intn(300)
	modes := []int{1, 4, 6}
	m := 0
	processPhoto(inputFile, outputFile, n, modes[m])

	// now post it on twitter
	tweettext := "n=" + strconv.Itoa(n) + " mode=" + strconv.Itoa(modes[m]) + " (original: " + photourl + ")"
	client := configure()

	media, _, err := client.Media.UploadFile(outputFile)
	if err != nil {
		log.Fatalf("UploadFile -> %v", err)
	}
	var params = twitter.StatusUpdateParams{MediaIds: []int64{media.MediaID}}

	_, _, err = client.Statuses.Update(tweettext, &params)
	if err != nil {
		log.Fatalf("Statuses.Update -> %v", err)
	} else {
		fmt.Printf("TWEETED: %+v\n", tweettext)
	}
}
