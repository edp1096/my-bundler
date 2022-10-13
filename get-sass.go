package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"builder/util/arc"
	"builder/util/dl"
)

type ReleaseStatus struct {
	Version string `json:"tag_name"`
	Assets  []struct {
		Name        string `json:"name"`
		DownloadURL string `json:"browser_download_url"`
	}
}

var dartsassRepos = "sass/dart-sass-embedded"

func checkSassExists() error {
	var err error

	dsDir := "sass_embedded"
	checkPaths := []string{
		dsDir,
		dsDir + "/dart-sass-embedded.bat",
		dsDir + "/src/dart-sass-embedded.snapshot",
		dsDir + "/src/dart.exe",
	}

	for _, path := range checkPaths {
		if _, err = os.Stat(path); os.IsNotExist(err) {
			return err
		}
	}

	return nil
}

func getSass() {
	var err error

	releaseUri := "https://api.github.com/repos/" + dartsassRepos + "/releases/latest"
	r, err := http.Get(releaseUri)
	if err != nil {
		log.Fatalln(err)
	}
	defer r.Body.Close()

	var release ReleaseStatus
	err = json.NewDecoder(r.Body).Decode(&release)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(release.Version)

	downloadURL := ""
	for _, asset := range release.Assets {
		if strings.Contains(asset.DownloadURL, "windows-x64") {
			downloadURL = asset.DownloadURL
			break
		}
	}

	// Dart sass embedded - https://github.com/sass/dart-sass-embedded/releases
	fmt.Println("Download, Dart-sass-embedded")
	saveName := "sass_embedded.zip"
	err = dl.DownloadFile("sass_embedded.zip", downloadURL)
	if err != nil {
		log.Fatalln(err)
	}

	// destination := strings.TrimSuffix(saveName, filepath.Ext(saveName))
	err = arc.DecompressZIP(saveName, "")
	if err != nil {
		log.Fatalln(err)
	}
}
