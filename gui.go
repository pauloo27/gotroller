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
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
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

	mainBox, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
	handleFatal(err)

	contentBox, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 1)
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

	metadata, err := player.GetMetadata()

	title := ""
	if err == nil {
		titleData := metadata["xesam:title"].Value()
		if titleData != nil {
			title = titleData.(string)
		}
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

	canSeeLength := true

	updateProgress := func() {
		length, err := player.GetLength()

		if err != nil {
			canSeeLength = false
		}

		position, err := player.GetPosition()

		if err != nil {
			canSeeLength = false
		}

		percent := 100.0 * position / length

		progressBar.SetValue(percent)
		lastPosition = percent
	}

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
			length, _ := player.GetLength()
			player.SetPosition(currentPosition * length / 100)
		}
	})

	closeButton, err := gtk.ButtonNewWithLabel("Close")
	handleFatal(err)

	closeButton.Connect("clicked", func() {
		fmt.Println("Closed!")
		os.Exit(0)
	})

	mainBox.PackEnd(contentBox, true, true, 0)

	if enabled {
		contentBox.PackStart(label, true, true, 1)
		if canSeeLength {
			contentBox.PackStart(progressBar, true, true, 1)
		}
	}

	artUrl := ""
	artUrlEntry := metadata["mpris:artUrl"].Value()
	if artUrlEntry != nil {
		artUrl = artUrlEntry.(string)
	}

	if artUrl != "" {
		if strings.HasPrefix(artUrl, "http") {
			setupCacheFolder()
			artUrl, err = downloadAlbumArt(metadata["xesam:url"].Value().(string), artUrl)
			if err != nil {
				fmt.Println("Cannot download album art")
			}
		} else if strings.HasPrefix(artUrl, "file://") {
			artUrl = strings.TrimPrefix(artUrl, "file://")
		}

		// Check one more time because it may change in the if above
		if artUrl != "" {

			albumImagePix, err := gdk.PixbufNewFromFileAtScale(artUrl, 150, 150, true)
			if err != nil {
				fmt.Println("Cannot load album art")
			} else {
				image, err := gtk.ImageNewFromPixbuf(albumImagePix)
				if err != nil {
					fmt.Println("Cannot load album art (2)")
				} else {
					mainBox.PackStart(image, true, true, 0)
				}
			}
		}
	}

	bottomBox, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 1)
	handleFatal(err)
	bottomBox.PackStart(comboBox, true, true, 1)
	bottomBox.PackEnd(closeButton, true, true, 1)

	contentBox.PackEnd(bottomBox, true, true, 1)

	win.Add(mainBox)

	win.SetDefaultSize(450, 150)
	win.ShowAll()
	gtk.Main()
}
