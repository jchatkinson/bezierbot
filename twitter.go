package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
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

func encodePhoto(file string) (base64String string, err error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Println(err)
	}
	base64String = base64.StdEncoding.EncodeToString(data)
	return
}

// uploadPhoto uploads the photo to twitter and returns the json response as a Media type
func uploadPhoto(base64String string) (media Media, err error) {
	v := url.Values{}
	v.Set("media_data", base64String)
	var mediaResponse Media
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

func getAndTweetPhoto() {
	client := configure()
	url := searchForPhoto()
	uuid, inputFile := downloadPhoto(url)
	outputFile := "./img/output/" + uuid + ".out.jpg"
	n := 50 + mathrand.Intn(450)
	mode := 1
	processPhoto(inputFile, outputFile, n, mode)
	tweettext := "n=" + strconv.Itoa(n) + " mode=" + strconv.Itoa(mode) + " (original: " + url + ")"
	tweetPhoto(client, tweettext, outputFile)
}
