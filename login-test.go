package main

import (
	"fmt"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

// this gets the username from twitter. To use, `go run login-test.go`
func testlogin() {
	consumerKey := "Qmh1KE0vfbKH22YrJpyRkGuze"
	consumerSecret := "xaKQpg02yHeWMFEcb5CgEaxIUSi69p4gYTEh4kYXCC53Egurps"
	accessToken := "943990759975354369-Bs8xe2kqOkRDW6HasNtK7Di2ZXoHDqq"
	accessSecret := "RfyWA1VlPV2eudfUOhoetUVv3FZXjVF1KqMKImgNi16oe"

	config := oauth1.NewConfig(consumerKey, consumerSecret)
	token := oauth1.NewToken(accessToken, accessSecret)

	// OAuth1 http.Client will automatically authorize Requests
	httpClient := config.Client(oauth1.NoContext, token)

	// Twitter client
	client := twitter.NewClient(httpClient)

	// Verify Credentials
	verifyParams := &twitter.AccountVerifyParams{
		IncludeEmail: twitter.Bool(true),
	}
	user, _, _ := client.Accounts.VerifyCredentials(verifyParams)
	fmt.Printf("User's Name:%+v\n", user.ScreenName)
}
