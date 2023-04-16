package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/mholt/archiver/v4"
)

type ReplayObject struct {
	Version        int    `json:"version"`
	Timestamp      string `json:"timestamp"`
	MatchhistoryID int    `json:"matchhistory_id"`
	Map            struct {
		Filename               string `json:"filename"`
		LocalizedNameID        string `json:"localized_name_id"`
		LocalizedDescriptionID string `json:"localized_description_id"`
	} `json:"map"`
	Players []struct {
		Name      string `json:"name"`
		Faction   string `json:"faction"`
		Team      string `json:"team"`
		SteamID   int64  `json:"steam_id"`
		ProfileID int    `json:"profile_id"`
		Messages  []struct {
			Message string `json:"message"`
			Tick    int    `json:"tick"`
		} `json:"messages"` // Change the type to the actual type if known
	} `json:"players"`
	Length int `json:"length"`
}

func ParseReplay(filename string) {
	user := GetUsername()
	replayDir := filepath.Join(user, "Documents", "My Games", "Company of Heroes 3", "playback", "replays")
	replayFilePath := filepath.Join(replayDir, filename)

	_, err := os.Stat(replayFilePath)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Checking for flank binary...")
	// Download the latest version of flank if it doesn't exist
	_, err = os.Stat("flank.exe")
	if err != nil {
		fmt.Println("flank binary not found, downloading...")
		fp := <-downloadFlank()
		if fp == "error" {
			return
		}

		fmt.Println("Extracting flank binary...")

		fsys, err := archiver.FileSystem(fp)
		if err != nil {
			fmt.Println("Error:", err)
		}

		fs.WalkDir(fsys, ".", func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				fmt.Println("Error:", err)
			}
			if !d.IsDir() {
				if filepath.Base(path) == "flank.exe" {
					src, err := fsys.Open(path)
					if err != nil {
						fmt.Println("Error:", err)
					}
					defer src.Close()

					dst, err := os.Create(filepath.Base(path))
					if err != nil {
						fmt.Println("Error:", err)
					}
					defer dst.Close()

					_, err = io.Copy(dst, src)
					if err != nil {
						fmt.Println("Error:", err)
					}
				}
			}

			return nil
		})

		// Delete the downloaded file
		err = os.Remove(filepath.Join(os.TempDir(), "flank.tar.gz"))
		if err != nil {
			fmt.Println("Failed to delete file:", err)
		}

		if err != nil {
			fmt.Println("Error:", err)
		}
	}

	fmt.Println("Parsing replay file:", replayFilePath)

	// Parse the replay file with flank and return a Replay object
	// flank returns a JSON object, so you can use the json package to parse it
	// Example: flank.exe /path/to/replay.rec

	cmd := exec.Command("./flank.exe", replayFilePath)
	out, err := cmd.Output()
	if err != nil {
		fmt.Println("Error:", err)
	}

	var obj ReplayObject

	err = json.Unmarshal(out, &obj)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// For each player in the replay, print their name and steam ID
	for _, player := range obj.Players {
		fmt.Println(player.Name, player.Team, player.SteamID)
	}
}

func downloadFlank() <-chan string {
	c := make(chan string)
	go func() {
		fmt.Println("Downloading flank binary...")
		// Download the latest version of flank from https://github.com/ryantaylor/flank/releases
		// And extract it to the current directory
		// Windows: x86_64-pc-windows-gnu.tar.gz
		// Linux: x86_64-unknown-linux-gnu.tar.gz

		// Filename
		var filename string
		switch runtime.GOOS {
		case "windows":
			filename = "x86_64-pc-windows-gnu.tar.gz"
		case "linux":
			filename = "x86_64-unknown-linux-gnu.tar.gz"
		default:
			fmt.Println("Unsupported OS:", runtime.GOOS)
			c <- "error"
		}

		// Download the latest version of flank from GitHub releases
		url := fmt.Sprintf("https://github.com/ryantaylor/flank/releases/latest/download/%s", filename)
		resp, err := http.Get(url)
		if err != nil {
			fmt.Println("Failed to download flank:", err)
			c <- "error"
		}
		defer resp.Body.Close()

		fmt.Println(os.TempDir())

		// Create a new file to save the downloaded flank archive in tmp folder
		filepath := filepath.Join(os.TempDir(), "flank.tar.gz")
		file, err := os.Create(filepath)
		if err != nil {
			fmt.Println("Failed to create file:", err)
			c <- "error"
		}
		defer file.Close()

		// Copy the downloaded flank binary to the newly created file
		_, err = io.Copy(file, resp.Body)
		if err != nil {
			fmt.Println("Failed to copy file:", err)
			c <- "error"
		}

		fmt.Println("flank binary downloaded")
		c <- filepath
	}()
	return c
}
