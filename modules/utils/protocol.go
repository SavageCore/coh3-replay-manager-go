package utils

import (
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/sys/windows/registry"
)

func RegisterUrlProtocol() {
	// Get the current executable's path.
	exePath, err := os.Executable()
	if err != nil {
		panic(err)
	}

	// Get the executable file name
	exeName := filepath.Base(exePath)

	// Get the protocol identifier and the command to run.
	protocol := "coh3-replay-manager-go"
	command := fmt.Sprintf(`"%s" "%s"`, exePath, "%1")

	// Open the registry key for the current user.
	key, err := registry.OpenKey(registry.CURRENT_USER, "SOFTWARE\\Classes", registry.ALL_ACCESS)
	if err != nil {
		panic(err)
	}
	defer key.Close()

	// Using ReadSubKeyNames() to check if the protocol is already registered.
	subKeys, err := key.ReadSubKeyNames(-1)
	if err != nil {
		panic(err)
	}

	// If the protocol is already registered and command matches, print a message and exit.
	for _, subKey := range subKeys {
		if subKey == protocol {
			// Open the key for the protocol.
			coh3rmKey, err := registry.OpenKey(key, protocol, registry.ALL_ACCESS)
			if err != nil {
				panic(err)
			}
			defer coh3rmKey.Close()

			// Open the command key under the open key.
			openKey, err := registry.OpenKey(coh3rmKey, "shell\\open", registry.ALL_ACCESS)
			if err != nil {
				panic(err)
			}
			defer openKey.Close()

			// Open the command key under the open key.
			commandKey, err := registry.OpenKey(openKey, "command", registry.ALL_ACCESS)
			if err != nil {
				panic(err)
			}
			defer commandKey.Close()

			// Get the default value of the command key.
			commandValue, _, err := commandKey.GetStringValue("")
			if err != nil {
				panic(err)
			}

			// If the command matches, print a message and exit.
			if commandValue == command {
				fmt.Printf("Protocol '%s' is already registered with command '%s'\n", protocol, command)
				return
			}
		}
	}

	// Create the key for the protocol.
	coh3rmKey, _, err := registry.CreateKey(key, protocol, registry.ALL_ACCESS)
	if err != nil {
		panic(err)
	}
	defer coh3rmKey.Close()

	// Set the default value of the coh3rm key to the URL protocol description
	if err = coh3rmKey.SetStringValue("", "URL:"+protocol); err != nil {
		panic(err)
	}

	// Set the URL protocol handler to the path of your application
	if err = coh3rmKey.SetStringValue("URL Protocol", ""); err != nil {
		panic(err)
	}

	// Create the DefaultIcon key
	defaultIconKey, _, err := registry.CreateKey(coh3rmKey, "DefaultIcon", registry.ALL_ACCESS)
	if err != nil {
		panic(err)
	}

	// Set the default value of the DefaultIcon key to the file name such as "myapp.exe,1", removing the path
	if err = defaultIconKey.SetStringValue("", fmt.Sprintf("%s,1", exeName)); err != nil {
		panic(err)
	}

	// Create the shell key for the protocol.
	shellKey, _, err := registry.CreateKey(coh3rmKey, "shell", registry.ALL_ACCESS)
	if err != nil {
		panic(err)
	}
	defer shellKey.Close()

	// Create the open key under the shell key.
	openKey, _, err := registry.CreateKey(shellKey, "open", registry.ALL_ACCESS)
	if err != nil {
		panic(err)
	}
	defer openKey.Close()

	// Create the command key under the open key
	commandKey, _, err := registry.CreateKey(openKey, "command", registry.ALL_ACCESS)
	if err != nil {
		panic(err)
	}
	defer commandKey.Close()

	// Set the default value of the command key to the command to run.
	if err = commandKey.SetStringValue("", command); err != nil {
		panic(err)
	}

	// Print a message indicating that the protocol handler was registered.
	fmt.Printf("Registered protocol '%s' with command '%s'\n", protocol, command)
}
