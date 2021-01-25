package window

import (
	"strings"

	"github.com/Pauloo27/go-mpris"
	"github.com/godbus/dbus/v5"
	"github.com/gotk3/gotk3/gtk"
)

func createArtist() *gtk.Label {
	artistLabel, err := gtk.LabelNew("")
	handleError(err)

	onUpdate(func(player *mpris.Player, metadata map[string]dbus.Variant) {
		var displayArtist string
		switch artist := metadata["xesam:artist"].Value(); artist.(type) {
		case []string:
			displayArtist = strings.Join(artist.([]string), ", ")
		case string:
			displayArtist = artist.(string)
		case nil:
			displayArtist = "--"
		}

		artistLabel.SetText(displayArtist)
	})

	return artistLabel
}
