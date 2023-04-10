package game

import (
	"golang.org/x/sys/windows/registry"
)

func GetInstallPath() string {
	// Get the install path of the game from the registry key:
	// HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall\Steam App 1677280
	key, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall\Steam App 1677280`, registry.QUERY_VALUE)
	if err != nil {
		panic(err)
	}
	defer key.Close()

	value, _, err := key.GetStringValue("InstallLocation")
	if err != nil {
		panic(err)
	}

	return value
}
