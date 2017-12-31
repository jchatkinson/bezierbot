package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func searchForPhoto() (url string) {
	var data map[string]interface{}

	photoClient := http.Client{
		Timeout: time.Second * 2, // Maximum of 2 secs
	}

	req, err := http.NewRequest(http.MethodGet, BingURL, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("User-Agent", "Golang_Spider_Bot/0.1")

	res, getErr := photoClient.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}
	if err := json.Unmarshal(body, &data); err != nil {
		panic(err)
	}
	fmt.Println(data)

	return
}
