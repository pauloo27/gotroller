package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/user"
)

var cacheFolder string

func setupCacheFolder() {
	usr, err := user.Current()
	handleFatal(err)
	cacheFolder = fmt.Sprintf("%s/.cache/gotroller", usr.HomeDir)

	if _, err := os.Stat(cacheFolder); os.IsNotExist(err) {
		os.Mkdir(cacheFolder, os.ModePerm)
	}
}

func downloadAlbumArt(mediaUrl, artUrl string) (string, error) {
	outputFile := fmt.Sprintf("%s/.album-%s", cacheFolder, url.QueryEscape(mediaUrl))

	if _, err := os.Stat(outputFile); err == nil {
		fmt.Println("Cache found")
		return outputFile, nil
	}

	res, err := http.Get(artUrl)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	file, err := os.Create(outputFile)
	if err != nil {
		return "", err
	}

	defer file.Close()

	_, err = io.Copy(file, res.Body)
	if err != nil {
		return "", err
	}

	return outputFile, nil
}
