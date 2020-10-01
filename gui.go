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

func appendAlbumArt(parent *gtk.Box, metadata map[string]dbus.Variant) {
	var err error

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
					parent.PackStart(image, true, true, 0)
				}
			}
		}
	}
}

func appendControllers(parent *gtk.Box, player *mpris.Player) {
	prevButton, err := gtk.ButtonNewFromIconName("media-seek-backward", gtk.ICON_SIZE_MENU)
	handleFatal(err)
	prevButton.Connect("clicked", func() {
		player.Previous()
		os.Exit(0)
	})

	buttonIcon := "media-playback-pause"

	playback, err := player.GetPlaybackStatus()
	handleFatal(err)

	if playback != mpris.PlaybackPlaying {
		buttonIcon = "media-playback-start"
	}

	playPauseButton, err := gtk.ButtonNewFromIconName(buttonIcon, gtk.ICON_SIZE_MENU)
	handleFatal(err)
	playPauseButton.Connect("clicked", func() {
		player.PlayPause()
		os.Exit(0)
	})

	nextButton, err := gtk.ButtonNewFromIconName("media-seek-forward", gtk.ICON_SIZE_MENU)
	handleFatal(err)
	nextButton.Connect("clicked", func() {
		player.Next()
		os.Exit(0)
	})

	buttonBox, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
	buttonBox.PackStart(prevButton, false, false, 1)
	buttonBox.PackStart(playPauseButton, false, false, 1)
	buttonBox.PackStart(nextButton, false, false, 1)
	buttonBox.SetHAlign(gtk.ALIGN_CENTER)

	parent.PackStart(buttonBox, true, true, 1)
}

func appendCloseButton(parent *gtk.Box) {
	closeButton, err := gtk.ButtonNewWithLabel("Close")
	handleFatal(err)

	closeButton.Connect("clicked", func() {
		fmt.Println("Closed!")
		os.Exit(0)
	})

	parent.PackEnd(closeButton, true, true, 1)
}

func appendProgressBar(parent *gtk.Box, player *mpris.Player) {
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

	if canSeeLength {
		parent.PackStart(progressBar, true, true, 1)
	}
}

func appendPlayerSelector(parent *gtk.Box, players []string) (string, bool) {
	comboBox, err := gtk.ComboBoxTextNew()

	handleFatal(err)

	playerName := players[0]

	comboBox.Append("Disable", "Disable")

	selectedPlayer := getSelectedPlayer()
	enabled := selectedPlayer != "Disable"

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

	parent.PackStart(comboBox, true, true, 1)

	return playerName, enabled
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

	bottomBox, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 1)
	handleFatal(err)

	players, err := mpris.List(conn)
	handleFatal(err)

	if len(players) == 0 {
		fmt.Println("No players found =(")
		os.Exit(0)
		return
	}

	playerName, enabled := appendPlayerSelector(bottomBox, players)

	player := mpris.New(conn, playerName)

	metadata, err := player.GetMetadata()

	title := ""
	if err == nil {
		titleData := metadata["xesam:title"].Value()
		if titleData != nil {
			title = titleData.(string)
		}
	}

	artist := metadata["xesam:artist"].Value()
	if artist != nil {
		title = fmt.Sprintf("%s - %s", artist.([]string)[0], title)
	}

	if len(title) > 35 {
		title = title[0:32] + "..."
	}

	label, err := gtk.LabelNew(title)
	handleFatal(err)

	mainBox.PackEnd(contentBox, true, true, 0)

	if enabled {
		contentBox.PackStart(label, true, true, 1)
		appendProgressBar(contentBox, player)

		appendControllers(contentBox, player)

		appendAlbumArt(mainBox, metadata)
	}

	appendCloseButton(bottomBox)

	contentBox.PackEnd(bottomBox, true, true, 1)

	win.Add(mainBox)

	win.SetDefaultSize(450, 150)
	win.ShowAll()
	gtk.Main()
}
