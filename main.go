package main

import (
	"coh3-replay-manager-go/modules/game"
	"coh3-replay-manager-go/modules/replay"
	"coh3-replay-manager-go/modules/utils"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/Teages/go-autostart"
	"github.com/gen2brain/beeep"
	"github.com/getlantern/systray"
)

func main() {
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
	mAbout := systray.AddMenuItem("About", "")
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
	re := regexp.MustCompile(`^play/(\d+)/v/(\d+)$`)

	if re.MatchString(input) {
		// Get the ID from the URL
		id := re.FindStringSubmatch(input)[1]
		// Get the game version from the URL
		replayGameVersion := re.FindStringSubmatch(input)[2]
		gameVersion := game.GetGameVersion()

		// Check if the game versions match
		if gameVersion != replayGameVersion {
			title := "⛔ Replay cannot be played"
			message := fmt.Sprintf("⚠️ Game version: %s does not equal Replay's version: %s", gameVersion, replayGameVersion)
			err := beeep.Notify(title, message, "")
			if err != nil {
				panic(err)
			}

			return
		}

		// Download the replay
		replay.Download(id)
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
		panic(err)
	}
	return iconData
}
