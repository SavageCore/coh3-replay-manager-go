package replay

import (
	"coh3-replay-manager-go/modules/shared"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func InitialiseFolderWatcher() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		panic(err)
	}
	defer watcher.Close()

	user := shared.GetUsername()

	dir := filepath.Join(user, "Documents", "My Games", "Company of Heroes 3", "playback")
	err = watcher.Add(dir)
	if err != nil {
		panic(err)
	}

	// Ensure replays folder exists under dir
	replaysDir := filepath.Join(dir, "replays")
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

		newPath := filepath.Join(dir, "replays", fmt.Sprintf("saved-%sreplay-%d%s", prefix, time.Now().UnixNano(), ext))
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

					// EventsEmit(ctx, "newReplay", nil)
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

func InitialiseDownloadWatcher(ctx context.Context) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		panic(err)
	}
	defer watcher.Close()

	user := shared.GetUsername()

	dir := filepath.Join(user, "Documents", "My Games", "Company of Heroes 3", "playback", "replays")
	err = watcher.Add(dir)
	if err != nil {
		panic(err)
	}

	// Define a map to store the file modification times
	modTimes := make(map[string]time.Time)

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}

			// Check if the event is a file creation or write
			if event.Op&fsnotify.Create == fsnotify.Create || event.Op&fsnotify.Write == fsnotify.Write {
				// Check if the file extension is .rec and starts with downloaded-
				if filepath.Ext(event.Name) == ".rec" && filepath.Base(event.Name)[:11] == "downloaded-" {
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

					filename := filepath.Base(event.Name)

					replay, err := Parse(filename)
					if err != nil {
						fmt.Println("Error parsing replay:", err)
					}

					// Broadcast the replay to the frontend
					runtime.EventsEmit(ctx, "replay:downloaded", replay)
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
