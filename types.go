package main

import "net/url"

// Media type for uploading photos
type Media struct {
	MediaID       int64  `json:"media_id"`
	MediaIDString string `json:"media_id_string"`
	Size          int    `json:"size"`
	Image         Image  `json:"image"`
}

// Image type for uploading photos
type Image struct {
	W         int    `json:"w"`
	H         int    `json:"h"`
	ImageType string `json:"image_type"`
}

type response struct {
	data interface{}
	err  error
}

// Query type for a query
type Query struct {
	url         string
	form        url.Values
	data        interface{}
	method      int
	response_ch chan response
}
