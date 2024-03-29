package main

import (
	"fmt"
	"os"

	"github.com/pauloo27/go-mpris"
	"github.com/pauloo27/gotroller"
	"github.com/pauloo27/gotroller/cli/dmenu"
	"github.com/pauloo27/gotroller/cli/operation"
	"github.com/pauloo27/gotroller/cli/polybar"
	"github.com/pauloo27/gotroller/cli/waybar"
)

type Mode func()

var modes = map[string]Mode{
	"polybar-dmenu": polybar.WithDmenu,
	"polybar-gui":   polybar.WithGUI,
	"waybar":    waybar.Start,
	"set-player":    setPlayer,
	"volume":        operation.SetVolume,
	"position":      operation.SetPosition,
	"dmenu-select":  dmenu.Select,
	"play-pause":    playPause,
	"next":          next,
	"prev":          prev,
}

func mustLoadPlayer() *mpris.Player {
	player, err := gotroller.GetBestPlayer()
	if err != nil {
		panic(err)
	}

	if player == nil {
		fmt.Println("No players found")
		os.Exit(-1)
	}

	return player
}

func playPause() {
	player := mustLoadPlayer()
	err := player.PlayPause()
	if err != nil {
		panic(err)
	}
}

func next() {
	player := mustLoadPlayer()
	err := player.Next()
	if err != nil {
		panic(err)
	}
}

func prev() {
	player := mustLoadPlayer()
	err := player.Previous()
	if err != nil {
		panic(err)
	}
}

func setPlayer() {
	if len(os.Args) < 3 {
		fmt.Println("Missing player identity")
		os.Exit(-1)
	}
	identity := os.Args[2]
	if identity == "auto" {
		gotroller.RemovePreferedPlayerName()
	} else {
		gotroller.SetPreferedPlayerName(identity)
	}
}

func listModes() {
	fmt.Print("Valid modes: ")

	var prefix string
	for modeName := range modes {
		fmt.Printf("%s%s", prefix, modeName)
		if prefix == "" {
			prefix = ", "
		}
	}

	fmt.Println()
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Missing mode.")
		listModes()
		os.Exit(-1)
	}

	mode, ok := modes[os.Args[1]]
	if !ok {
		fmt.Println("Invalid mode.")
		listModes()
		os.Exit(-1)
	}

	mode()
}
