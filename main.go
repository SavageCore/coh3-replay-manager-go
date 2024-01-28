package main

import (
	"coh3-replay-manager-go/modules/game"
	"coh3-replay-manager-go/modules/replay"
	"coh3-replay-manager-go/modules/utils"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/Teages/go-autostart"
	"github.com/fynelabs/selfupdate"
	"github.com/gen2brain/beeep"
	"github.com/getlantern/systray"
)

const CurrentVersion = "v0.1.4"

func main() {
	if os.Getenv("DEV_MODE") != "true" {
		autoUpdate()

		// Run auto update every 24 hours
		go func() {
			for {
				time.Sleep(24 * time.Hour)
				autoUpdate()
			}
		}()
	}

	// Check if the app was opened with a command line argument
	if len(os.Args) > 1 && strings.HasPrefix(os.Args[1], "coh3-replay-manager-go://") {
		// Get the string containing the replay id and game version from the command line argument (without the protocol)
		input := strings.TrimPrefix(os.Args[1], "coh3-replay-manager-go://")

		parseUrlInput(input)

		// Exit the application
		return
	}

	systray.Run(onReady, onExit)
}

func onReady() {
	appPath, err := os.Executable()
	if err != nil {
		fmt.Println("Error getting app path:", err)
		return
	}

	app := &autostart.App{
		Name: "coh3-replay-manager-go",
		Exec: []string{appPath},
	}

	appIcon := getIconData("icon.ico")
	systray.SetIcon(appIcon)

	systray.SetTitle("Company of Heroes 3 Replay Manager")
	systray.SetTooltip("Company of Heroes 3 Replay Manager")

	mSetStartup := systray.AddMenuItem("Launch on startup", "Start this app when your computer starts")
	systray.AddSeparator()
	mAbout := systray.AddMenuItem(fmt.Sprintf("About (%s)", CurrentVersion), "")
	mQuit := systray.AddMenuItem("Exit", "")

	if app.IsEnabled() {
		mSetStartup.Check()
	} else {
		mSetStartup.Uncheck()
	}

	go func() {
		for range mSetStartup.ClickedCh {
			if app.IsEnabled() {
				mSetStartup.Uncheck()

				if err := app.Disable(); err != nil {
					fmt.Println(err)
				}
			} else {
				mSetStartup.Check()
				if err := app.Enable(); err != nil {
					fmt.Println(err)
				}
			}
		}
	}()

	go func() {
		for range mAbout.ClickedCh {
			utils.OpenUrl("https://github.com/SavageCore/coh3-replay-manager-go")
		}
	}()

	go func() {
		for range mQuit.ClickedCh {
			systray.Quit()
		}
	}()

	utils.RegisterUrlProtocol()
	go replay.InitialiseFolderWatcher()
}

func onExit() {
	os.Exit(0)
}

func parseUrlInput(input string) {
	re := regexp.MustCompile(`^([A-Za-z]+)/(\d+)/v/(\d.+)$`)

	if re.MatchString(input) {
		action := re.FindStringSubmatch(input)[1]
		// Get the ID from the URL
		id := re.FindStringSubmatch(input)[2]
		// Get the game version from the URL
		replayGameVersion := re.FindStringSubmatch(input)[3]
		gameVersion := game.GetGameVersion()

		// Check if the game versions match
		if gameVersion != replayGameVersion {
			title := "⛔ Replay cannot be played"
			message := fmt.Sprintf("⚠️ Game version: %s does not equal Replay's version: %s", gameVersion, replayGameVersion)
			err := beeep.Notify(title, message, "")
			if err != nil {
				fmt.Println(err)
			}

			return
		}

		// If the action is "play", play the replay
		if action == "play" {
			filename := replay.Download(id)
			replay.Play(filename)
		}

		if action == "download" {
			replay.Download(id)
		}

	} else {
		// Invalid URL
		fmt.Println("Invalid URL")
		os.Exit(1)
	}
}

func getIconData(file string) []byte {
	if file == "" {
		file = "icon.ico"
	}

	iconData, err := Asset("assets/icons/" + file)
	if err != nil {
		fmt.Println(err)
	}
	return iconData
}

func autoUpdate() bool {
	release, err := utils.GetLatestRelease()
	if err != nil {
		fmt.Println(err)
	}

	if release.TagName != CurrentVersion {
		// New version available
		title := "New version available"
		message := fmt.Sprintf("Version %s is available. Downloading and restarting the app.", release.TagName)
		err := beeep.Notify(title, message, "")
		if err != nil {
			fmt.Println(err)
		}

		for _, asset := range release.Assets {
			if asset.Name == "coh3-replay-manager-go_Windows_x86_64.zip" {
				// Set file path to tmp folder
				downloadPath := filepath.Join(os.TempDir(), asset.Name)

				utils.DownloadFile(asset.DownloadURL, downloadPath)
				fmt.Println("Downloaded file to", downloadPath)

				// Unzip the file
				utils.ExtractZip(downloadPath, os.TempDir())

				// Delete the zip file
				err := os.Remove(downloadPath)
				if err != nil {
					fmt.Println("Failed to delete file:", err)
					return false
				}

				extractedFilePath := filepath.Join(os.TempDir(), "coh3-replay-manager-go.exe")

				// Read the extracted file coh3-replay-manager-go.exe
				file, err := os.Open(extractedFilePath)
				if err != nil {
					fmt.Println("Failed to open file:", err)
					return false
				}
				defer file.Close()

				reader := io.Reader(file)

				// Update the app
				err = selfupdate.Apply(reader, selfupdate.Options{})
				if err != nil {
					// error handling
					fmt.Println("Failed to update:", err)
				}

				err = file.Close()
				if err != nil {
					fmt.Println("Failed to close file:", err)
					return false
				}

				// Get the path to the current executable
				exePath, err := os.Executable()
				if err != nil {
					fmt.Println("Failed to get executable path:", err)
					return false
				}

				// Delete the extracted file
				err = os.Remove(extractedFilePath)
				if err != nil {
					fmt.Println("Failed to delete file:", err)
					return false
				}

				// Restart the app
				cmd := exec.Command(exePath)
				_, err = cmd.Output()
				if err != nil {
					fmt.Println("Failed to restart:", err)
					return false
				}

				os.Exit(0)
			}
		}

		return true
	} else {
		return false
	}
}
