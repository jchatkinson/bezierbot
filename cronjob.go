package main

import mathrand "math/rand"

func getPhotoOfTheDay() error {
	var err error
	photourl := searchForPhoto()
	uuid, inputFile := downloadPhoto(photourl)
	outputFile := "./img/output/" + uuid + "-out.jpg"
	n := 200 + mathrand.Intn(300)
	modes := []int{1, 4, 6}
	m := 0
	processPhoto(inputFile, outputFile, n, modes[m])
	return err
}
