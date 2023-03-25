package replay

import (
	"fmt"

	"github.com/gen2brain/beeep"
	"github.com/skratchdot/open-golang/open"
)

func Play(fileName string) {
	url := fmt.Sprintf("steam://run/1677280//-replay playback:%s", fileName)

	open.Run(url)

	title := "Replay launched 🚀"
	message := "⚠️ Look out for a confirmation window from Steam, allowing you to launch the game. ⚠️"
	err := beeep.Notify(title, message, "")
	if err != nil {
		panic(err)
	}
}
