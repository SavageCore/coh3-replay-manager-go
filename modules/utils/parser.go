package utils

import (
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Player struct {
	Name    string
	Faction string
}

func ParseReplay(filename string) {
	fmt.Println("")
	fmt.Println("")
	fmt.Println("Parsing replay: " + filename)

	user := GetUsername()
	replayDir := filepath.Join(user, "Documents", "My Games", "Company of Heroes 3", "playback", "replays")
	replayFilePath := filepath.Join(replayDir, filename)

	fi, err := os.Stat(replayFilePath)
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.Open(replayFilePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
	}
	defer file.Close()

	data := make([]byte, fi.Size())
	_, err = file.Read(data)
	if err != nil {
		log.Fatal(err)
	}

	dataHex := hex.EncodeToString(data)

	// Get Players
	playerSearch := hex.EncodeToString([]byte("default_ai_personality"))

	var players []Player
	pos := strings.Index(dataHex, playerSearch)

	// If pos == -1, there are no players in the replay, so we can stop here
	if pos == -1 {
		fmt.Println("No players in replay - likely a campaign replay")
		return
	}

	for pos != -1 {
		// Get the player's team

		// Get the player's faction
		curPos := pos - 2

		for dataHex[curPos:curPos+2] != "22" {
			curPos = curPos - 2
		}

		var faction string

		curPos -= 4
		for !strings.HasPrefix(faction, "0000") {
			curPos -= 2
			faction = dataHex[curPos:curPos+2] + faction
		}

		faction = strings.ReplaceAll(faction, "00", "")
		faction = string(FromHex(faction))

		// Get the player's name
		var playerName string
		curPos -= 24
		for !strings.HasPrefix(playerName, "0000") {
			curPos -= 2
			playerName = dataHex[curPos:curPos+2] + playerName
		}

		playerName = strings.ReplaceAll(playerName, "00", "")
		playerName = string(FromHex(playerName))

		players = append(players, Player{
			Name:    playerName,
			Faction: faction,
		})

		// Get the next player
		dataHex = dataHex[pos+len(playerSearch):]
		pos = strings.Index(dataHex, playerSearch)
	}

	// For each players, print their name and faction
	for _, player := range players {
		fmt.Println("Player: " + player.Name)
		fmt.Println("Faction: " + player.Faction)
	}

	// Get Map
	mapSearch := hex.EncodeToString([]byte("data:"))
	pos = strings.Index(dataHex, mapSearch)

	var mapName string

	if pos != -1 {
		pos += 10
		var mapStr []byte

		for dataHex[pos:pos+2] != "09" {
			mapStr = append(mapStr, dataHex[pos:pos+2]...)
			pos += 2
		}

		mapName = string(FromHex(string(mapStr)))
	}

	fmt.Println("Map: " + mapName)
}

func FromHex(s string) []byte {
	b, err := hex.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return b
}
