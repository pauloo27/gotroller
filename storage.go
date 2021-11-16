package gotroller

import (
	"io/ioutil"
	"os"
	"strings"

	"github.com/fsnotify/fsnotify"
)

const PREFERED_PLAYER_STORE_PATH = "/dev/shm/gotroller-player.txt"

func SetPreferedPlayerName(name string) error {
	data := []byte(name)
	return ioutil.WriteFile(PREFERED_PLAYER_STORE_PATH, data, 0644)
}

func GetPreferedPlayerName() (string, error) {
	buffer, err := ioutil.ReadFile(PREFERED_PLAYER_STORE_PATH)
	if err != nil {
		if os.IsNotExist(err) {
			RemovePreferedPlayerName() // create a file with empty content
		}
		return "", nil
	}
	return strings.TrimSuffix(string(buffer), "\n"), nil
}

func RemovePreferedPlayerName() error {
	return SetPreferedPlayerName("")
}

func HideGotroller() error {
	return SetPreferedPlayerName("Disabled")
}

func ListenToChanges(ch chan fsnotify.Event) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	go func() {
		for {
			event, ok := <-watcher.Events
			if !ok {
				return
			}
			ch <- event
		}
	}()

	return watcher.Add(PREFERED_PLAYER_STORE_PATH)
}
