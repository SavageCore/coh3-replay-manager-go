package replay

import (
	"coh3-replay-manager-go/modules/utils"
	"fmt"
	"os"
	"path/filepath"

	"go.etcd.io/bbolt"
)

func Remove(fileName string) {
	user := utils.GetUsername()
	replayDir := filepath.Join(user, "Documents", "My Games", "Company of Heroes 3", "playback")
	filePath := filepath.Join(replayDir, fileName)

	// If fileName includes downloaded-replay- prefix, remove it from the database
	if fileName[:18] == "downloaded-replay-" {
		db, err := bbolt.Open(filepath.Join(replayDir, "coh3rm.db"), 0600, nil)
		if err != nil {
			fmt.Println(err)
			panic(err)
		}
		defer db.Close()

		// Extract the replay ID from the file name
		replayID := fileName[18 : len(fileName)-4]

		fmt.Println("Removing: " + replayID + " from the database...")

		err = db.Update(func(tx *bbolt.Tx) error {
			b := tx.Bucket([]byte("replayNames"))
			err := b.Delete([]byte(replayID))
			return err
		})

		if err != nil {
			fmt.Println(err)
			panic(err)
		}
	}

	// Check if file exists then remove it
	_, err := os.Stat(filePath)
	if err == nil {
		err := os.Remove(filePath)
		if err != nil {
			panic(err)
		}
	}
}
