package main

import (
	"coh3-replay-manager-go/modules/game"
	"coh3-replay-manager-go/modules/replay"
	"coh3-replay-manager-go/modules/utils"
	"context"
	"embed"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"syscall"
	"time"

	"github.com/Masterminds/semver/v3"
	"github.com/Teages/go-autostart"
	"github.com/fynelabs/selfupdate"
	"github.com/gen2brain/beeep"
	"github.com/getlantern/systray"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

const CurrentVersion = "v0.1.5-pre"

//go:embed all:frontend/dist
var assets embed.FS

// Create an instance of the app structure
var app = NewApp()
var replayWindowOpen = false
var mReplayListView *systray.MenuItem

// var wv webview.WebView

func main() {
	// Check if the app was opened with a command line argument
	if len(os.Args) > 1 && strings.HasPrefix(os.Args[1], "coh3-replay-manager-go://") {
		// Print
		fmt.Println("App opened with command line argument:", os.Args[1])
		// Get the string containing the replay id and game version from the command line argument (without the protocol)
		input := strings.TrimPrefix(os.Args[1], "coh3-replay-manager-go://")

		parseUrlInput(input)

		// Exit the application
		return
	}

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

	systray.Register(onReady, onExit)

	StartHidden := true

	// If in development mode, StartHidden is set to false so window is visible
	if os.Getenv("DEV_MODE") == "true" {
		StartHidden = false
	}

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "Company of Heroes 3 Replay Manager (" + CurrentVersion + ")",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		OnShutdown:       app.shutdown,
		OnBeforeClose:    beforeClose,
		Bind: []interface{}{
			app,
		},
		StartHidden:       StartHidden,
		HideWindowOnClose: false,
	})

	if err != nil {
		println("Error:", err.Error())
	}
}

func onReady() {
	appPath, err := os.Executable()
	if err != nil {
		fmt.Println("Error getting app path:", err)
		return
	}

	autostartApp := &autostart.App{
		Name: "coh3-replay-manager-go",
		Exec: []string{appPath},
	}

	appIcon := getIconData("icon.ico")
	systray.SetIcon(appIcon)

	systray.SetTitle("Company of Heroes 3 Replay Manager")
	systray.SetTooltip("Company of Heroes 3 Replay Manager")

	mSetStartup := systray.AddMenuItem("Launch on startup", "Start this app when your computer starts")
	mReplayListView = systray.AddMenuItem("View replays", "View saved and downloaded replays")
	systray.AddSeparator()
	mAbout := systray.AddMenuItem(fmt.Sprintf("About (%s)", CurrentVersion), "")
	mQuit := systray.AddMenuItem("Exit", "")

	if autostartApp.IsEnabled() {
		mSetStartup.Check()
	} else {
		mSetStartup.Uncheck()
	}

	go func() {
		for range mSetStartup.ClickedCh {
			if autostartApp.IsEnabled() {
				mSetStartup.Uncheck()

				if err := autostartApp.Disable(); err != nil {
					fmt.Println(err)
				}
			} else {
				mSetStartup.Check()
				if err := autostartApp.Enable(); err != nil {
					fmt.Println(err)
				}
			}
		}
	}()

	go func() {
		for range mReplayListView.ClickedCh {
			if replayWindowOpen {
				runtime.Hide(app.ctx)
				replayWindowOpen = false
				mReplayListView.SetTitle("View replays")
			} else {
				runtime.Show(app.ctx)
				replayWindowOpen = true
				mReplayListView.SetTitle("Hide replays")

				replay.List()
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
	re := regexp.MustCompile(`^([A-Za-z]+)/(\d+)/v/(\d+)$`)

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

			// return
		}

		// If the action is "play", play the replay
		if action == "play" {
			filename := replay.Download(id, replayGameVersion)
			replay.Play(filename)
		}

		if action == "download" {
			replay.Download(id, replayGameVersion)
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

	currentVersionSem, err := semver.NewVersion(CurrentVersion)
	if err != nil {
		fmt.Println(err)
	}

	latestVersionSem, err := semver.NewVersion(release.TagName)
	if err != nil {
		fmt.Println(err)
	}

	if currentVersionSem.LessThan(latestVersionSem) {
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
				cmd.SysProcAttr = &syscall.SysProcAttr{
					HideWindow:    true,
					CreationFlags: 0x08000000,
				}
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

func beforeClose(ctx context.Context) bool {
	// Hide the window
	runtime.Hide(ctx)

	// Update the systray menu item
	replayWindowOpen = false
	mReplayListView.SetTitle("View replays")

	// Return true to prevent the app from closing
	return true
}
