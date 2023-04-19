package replay

import (
	"coh3-replay-manager-go/modules/shared"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Replay struct
type Replay struct {
	Filename  string
	Version   int
	Timestamp string
	Length    int
	Map       Map
	Players   []Player
}

type Map struct {
	Filename               string
	LocalizedNameID        string
	LocalizedDescriptionID string
}

type Player struct {
	Name      string
	Faction   string
	Team      string
	SteamID   string
	ProfileID int
	Messages  []Message
}

type Message struct {
	Message string `json:"message"`
	Tick    int    `json:"tick"`
}

// Function to list all replays by ".rec" files in the replay directory
func List() []Replay {
	var replays []Replay
	replayDir := shared.GetReplayDir()

	files, err := os.ReadDir(replayDir)
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".rec" {
			// Parse and add the replay to the list
			// If file.Name() includes "saved" then it's a local replay, do not parse
			if strings.Contains(file.Name(), "saved") {
				continue
			}

			// Check if the replay is cached in the database
			cachedReplay, err := GetCached(file.Name())
			if err != nil {
				fmt.Println("Error getting cached replay:", err)
			}

			if cachedReplay.Players != nil {
				fmt.Println("Using cached replay: " + file.Name())

				players := []Player{}

				for _, player := range cachedReplay.Players {
					players = append(players, Player{
						Name:      player.Name,
						Faction:   player.Faction,
						Team:      player.Team,
						SteamID:   fmt.Sprintf("%d", player.SteamID),
						ProfileID: player.ProfileID,
					})
				}

				replays = append(replays, Replay{
					Filename:  file.Name(),
					Version:   cachedReplay.Version,
					Timestamp: cachedReplay.Timestamp,
					Length:    cachedReplay.Length,
					Map: Map{Filename: cachedReplay.Map.Filename,
						LocalizedNameID:        cachedReplay.Map.LocalizedNameID,
						LocalizedDescriptionID: cachedReplay.Map.LocalizedDescriptionID,
					},
					Players: players})

				continue
			}

			replay, err := Parse(file.Name())
			if err != nil {
				fmt.Println("Error parsing replay:", err)
			}

			players := []Player{}

			for _, player := range replay.Players {
				players = append(players, Player{
					Name:      player.Name,
					Faction:   player.Faction,
					Team:      player.Team,
					SteamID:   fmt.Sprintf("%d", player.SteamID),
					ProfileID: player.ProfileID,
				})
			}

			replays = append(replays, Replay{
				Filename:  file.Name(),
				Version:   replay.Version,
				Timestamp: replay.Timestamp,
				Length:    replay.Length,
				Map: Map{Filename: replay.Map.Filename,
					LocalizedNameID:        replay.Map.LocalizedNameID,
					LocalizedDescriptionID: replay.Map.LocalizedDescriptionID,
				},
				Players: players})

			// Cache the replay in the database so we can skip parsing it next time
			Cache(file.Name(), replay)
		}
	}

	return replays
}
