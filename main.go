package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

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

func printToPolybar(name string, player *mpris.Player) {
	playerButton := PolybarActionButton{
		Index:   1,
		Display: "",
		Command: "gotroller player | dmenu > /dev/shm/gotroller-player.txt",
	}

	icon := stoppedIcon
	if player == nil {
		fmt.Printf("%s %s\n", playerButton.String(), icon)
		return
	}

	identity := strings.Split(name, ".")[3]

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

	title := metadata["xesam:title"].Value().(string)

	if len(title) > 35 {
		title = title[0:32] + "..."
	}

	displayName := title

	if metadata["xesam:artist"].Value() != nil {
		var artistName string

		artist := metadata["xesam:artist"].Value()

		switch artist.(type) {
		case string:
			artistName = artist.(string)
		case []string:
			artistName = artist.([]string)[0]
		}

		displayName += " from "
		displayName += artistName
	}

	fmt.Printf("%s %s %s %s %s %s\n", playerButton.String(), displayName, prevButton.String(), togglePauseButton.String(), nextButton.String(), volumeButton.String())
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

	if len(os.Args) >= 2 {
		if os.Args[1] == "player" {
			fmt.Println("Disable")
			for _, player := range names {
				fmt.Println(player)
			}
			return
		}
	}

	if len(names) == 0 {
		printToPolybar("", nil)
		return
	}

	selectedPlayer := ""

	buffer, err := ioutil.ReadFile("/dev/shm/gotroller-player.txt")
	if err == nil {
		selectedPlayer = strings.TrimSuffix(string(buffer), "\n")
	}

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
