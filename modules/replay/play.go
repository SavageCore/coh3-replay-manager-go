package replay

import (
	"coh3-replay-manager-go/modules/utils"
	"fmt"
	"os/exec"
	"strings"

	"github.com/gen2brain/beeep"
)

// Launch the game
// Sample command:
// C:\Program Files (x86)\Steam\steam.exe -applaunch 1677280 -replay playback:file.rec
func Play(fileName string) {
	// Get Steam path
	steamPath := utils.GetSteamPath()
	launchOptions := fmt.Sprintf(`-applaunch 1677280 -replay playback:%s`, fileName)

	args := strings.Split(launchOptions, " ")

	cmd := exec.Command(steamPath, args...)

	_, cmdErr := cmd.Output()

	if cmdErr != nil {
		fmt.Println(cmdErr.Error())
	}

	title := "Replay launched 🚀"
	message := "⚠️ Look out for a confirmation window from Steam, allowing you to launch the game. ⚠️"
	err := beeep.Notify(title, message, "")
	if err != nil {
		panic(err)
	}
}
