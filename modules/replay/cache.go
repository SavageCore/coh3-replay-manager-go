package replay

import (
	"coh3-replay-manager-go/modules/shared"
	"encoding/json"
	"fmt"
	"path/filepath"

	"go.etcd.io/bbolt"
)

// Function to cache replay objects in bbolt database, saved against the filename
func Cache(filename string, replay ReplayObject) {
	db, err := bbolt.Open(filepath.Join(shared.GetReplayDir(), "coh3rm.db"), 0600, nil)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	defer db.Close()

	err = db.Update(func(tx *bbolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte("replayObjects"))
		if err != nil {
			fmt.Println(err)
			return err
		}

		// Convert the replay object to a JSON string
		replayJson, _ := json.Marshal(replay)

		// Store the replay object against the replay filename
		err = bucket.Put([]byte(filename), []byte(replayJson))
		if err != nil {
			fmt.Println(err)
			return err
		}

		return nil
	})

	if err != nil {
		fmt.Println(err)
	}
}

// Function to get a replay object from the bbolt database, using the filename
func GetCached(filename string) (ReplayObject, error) {
	var replay ReplayObject

	db, err := bbolt.Open(filepath.Join(shared.GetReplayDir(), "coh3rm.db"), 0600, nil)
	if err != nil {
		fmt.Println(err)
		return ReplayObject{}, err
	}
	defer db.Close()

	err = db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte("replayObjects"))
		if bucket == nil {
			return nil
		}

		// Get the replay object from the database
		replayJson := bucket.Get([]byte(filename))

		// Convert the JSON string to a ReplayObject
		json.Unmarshal(replayJson, &replay)

		return nil
	})

	if err != nil {
		fmt.Println(err)
		return ReplayObject{}, err
	}

	return replay, nil
}
