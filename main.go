package main

import (
	"fmt"
	"log"

	"github.com/Pauloo27/go-mpris"
	"github.com/godbus/dbus"
)

const pauseIcon = ""
const playIcon = ""
const stoppedIcon = ""

func printToPolybar(player *mpris.Player) {
	icon := stoppedIcon
	if player == nil {
		fmt.Printf("%s\n", icon)
		return
	}

	status := player.GetPlaybackStatus()

	if status == mpris.PlaybackPlaying {
		icon = pauseIcon
	} else if status == mpris.PlaybackPaused {
		icon = playIcon
	} else if status == mpris.PlaybackStopped {
		icon = stoppedIcon
	} else {
		log.Fatalf("Invalid playback status %s / %s", status, mpris.PlaybackPaused)
	}
	fmt.Printf("%s\n", icon)
}

func main() {
	conn, err := dbus.SessionBus()
	if err != nil {
		panic(err)
	}
	names, err := mpris.List(conn)
	if err != nil {
		panic(err)
	}
	if len(names) == 0 {
		printToPolybar(nil)
		return
	}
	player := mpris.New(conn, names[0])
	printToPolybar(player)
}
