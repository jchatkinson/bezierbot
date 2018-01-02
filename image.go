package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"

	"github.com/tidwall/gjson"
)

func searchForPhoto() (url string) {
	var data map[string]interface{}

	photoClient := http.Client{
		Timeout: time.Second * 2, // Maximum of 10 secs
	}

	req, err := http.NewRequest(http.MethodGet, BingURL, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2228.0 Safari/537.36")
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
	json := string(body)
	url = gjson.Get(json, "images.0.url").String()
	url = "https://bing.com" + url
	// fmt.Println(json)
	fmt.Println(url)

	return
}

func downloadPhoto(url string) (uuid, imgpath string) {
	response, e := http.Get(url)
	if e != nil {
		log.Fatalf("http.Get -> %v", e)
	}
	defer response.Body.Close()

	//open file for writing
	uuid, err := newUUID()
	filename := uuid + ".jpg"
	imgpath = filepath.Join(InputFolder, filename)

	// Use io.Copy to just dump the response body to the file. This supports huge files
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatalf("ioutil.Readall -> %v", err)
	}
	ioutil.WriteFile(imgpath, data, 0666)

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
