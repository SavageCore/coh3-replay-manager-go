package replay

import (
	"coh3-replay-manager-go/modules/utils"
	"fmt"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"regexp"

	"go.etcd.io/bbolt"
)

func Download(id string) {
	user := utils.GetUsername()
	fileName := fmt.Sprintf("downloaded-replay-%s.rec", id)
	replayDir := filepath.Join(user, "Documents", "My Games", "Company of Heroes 3", "playback", "replays")

	// If another replay with the same id exists, play it instead of downloading it again
	files, err := os.ReadDir(replayDir)
	if err != nil {
		panic(err)
	}

	re := regexp.MustCompile(fmt.Sprintf(`^downloaded-replay-%s.rec$`, id))

	for _, file := range files {
		if re.MatchString(file.Name()) {
			fmt.Printf("Found replay file with ID %s: %s\n", id, file.Name())
			Play(file.Name())
			return
		}
	}

	// Proceed with the download
	filePath := filepath.Join(replayDir, fileName)
	url := fmt.Sprintf("https://cohdb.com/replays/%s/download", id)

	// Create the file
	file, err := os.Create(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Download the file
	response, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	// Write the response body to the file
	_, err = file.ReadFrom(response.Body)
	if err != nil {
		panic(err)
	}

	// Get the original file name from the Content-Disposition header
	contentDisposition := response.Header.Get("Content-Disposition")

	_, params, err := mime.ParseMediaType(contentDisposition)
	if err != nil {
		panic(err)
	}

	originalFileName := params["filename"]

	// Create a simple ley/value database to store the original file names against the replay IDs within replay directory
	db, err := bbolt.Open(filepath.Join(replayDir, "coh3rm.db"), 0600, nil)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	defer db.Close()

	err = db.Update(func(tx *bbolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte("replayNames"))
		if err != nil {
			fmt.Println(err)
			return err
		}

		err = bucket.Put([]byte(id), []byte(originalFileName))
		if err != nil {
			fmt.Println(err)
			return err
		}

		return nil
	})

	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	// Play the replay
	Play(fileName)
}
