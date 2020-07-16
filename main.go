package main

import (
	"fmt"
	"log"

	"github.com/Pauloo27/go-mpris"
	"github.com/godbus/dbus"
)

const pausedIcon = ""
const playingIcon = ""
const stoppedIcon = ""

type PolybarActionButton struct {
	Index            uint
	Display, Command string
}

func (a PolybarActionButton) String() string {
	return fmt.Sprintf("%%{A%d:%s:} %s %%{A}", a.Index, a.Command, a.Display)
}

func printToPolybar(player *mpris.Player) {
	icon := stoppedIcon
	if player == nil {
		fmt.Printf("%s\n", icon)
		return
	}

	status := player.GetPlaybackStatus()
	var playPause string

	if status == mpris.PlaybackPlaying {
		icon = playingIcon
		playPause = "pause"
	} else if status == mpris.PlaybackPaused {
		icon = pausedIcon
		playPause = "play"
	} else if status == mpris.PlaybackStopped {
		icon = stoppedIcon
		playPause = "play"
	} else {
		log.Fatalf("Invalid playback status %s / %s", status, mpris.PlaybackPaused)
	}

	metadata := player.GetMetadata()

	togglePauseButton := PolybarActionButton{
		Index:   1,
		Display: icon,
		Command: fmt.Sprintf("playerctl -p %s %s", player.GetIdentity(), playPause),
	}

	fmt.Printf("%s %v\n", togglePauseButton.String(), metadata["xesam:title"].Value())
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
