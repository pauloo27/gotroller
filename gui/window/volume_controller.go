package window

import (
	"github.com/godbus/dbus/v5"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

func createVolumeItem() (*gtk.Box, *gtk.Scale) {
	container, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	handleError(err)

	volume, err := gtk.ScaleNewWithRange(gtk.ORIENTATION_VERTICAL, 0.0, 1.0, 0.01)
	handleError(err)

	volume.SetDrawValue(false)
	volume.SetInverted(true)

	icon, err := gtk.ImageNewFromIconName("audio-volume-high", gtk.ICON_SIZE_BUTTON)
	handleError(err)

	container.PackStart(volume, true, true, 3)
	container.PackEnd(icon, false, false, 3)

	return container, volume
}

func createVolumeController() *gtk.Box {
	mainContainer, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	handleError(err)

	container, playerVolume := createVolumeItem()
	mainContainer.PackStart(container, true, true, 0)

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
