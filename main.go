package main

import (
	"fmt"
	"os"

	"github.com/Pauloo27/go-mpris"
	"github.com/godbus/dbus/v5"
)

const (
	pausedIcon  = ""
	playingIcon = ""
	stoppedIcon = ""

	selectedPlayerPath = "/dev/shm/gotroller-player.txt"
)

func main() {
	conn, err := dbus.SessionBus()
	if err != nil {
		panic(err)
	}
	names, err := mpris.List(conn)

	if err != nil {
		panic(err)
	}

	if len(os.Args) >= 2 {
		if os.Args[1] == "player" {
			fmt.Println("Disable")
			for _, player := range names {
				fmt.Println(player)
			}
			return
		} else if os.Args[1] == "gui" {
			fmt.Println("Opening GUI")
			showGUI(conn)
			return
		}
	}

	if len(names) == 0 {
		printToPolybar("", nil)
		return
	}

	selectedPlayer := getSelectedPlayer()

	if selectedPlayer == "Disable" {
		printToPolybar("Disable", nil)
		return
	}

	var playerName string
	if selectedPlayer != "" {
		for _, name := range names {
			if name == selectedPlayer {
				playerName = name
				break
			}
		}
	}

	if playerName == "" {
		playerName = names[0]
	}

	player := mpris.New(conn, playerName)
	printToPolybar(playerName, player)
}
