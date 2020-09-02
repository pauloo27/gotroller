package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/Pauloo27/go-mpris"
)

type PolybarActionButton struct {
	Index            uint
	Display, Command string
}

func (a PolybarActionButton) String() string {
	return fmt.Sprintf("%%{A%d:%s:}%s%%{A}", a.Index, a.Command, a.Display)
}

func printToPolybar(name string, player *mpris.Player) {
	playerButton := PolybarActionButton{
		Index:   1,
		Display: "",
		Command: "gotroller gui",
	}

	icon := stoppedIcon
	if player == nil {
		fmt.Printf("%s %s\n", playerButton.String(), icon)
		return
	}

	identity := strings.TrimPrefix(name, "org.mpris.MediaPlayer2.")

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

	restartButton := PolybarActionButton{
		Index:   3,
		Display: "",
		Command: fmt.Sprintf("playerctl -p %s position 0", identity),
	}

	prevButton := PolybarActionButton{
		Index:   1,
		Display: restartButton.String(),
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

	title := ""
	titleData := metadata["xesam:title"].Value()
	if titleData != nil {
		title = titleData.(string)
	}

	if len(title) > 25 {
		title = title[0:22] + "..."
	}

	displayName := title

	if metadata["xesam:artist"].Value() != nil {
		var artistName string

		artist := metadata["xesam:artist"].Value()

		switch artist := artist.(type) {
		case string:
			artistName = artist
		case []string:
			artistName = artist[0]
		}

		displayName += " from "
		displayName += artistName
	}

	fmt.Printf("%s %s %s %s %s %s\n", playerButton.String(), displayName, prevButton.String(), togglePauseButton.String(), nextButton.String(), volumeButton.String())
}
