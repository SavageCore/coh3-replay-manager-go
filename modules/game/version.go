package game

import (
	"fmt"
	"os/exec"
	"strings"
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

	// wmic datafile where name="C:\\Program Files (x86)\\Steam\\steamapps\\common\\Company of Heroes 3\\RelicCoH3.exe" get Version /value
	cmd := exec.Command("wmic", "datafile", "where", "name=\""+exePath+"\"", "get", "Version", "/value")
	out, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
	}

	// Split the output into lines
	lines := strings.Split(string(out), "\r\n")
	parts := strings.Split(lines[2], "=")

	var version string

	if len(parts) > 0 {
		// parts[1] has the format 5.1.10907.0 we need the 5 digit version
		version = strings.Split(parts[1], ".")[2]
	} else {
		fmt.Println("File not found or version information not available.")
	}

	return string(version)
}
