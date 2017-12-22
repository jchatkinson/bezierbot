package main

import (
	"crypto/rand"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

// InputFolder is the location for saved input photos
var InputFolder = "./img/input/"

// OutputFolder is the location for output photos
var OutputFolder = "./img/output/"

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

func downloadPhoto(url string) {
	response, e := http.Get(url)
	if e != nil {
		log.Fatal(e)
	}
	defer response.Body.Close()

	//open file for writing
	//uuid, err := newUUID()
	filename := "test.jpg"
	imgpath := filepath.Join(InputFolder, filename)
	file, err := os.Create(imgpath)
	if err != nil {
		log.Fatal(err)
	}

	// Use io.Copy to just dump the response body to the file. This supports huge files
	_, err = io.Copy(file, response.Body)
	if err != nil {
		log.Fatal(err)
	}
	file.Close()
	fmt.Println("Success!")
}

// newUUID generates a random UUID according to RFC 4122
func newUUID() (string, error) {
	uuid := make([]byte, 16)
	n, err := io.ReadFull(rand.Reader, uuid)
	if n != len(uuid) || err != nil {
		return "", err
	}
	// variant bits; see section 4.1.1
	uuid[8] = uuid[8]&^0xc0 | 0x80
	// version 4 (pseudo-random); see section 4.1.3
	uuid[6] = uuid[6]&^0xf0 | 0x40
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:]), nil
}

// processPhoto uses primative to process the photo
func processPhoto() {

}

func main() {
	// client := configure()
	downloadPhoto("https://blog.golang.org/gopher/header.jpg")
	processPhoto()

	primative
}
