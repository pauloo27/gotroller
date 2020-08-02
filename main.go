package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/Pauloo27/go-mpris"
	"github.com/godbus/dbus"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

const pausedIcon = ""
const playingIcon = ""
const stoppedIcon = ""

func getSelectedPlayer() string {
	selectedPlayer := ""

	buffer, err := ioutil.ReadFile("/dev/shm/gotroller-player.txt")
	if err == nil {
		selectedPlayer = strings.TrimSuffix(string(buffer), "\n")
	}

	return selectedPlayer
}

func handleFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func showGUI(conn *dbus.Conn) {
	gtk.Init(nil)
	win, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	handleFatal(err)

	_, err = win.Connect("destroy", func() {
		fmt.Println("Closed")
		gtk.MainQuit()
	})
	handleFatal(err)

	win.SetTitle("Gotroller")
	grid, err := gtk.GridNew()
	handleFatal(err)

	selectedPlayer := getSelectedPlayer()

	players, err := mpris.List(conn)
	handleFatal(err)

	if len(players) == 0 {
		fmt.Println("No players found =(")
		os.Exit(0)
		return
	}

	comboBox, err := gtk.ComboBoxTextNew()

	handleFatal(err)

	playerName := players[0]

	for _, player := range players {
		comboBox.Append(player, player)
		if player == selectedPlayer {
			playerName = player
		}
	}

	comboBox.SetActiveID(playerName)
	player := mpris.New(conn, playerName)

	metadata := player.GetMetadata()

	title := ""
	titleData := metadata["xesam:title"].Value()
	if titleData != nil {
		title = titleData.(string)
	}

	if len(title) > 35 {
		title = title[0:32] + "..."
	}

	label, err := gtk.LabelNew(title)
	handleFatal(err)

	progressBar, err := gtk.ScaleNewWithRange(gtk.ORIENTATION_HORIZONTAL, 0, 100, 1)
	handleFatal(err)

	progressBar.SetHExpand(true)
	progressBar.SetDrawValue(false)

	var lastPosition float64

	updateProgress := func() {
		length := player.GetLength()
		position := float64(player.GetPosition())

		percent := 100.0 * position / length

		progressBar.SetValue(percent)
		lastPosition = percent
	}

	go func() {
		for {
			glib.IdleAdd(updateProgress)
			time.Sleep(1 * time.Second)
		}
	}()

	progressBar.Connect("value-changed", func() {
		if lastPosition == 0 {
			return
		}
		currentPosition := progressBar.GetValue()
		if currentPosition-lastPosition >= 1 || currentPosition-lastPosition <= -1 {
			player.SetPosition(currentPosition * player.GetLength() / 100)
		}
	})

	grid.Attach(comboBox, 0, 0, 1, 1)
	grid.Attach(label, 0, 1, 1, 1)
	grid.Attach(progressBar, 0, 2, 10, 1)

	win.Add(grid)

	win.SetDefaultSize(400, 150)
	win.ShowAll()
	gtk.Main()
}

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

	title := ""
	titleData := metadata["xesam:title"].Value()
	if titleData != nil {
		title = titleData.(string)
	}

	if len(title) > 35 {
		title = title[0:32] + "..."
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
