package window

import (
	"strings"

	"github.com/Pauloo27/gotroller/cli/utils"
	"github.com/godbus/dbus/v5"
	"github.com/gotk3/gotk3/gtk"
)

func createArtist() *gtk.Label {
	artistLabel, err := gtk.LabelNew("")
	handleError(err)

	onUpdate(func(metadata map[string]dbus.Variant) {
		var displayArtist string
		switch artist := metadata["xesam:artist"].Value(); artist.(type) {
		case []string:
			displayArtist = strings.Join(artist.([]string), ", ")
		case string:
			displayArtist = artist.(string)
		case nil:
			displayArtist = "--"
		}

		artistLabel.SetTooltipText(displayArtist)
		displayArtist = utils.EnforceSize(displayArtist, maxArtistSize)
		artistLabel.SetText(displayArtist)
		ctx, err := artistLabel.GetStyleContext()
		handleError(err)
		ctx.AddClass("artist")
	})

	return artistLabel
}
