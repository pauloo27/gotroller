package window

import (
	"github.com/Pauloo27/go-mpris"
	"github.com/godbus/dbus/v5"
	"github.com/gotk3/gotk3/gtk"
)

func createTitle() *gtk.Label {
	titleLabel, err := gtk.LabelNew("")
	handleError(err)

	onUpdate(func(player *mpris.Player, metadata map[string]dbus.Variant) {
		var displayTitle string
		switch title := metadata["xesam:title"].Value(); title.(type) {
		case string:
			displayTitle = title.(string)
		case nil:
			displayTitle = "--"
		}

		titleLabel.SetText(displayTitle)
	})

	return titleLabel
}
