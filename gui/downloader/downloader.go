package downloader

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

const ROOT_DOWNLOAD_FOLDER = "/tmp/gotroller-albums"

var rootFolderExist = false

func DownloadRemoteArt(artURL string) (string, error) {
	if !rootFolderExist {
		if _, err := os.Stat(ROOT_DOWNLOAD_FOLDER); os.IsNotExist(err) {
			err := os.Mkdir(ROOT_DOWNLOAD_FOLDER, os.ModePerm)
			if err != nil {
				return "", err
			}
		}
		rootFolderExist = true
	}

	outputPath := fmt.Sprintf("%s/%s", ROOT_DOWNLOAD_FOLDER, url.QueryEscape(artURL))
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		res, err := http.Get(artURL)
		if err != nil {
			return "", nil
		}
		defer res.Body.Close()

		file, err := os.Create(outputPath)
		if err != nil {
			return "", nil
		}
		defer file.Close()

		_, err = io.Copy(file, res.Body)
		if err != nil {
			return "", nil
		}

		return outputPath, nil
	}

	return outputPath, nil
}
