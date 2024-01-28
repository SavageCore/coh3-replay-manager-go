package replay

import (
	"coh3-replay-manager-go/modules/utils"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/fsnotify/fsnotify"
)

func InitialiseFolderWatcher() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		panic(err)
	}
	defer watcher.Close()

	user := utils.GetUsername()

	dir := filepath.Join(user, "Documents", "My Games", "Company of Heroes 3", "playback")
	err = watcher.Add(dir)
	if err != nil {
		panic(err)
	}

	// Ensure replays folder exists under dir
	replaysDir := dir
	if _, err := os.Stat(replaysDir); os.IsNotExist(err) {
		os.Mkdir(replaysDir, os.ModePerm)
	}

	// Define a map to store the file modification times
	modTimes := make(map[string]time.Time)

	// Define a function to rename the file with a unique name
	renameFile := func(path string) error {
		ext := filepath.Ext(path)
		filename := filepath.Base(path)
		prefix := ""

		if filename == "temp_campaign.rec" {
			prefix = "campaign-"
		}

		newPath := filepath.Join(dir, fmt.Sprintf("saved-%sreplay-%d%s", prefix, time.Now().UnixNano(), ext))
		return os.Rename(path, newPath)
	}

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			// Check if the event is a file modification
			if event.Op&fsnotify.Write == fsnotify.Write {
				// Check if the file is the one you're interested in (temp.rec or temp_campaign.rec)
				if filepath.Base(event.Name) == "temp.rec" || filepath.Base(event.Name) == "temp_campaign.rec" {
					// Get the file modification time
					fi, err := os.Stat(event.Name)
					if err != nil {
						fmt.Println("Error stating file:", err)
					}
					modTime := fi.ModTime()
					fileSize := fi.Size()

					// Check if the file is empty
					if fileSize == 0 {
						fmt.Println("File is empty, skipping...")
						continue
					}

					// Check if the file has been modified since last time
					lastModTime, ok := modTimes[event.Name]
					if ok && modTime.Equal(lastModTime) {
						continue
					}
					modTimes[event.Name] = modTime

					// Rename the file with a unique name
					err = renameFile(event.Name)
					if err != nil {
						fmt.Println("Error renaming file:", err)
					}
				}
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			panic(err)
		}
	}
}
