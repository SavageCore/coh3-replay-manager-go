package game

import (
	"fmt"
	"strings"

	"github.com/StackExchange/wmi"
)

type CIM_DataFile struct {
	Version string
}

func GetGameVersion() string {
	installPath := GetInstallPath()

	// Get the game version from the game executable
	exePath := installPath + "\\RelicCoH3.exe"

	// Replace all backslashes with double backslashes
	exePath = strings.ReplaceAll(exePath, "\\", "\\\\")

	var result []CIM_DataFile
	query := fmt.Sprintf("SELECT Version FROM CIM_DataFile WHERE Name='%s'", exePath)
	err := wmi.Query(query, &result)
	if err != nil {
		panic(err)
	}

	var version string

	if len(result) > 0 {
		// fileVersionString := strings.Split(strings.TrimSpace(result[0].Version), "=")[1]
		// version = strings.Split(fileVersionString, ".")[2]
		parts := strings.Split(result[0].Version, ".")
		version = parts[2]
	} else {
		fmt.Println("File not found or version information not available.")
	}

	return string(version)
}
