package gotroller

import (
	"io/ioutil"
	"os"
	"strings"
)

const PREFERED_PLAYER_STORE_PATH = "/dev/shm/gotroller-player.txt"

func SetPreferedPlayerName(name string) error {
	data := []byte(name)
	return ioutil.WriteFile(PREFERED_PLAYER_STORE_PATH, data, 0644)
}

func GetPreferedPlayerName() (string, error) {
	buffer, err := ioutil.ReadFile(PREFERED_PLAYER_STORE_PATH)
	if err != nil {
		return "", nil
	}
	return strings.TrimSuffix(string(buffer), "\n"), nil
}

func RemovePreferedPlayerName() error {
	return os.Remove(PREFERED_PLAYER_STORE_PATH)
}

func HideGotroller() error {
	return SetPreferedPlayerName("Disabled")
}
