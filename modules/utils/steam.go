package utils

import "golang.org/x/sys/windows/registry"

func GetSteamPath() string {
	key, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall\Steam App 1677280`, registry.QUERY_VALUE)
	if err != nil {
		panic(err)
	}
	defer key.Close()

	value, _, err := key.GetStringValue("UninstallString")
	if err != nil {
		panic(err)
	}

	// Sample UninstallString value:
	// "C:\Program Files (x86)\Steam\steam.exe" steam://uninstall/1677280

	// Return the path to the Steam executable
	return value[1 : len(value)-len("\" steam://uninstall/1677280")]
}
