package shared

import (
	"path/filepath"
)

func GetReplayDir() string {
	user := GetUsername()
	replayDir := filepath.Join(user, "Documents", "My Games", "Company of Heroes 3", "playback", "replays")

	return replayDir
}

func GetPlaybackDir() string {
	user := GetUsername()
	playbackDir := filepath.Join(user, "Documents", "My Games", "Company of Heroes 3", "playback")

	return playbackDir
}
