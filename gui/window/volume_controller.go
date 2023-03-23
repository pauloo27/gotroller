package window

import (
	"github.com/godbus/dbus/v5"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

func createVolumeController() *gtk.Box {
	mainContainer, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	handleError(err)

	playerVolume, err := gtk.ScaleNewWithRange(gtk.ORIENTATION_VERTICAL, 0.0, 1.0, 0.01)
	handleError(err)
	playerVolume.SetInverted(true)

	mainContainer.PackStart(playerVolume, true, true, 0)

	playerVolume.SetDrawValue(false)
	expectedVolume := 0.0

	playerVolume.Connect("value-changed", func() {
		volume := playerVolume.GetValue()
		if volume != expectedVolume {
			playerInstance.SetVolume(volume)
		}
	})

	onUpdate(func(metadata map[string]dbus.Variant) {
		volume, err := playerInstance.GetVolume()
		if err != nil {
			return
		}

		expectedVolume = volume
		glib.IdleAdd(func() {
			playerVolume.SetValue(expectedVolume)
		})
	})

	return mainContainer
}
