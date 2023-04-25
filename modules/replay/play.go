package replay

import (
	"coh3-replay-manager-go/modules/utils"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/mitchellh/go-ps"
)

// Launch the game
// Sample command:
// C:\Program Files (x86)\Steam\steam.exe -applaunch 1677280 -replay playback:file.rec
func Play(fileName string) {
	// Get Steam path
	steamPath := utils.GetSteamPath()
	launchOptions := fmt.Sprintf(`-applaunch 1677280 -replay playback:%s`, fileName)

	args := strings.Split(launchOptions, " ")

	user := utils.GetUsername()
	playbackDir := filepath.Join(user, "Documents", "My Games", "Company of Heroes 3", "playback")
	replayDir := filepath.Join(playbackDir, "replays")

	srcFile := filepath.Join(replayDir, fileName)
	destinationFile := filepath.Join(playbackDir, fileName)
	// Copy replay file from replayDir to playbackDir
	utils.CopyFile(srcFile, destinationFile)

	// Launch the game
	steamCmd := exec.Command(steamPath, args...)

	cmdErr := steamCmd.Start()

	if cmdErr != nil {
		fmt.Println(cmdErr.Error())
	}

	// Wait for Steam to exit
	cmdErr = steamCmd.Wait()
	if cmdErr != nil {
		fmt.Println("Steam exited with error:", cmdErr)
		return
	}

	processName := "RelicCoH3.exe"

	// Wait for the game to launch
	// Once the game has launched, we can delete the replay file
	for {
		// Get a list of all currently running processes
		processes, err := ps.Processes()
		if err != nil {
			fmt.Println("Failed to get processes:", err)
			return
		}

		// Loop through the list of processes and check for the target process name
		found := false
		for _, process := range processes {
			if process.Executable() == processName {
				found = true

				break
			}
		}

		if found {
			// Delete the replay file from the playback directory
			err = os.Remove(destinationFile)
			if err != nil {
				fmt.Println(err)
			}

			break
		}

		// Sleep for a while before checking again
		time.Sleep(time.Second * 30)
	}
}
