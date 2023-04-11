package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Release represents a GitHub release
type Release struct {
	TagName string `json:"tag_name"`
	Name    string `json:"name"`
	Assets  []struct {
		Name        string `json:"name"`
		DownloadURL string `json:"browser_download_url"`
	} `json:"assets"`
}

func GetLatestRelease() (Release, error) {
	var release Release
	repoOwner := "SavageCore"
	repoName := "coh3-replay-manager-go"

	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", repoOwner, repoName)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error:", err)
		return release, err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&release)
	if err != nil {
		fmt.Println("Error:", err)
		return release, err
	}

	return release, nil
}
