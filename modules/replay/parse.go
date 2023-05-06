package replay

import (
	"coh3-replay-manager-go/modules/shared"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/mholt/archiver/v4"
)

type ReplayObject struct {
	Version        int    `json:"Version"`
	Timestamp      string `json:"Timestamp"`
	MatchhistoryID int    `json:"MatchhistoryID"`
	Map            struct {
		Filename               string `json:"Filename"`
		LocalizedNameID        string `json:"LocalizedNameID"`
		LocalizedDescriptionID string `json:"LocalizedDescriptionID"`
	} `json:"Map"`
	Players []struct {
		Name      string `json:"Name"`
		Faction   string `json:"Faction"`
		Team      string `json:"Team"`
		SteamID   int64  `json:"SteamID"`
		ProfileID int    `json:"ProfileID"`
		Messages  []struct {
			Message string `json:"Message"`
			Tick    int    `json:"Tick"`
		} `json:"Messages"` // Change the type to the actual type if known
	} `json:"Players"`
	Length int `json:"Length"`
}

func formatTime(timeString string) string {
	// Best guess regexes to convert the time string to a format that can be parsed by Date
	// Examples:
	// 13.04.2023 20:08 DD.MM.YYYY HH:MM
	// 2023/4/9上午 12:45 YYYY/MM/DDAM HH:MM
	// 08/04/2023 19:28 DD/MM/YYYY HH:MM
	// 4/12/2023 6:47 PM MM/DD/YYYY HH:MM AM/PM
	// 10/3/2023 12:15 πμ
	// 2023-03-13 오후 8:52 YYYY-MM-DD PM H:MM
	// 2023-03-13 8:52 YYYY-MM-DD H:MM
	// 3/03/2023 5:39 pm MM/DD/YYYY H:MM AM/PM

	regex := regexp.MustCompile(`(\d{2})\.(\d{2})\.(\d{4}) (\d{2}):(\d{2})`)
	regex2 := regexp.MustCompile(`(\d{4})/(\d{1,2})/(\d{1,2})`)
	regex3 := regexp.MustCompile(`(\d{1,2})/(\d{1,2})/(\d{4}) (\d{2}):(\d{2})`)
	regex4 := regexp.MustCompile(`(\d{1,2})/(\d{1,2})/(\d{4}) (\d{1,2}):(\d{2}) (AM|PM|πμ)`)
	regex5 := regexp.MustCompile(`(\d{4})-(\d{1,2})-(\d{1,2}) (?:오전|오후) (\d{1,2}):(\d{1,2})`)
	regex6 := regexp.MustCompile(`(\d{4})-(\d{1,2})-(\d{1,2}) (\d{1,2}):(\d{1,2})`)
	regex7 := regexp.MustCompile(`(\d{1,2})/(\d{1,2})/(\d{4}) (\d{1,2}):(\d{2}) (am|pm)`)

	guess1 := regex.FindStringSubmatch(timeString)
	guess2 := regex2.FindStringSubmatch(timeString)
	guess3 := regex3.FindStringSubmatch(timeString)
	guess4 := regex4.FindStringSubmatch(timeString)
	guess5 := regex5.FindStringSubmatch(timeString)
	guess6 := regex6.FindStringSubmatch(timeString)
	guess7 := regex7.FindStringSubmatch(timeString)

	originalTimeString := timeString

	if len(guess1) > 0 {
		timeString = fmt.Sprintf("%s-%s-%sT%s:%s:00Z", pad(guess1[3]), pad(guess1[2]), pad(guess1[1]), pad(guess1[4]), pad(guess1[5]))
	} else if len(guess2) > 0 {
		timeString = fmt.Sprintf("%s-%s-%sT00:00:00Z", pad(guess2[1]), pad(guess2[2]), pad(guess2[3]))
	} else if len(guess3) > 0 {
		if num, err := strconv.Atoi(guess3[1]); err == nil && num > 12 {
			timeString = fmt.Sprintf("%s-%s-%sT%s:%s:00Z", pad(guess3[3]), pad(guess3[2]), pad(guess3[1]), pad(guess3[4]), pad(guess3[5]))
		} else {
			timeString = fmt.Sprintf("%s-%s-%sT%s:%s:00Z", pad(guess3[3]), pad(guess3[1]), pad(guess3[2]), pad(guess3[4]), pad(guess3[5]))
		}
	} else if len(guess4) > 0 {
		if num, err := strconv.Atoi(guess4[1]); err == nil && num > 12 {
			timeString = fmt.Sprintf("%s-%s-%sT%s:%s:00Z", pad(guess4[3]), pad(guess4[2]), pad(guess4[1]), pad(guess4[4]), pad(guess4[5]))
		} else {
			timeString = fmt.Sprintf("%s-%s-%sT%s:%s:00Z", pad(guess4[3]), pad(guess4[1]), pad(guess4[2]), pad(guess4[4]), pad(guess4[5]))
		}
	} else if len(guess5) > 0 {
		timeString = fmt.Sprintf("%s-%s-%sT%s:%s:00Z", pad(guess5[1]), pad(guess5[2]), pad(guess5[3]), pad(guess5[4]), pad(guess5[5]))
	} else if len(guess6) > 0 {
		timeString = fmt.Sprintf("%s-%s-%sT%s:%s:00Z", pad(guess6[1]), pad(guess6[2]), pad(guess6[3]), pad(guess6[4]), pad(guess6[5]))
	} else if len(guess7) > 0 {
		if num, err := strconv.Atoi(guess7[1]); err == nil && num > 12 {
			timeString = fmt.Sprintf("%s-%s-%sT%s:%s:00Z", pad(guess7[3]), pad(guess7[2]), pad(guess7[1]), pad(guess7[4]), pad(guess7[5]))
		} else {
			timeString = fmt.Sprintf("%s-%s-%sT%s:%s:00Z", pad(guess7[3]), pad(guess7[1]), pad(guess7[2]), pad(guess7[4]), pad(guess7[5]))
		}
	}

	t, err := time.Parse(time.RFC3339, timeString)
	if err != nil {
		fmt.Println("Error parsing time string:", originalTimeString)
		return originalTimeString
	}

	return t.Format("2006-01-02 15:04")
}

func pad(s string) string {
	if len(s) == 1 {
		return "0" + s
	}
	return s
}

func Parse(filename string) (ReplayObject, error) {
	replayDir := shared.GetReplayDir()
	replayFilePath := filepath.Join(replayDir, filename)

	_, err := os.Stat(replayFilePath)
	if err != nil {
		fmt.Println("Error:", err)
		return ReplayObject{}, err
	}

	// Download the latest version of flank if it doesn't exist
	_, err = os.Stat("flank.exe")
	if err != nil {
		fmt.Println("flank binary not found, downloading...")
		fp := <-downloadFlank()
		if fp == "error" {
			return ReplayObject{}, err
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

	// fmt.Println("Parsing replay file:", replayFilePath)

	// Parse the replay file with flank and return a Replay object
	// flank returns a JSON object, so you can use the json package to parse it
	// Example: flank.exe /path/to/replay.rec

	cmd := exec.Command("./flank.exe", replayFilePath)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow:    true,
		CreationFlags: 0x08000000,
	}
	out, err := cmd.Output()
	if err != nil {
		fmt.Println("Error:", err)
	}

	var replay ReplayObject

	err = json.Unmarshal(out, &replay)
	if err != nil {
		fmt.Println("Error:", err)
		return ReplayObject{}, err
	}

	formattedTimestamp := formatTime(replay.Timestamp)
	replay.Timestamp = formattedTimestamp

	// Example filename: "data:scenarios\multiplayer\twin_beach_2p_mkii\twin_beach_2p_mkii"
	// Another: "data:scenarios\multiplayer\(2) crossroads\(2) crossroads"
	// We want to return the last part of the path
	mapFilename := strings.Split(replay.Map.Filename, "\\")

	replay.Map.Filename = mapFilename[len(mapFilename)-1]

	return replay, nil
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
