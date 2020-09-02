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
	"github.com/gotk3/gotk3/gtk"
  "github.com/gotk3/gotk3/glib"
)

func getSelectedPlayer() string {
	selectedPlayer := ""

	buffer, err := ioutil.ReadFile(selectedPlayerPath)
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

	win.SetPosition(gtk.WIN_POS_MOUSE)

	_, err = win.Connect("destroy", func() {
		fmt.Println("Closed")
		gtk.MainQuit()
	})
	handleFatal(err)

	win.SetTitle("Gotroller")
	grid, err := gtk.GridNew()
	handleFatal(err)

	selectedPlayer := getSelectedPlayer()
	enabled := selectedPlayer != "Disable"

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

	comboBox.Append("Disable", "Disable")

	for _, player := range players {
		comboBox.Append(player, player)
		if player == selectedPlayer {
			playerName = player
		}
	}

	if enabled {
		comboBox.SetActiveID(playerName)
	} else {
		comboBox.SetActiveID("Disable")
	}

	comboBox.Connect("changed", func() {
		newSelection := comboBox.GetActiveText()
		go func() {
			data := []byte(newSelection)
			err := ioutil.WriteFile(selectedPlayerPath, data, 0644)
			handleFatal(err)
			fmt.Println("Player changed to", newSelection)
			os.Exit(0)
		}()
	})

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

	canSeeLength := metadata["mpris:length"].Value() != nil

	if canSeeLength {
		go func() {
			for {
				glib.IdleAdd(updateProgress)
				time.Sleep(1 * time.Second)
			}
		}()
	}

	progressBar.Connect("value-changed", func() {
		if lastPosition == 0 {
			return
		}
		currentPosition := progressBar.GetValue()
		if currentPosition-lastPosition >= 1 || currentPosition-lastPosition <= -1 {
			player.SetPosition(currentPosition * player.GetLength() / 100)
		}
	})

	closeButton, err := gtk.ButtonNewWithLabel("Close")
	handleFatal(err)

	closeButton.Connect("clicked", func() {
		fmt.Println("Closed!")
		os.Exit(0)
	})

	grid.Attach(comboBox, 0, 0, 1, 1)
	if enabled {
		grid.Attach(label, 0, 1, 1, 1)
		if canSeeLength {
			grid.Attach(progressBar, 0, 2, 10, 1)
		}
	}
	grid.Attach(closeButton, 5, 3, 1, 1)

	win.Add(grid)

	win.SetDefaultSize(400, 100)
	win.ShowAll()
	gtk.Main()
}
