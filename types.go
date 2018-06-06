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
	url        string
	form       url.Values
	data       interface{}
	method     int
	responseCh chan response
}

// Data type for the bing json data
type Data struct {
	Images []struct {
		Startdate     string        `json:"startdate"`
		Fullstartdate string        `json:"fullstartdate"`
		Enddate       string        `json:"enddate"`
		URL           string        `json:"url"`
		Urlbase       string        `json:"urlbase"`
		Copyright     string        `json:"copyright"`
		Copyrightlink string        `json:"copyrightlink"`
		Quiz          string        `json:"quiz"`
		Wp            bool          `json:"wp"`
		Hsh           string        `json:"hsh"`
		Drk           int           `json:"drk"`
		Top           int           `json:"top"`
		Bot           int           `json:"bot"`
		Hs            []interface{} `json:"hs"`
	} `json:"images"`
	Tooltips struct {
		Loading  string `json:"loading"`
		Previous string `json:"previous"`
		Next     string `json:"next"`
		Walle    string `json:"walle"`
		Walls    string `json:"walls"`
	} `json:"tooltips"`
}
