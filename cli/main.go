package main

import (
	"fmt"
	"os"

	"github.com/Pauloo27/gotroller"
	"github.com/Pauloo27/gotroller/cli/dmenu"
	"github.com/Pauloo27/gotroller/cli/polybar"
)

type Mode func()

var modes = map[string]Mode{
	"polybar-dmenu": polybar.WithDmenu,
	"polybar-gui":   polybar.WithGUI,
	"set-player":    setPlayer,
	"dmenu-select":  dmenu.Select,
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
