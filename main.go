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
	return fmt.Sprintf("%%{A%d:%s:}%s%%{A}", a.Index, a.Command, a.Display)
}

func printToPolybar(player *mpris.Player) {
	icon := stoppedIcon
	if player == nil {
		fmt.Printf("%s\n", icon)
		return
	}

	identity := player.GetIdentity()

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

	prevButton := PolybarActionButton{
		Index:   1,
		Display: "",
		Command: fmt.Sprintf("playerctl -p %s previous", identity),
	}

	togglePauseButton := PolybarActionButton{
		Index:   1,
		Display: icon,
		Command: fmt.Sprintf("playerctl -p %s %s", identity, playPause),
	}

	nextButton := PolybarActionButton{
		Index:   1,
		Display: "",
		Command: fmt.Sprintf("playerctl -p %s next", identity),
	}

	volumeUpButton := PolybarActionButton{
		Index:   4,
		Display: fmt.Sprintf(" %.f%%", player.GetVolume()*100),
		Command: fmt.Sprintf("playerctl -p %s volume 0.05+", identity),
	}

	volumeButton := PolybarActionButton{
		Index:   5,
		Display: volumeUpButton.String(),
		Command: fmt.Sprintf("playerctl -p %s volume 0.05-", identity),
	}

	metadata := player.GetMetadata()

	title := metadata["xesam:title"].Value().(string)

	displayName := title

	if metadata["xesam:artist"].Value() != nil {
		artist := metadata["xesam:artist"].Value().([]string)[0]
		displayName += " from "
		displayName += artist
	}

	fmt.Printf("%s %s %s %s %s\n", displayName, prevButton.String(), togglePauseButton.String(), nextButton.String(), volumeButton.String())
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
