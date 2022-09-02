package window

import (
	"os"
	"strings"

	"github.com/Pauloo27/gotroller"
	"github.com/godbus/dbus/v5"
	"github.com/gotk3/gotk3/gtk"
)

const MPRIS_PREFIX = "org.mpris.MediaPlayer2."

func createPlayerSelector() *gtk.ComboBoxText {
	box, err := gtk.ComboBoxTextNew()
	handleError(err)

	names, err := gotroller.ListPlayersName()
	handleError(err)

	box.Append("auto", "auto")
	for _, fullName := range names {
		box.Append(fullName, strings.TrimPrefix(fullName, MPRIS_PREFIX))
	}
	box.Append("disabled", "disabled")

	currentPlayer := ""
	onUpdate(func(metadata map[string]dbus.Variant) {
		if currentPlayer == "" {
			currentPlayer = playerInstance.GetName()
			box.SetActiveID(currentPlayer)
		}
	})

	box.Connect("changed", func() {
		id := box.GetActiveID()
		switch id {
		case currentPlayer:
			return
		case "auto":
			gotroller.RemovePreferedPlayerName()
		case "disabled":
			gotroller.HideGotroller()
		default:
			gotroller.SetPreferedPlayerName(id)
		}
		os.Exit(0)
	})

	return box
}
