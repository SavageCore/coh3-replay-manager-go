package replay

import (
	"coh3-replay-manager-go/modules/utils"
	"fmt"
	"os"
	"path/filepath"

	"go.etcd.io/bbolt"
)

// Function to list all replays in the database
// returns a map of replay IDs and their original filenames
func ListDownloaded() map[string]string {
	replays := make(map[string]string)
	user := utils.GetUsername()
	replayDir := filepath.Join(user, "Documents", "My Games", "Company of Heroes 3", "playback")

	// If the database doesn't exist, return an empty map
	if _, err := os.Stat(filepath.Join(replayDir, "coh3rm.db")); os.IsNotExist(err) {
		return replays
	}

	db, err := os.Open(filepath.Join(replayDir, "coh3rm.db"))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Open the database
	boltDB, err := bbolt.Open(db.Name(), 0600, nil)
	if err != nil {
		panic(err)
	}
	defer boltDB.Close()

	// List all replays in the database
	boltDB.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("replayNames"))
		if b == nil {
			fmt.Println("No replays found")
			return nil
		}

		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			replays[string(k)] = string(v)
		}

		// Return nil to signal that we've finished iterating over the bucket
		return nil
	})

	// Return the map of replays
	return replays
}

func ListLocal() map[string]string {
	replays := make(map[string]string)
	user := utils.GetUsername()
	replayDir := filepath.Join(user, "Documents", "My Games", "Company of Heroes 3", "playback")

	// Replay files have the following formats:
	// saved-replay-{unixnano}.rec
	// saved-campaign-replay-{unixnano}.rec

	// Get all files in the replay directory
	files, err := os.ReadDir(replayDir)
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return replays
	}

	var id = 0

	// Iterate over all files in the directory
	for _, file := range files {
		// If the file is a replay file, add it to the map
		if matchesPattern(file.Name(), "saved-replay-*.rec") || matchesPattern(file.Name(), "saved-campaign-replay-*.rec") {
			replays[fmt.Sprintf("%d", id)] = file.Name()
			id++
		}
	}

	// Return the map of replays
	return replays
}

// Helper function to check if a file name matches a pattern
func matchesPattern(fileName, pattern string) bool {
	matched, err := filepath.Match(pattern, fileName)
	if err != nil {
		fmt.Println("Error matching pattern:", err)
		return false
	}
	return matched
}
