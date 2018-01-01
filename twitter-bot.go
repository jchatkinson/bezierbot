package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
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

func encodePhoto(file string) (base64String string, err error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Println(err)
	}
	base64String = base64.StdEncoding.EncodeToString(data)
	return
}

// uploadPhoto uploads the photo to twitter and returns the json response as a Media type
func uploadPhoto(base64String string) (mediaResponse Media, err error) {
	v := url.Values{}
	v.Set("media_data", base64String)
	queue := make(chan Query)
	queryQueue := queue
	responseCh := make(chan response)
	queryQueue <- Query{UploadBaseUrl + "/media/upload.json", v, &mediaResponse, _POST, responseCh}
	return mediaResponse, (<-responseCh).err
}

func tweetPhoto(client *twitter.Client, text string, file string) {
	// encode the photo
	base64string, _ := encodePhoto(file)

	// upload the photo
	media, err := uploadPhoto(base64string)
	if err != nil {
		println("error uploading photo")
		println(err)
	}

	// set the media id for the tweet
	var vs *twitter.StatusUpdateParams
	vs.MediaIds = make([]int64, media.MediaID)
	// send a tweet with the media id and log the result
	tweet, _, err := client.Statuses.Update(text, vs)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("TWEETED: %+v\n", tweet.Text)
	}
}

func downloadPhoto(url string) (uuid, imgpath string) {
	response, e := http.Get(url)
	if e != nil {
		log.Fatal(e)
	}
	defer response.Body.Close()

	//open file for writing
	uuid, err := newUUID()
	filename := uuid + ".jpg"
	imgpath = filepath.Join(InputFolder, filename)
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
	fmt.Println("Image Downloaded!")
	return
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
func processPhoto(inputFile, outputFile string, n, m int) {
	cmd := "primitive"
	args := []string{"-i", inputFile, "-o", outputFile, "-n", strconv.Itoa(n), "-m", strconv.Itoa(m)}
	if err := exec.Command(cmd, args...).Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println("Image Processed!")
}

func postNewPhoto() {
	client := configure()
	photourl := searchForPhoto()
	uuid, inputFile := downloadPhoto(photourl)
	outputFile := "./img/output/" + uuid + ".out.jpg"
	n := 200 + mathrand.Intn(450)
	mode := 1
	processPhoto(inputFile, outputFile, n, mode)
	tweettext := "n=" + strconv.Itoa(n) + " mode=" + strconv.Itoa(mode) + " (original: " + photourl + ")"
	tweetPhoto(client, tweettext, outputFile)
}
