package window

import (
	"github.com/pauloo27/gotroller/cli/utils"
	"github.com/godbus/dbus/v5"
	"github.com/gotk3/gotk3/gtk"
)

func createTitle() *gtk.Label {
	titleLabel, err := gtk.LabelNew("")
	handleError(err)

	onUpdate(func(metadata map[string]dbus.Variant) {
		var displayTitle string
		switch title := metadata["xesam:title"].Value(); title.(type) {
		case string:
			displayTitle = title.(string)
		case nil:
			displayTitle = "--"
		}

		titleLabel.SetTooltipText(displayTitle)
		displayTitle = utils.EnforceSize(displayTitle, maxTitleSize)
		titleLabel.SetText(displayTitle)
		ctx, err := titleLabel.GetStyleContext()
		handleError(err)
		ctx.AddClass("title")
	})
	return titleLabel
}
