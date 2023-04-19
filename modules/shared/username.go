package shared

import "os"

func GetUsername() string {
	user, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	return user
}
