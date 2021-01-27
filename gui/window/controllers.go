package window

import (
	"github.com/Pauloo27/go-mpris"
	"github.com/godbus/dbus/v5"
	"github.com/gotk3/gotk3/gtk"
)

func loadIcon(name string) *gtk.Image {
	img, err := gtk.ImageNewFromIconName(name, gtk.ICON_SIZE_BUTTON)
	handleError(err)
	return img
}

var currentPlayer *mpris.Player

func createControllers() *gtk.Box {
	container, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 5)
	handleError(err)

	prev, err := gtk.ButtonNewFromIconName("media-seek-backward", gtk.ICON_SIZE_BUTTON)
	handleError(err)

	prev.Connect("clicked", func() {
		if currentPlayer == nil {
			return
		}
		currentPlayer.Previous()
	})

	next, err := gtk.ButtonNewFromIconName("media-seek-forward", gtk.ICON_SIZE_BUTTON)
	handleError(err)

	next.Connect("clicked", func() {
		if currentPlayer == nil {
			return
		}
		currentPlayer.Next()
	})

	playingIcon := loadIcon("media-playback-pause")
	pausedIcon := loadIcon("media-playback-start")

	playPause, err := gtk.ButtonNew()
	handleError(err)

	playPause.Connect("clicked", func() {
		if currentPlayer == nil {
			return
		}
		currentPlayer.PlayPause()
	})

	onUpdate(func(player *mpris.Player, metadata map[string]dbus.Variant) {
		currentPlayer = player
		status, err := player.GetPlaybackStatus()
		handleError(err)

		if status == mpris.PlaybackPaused {
			playPause.SetImage(pausedIcon)
		} else {
			playPause.SetImage(playingIcon)
		}
	})

	container.PackStart(prev, false, false, 0)
	container.PackStart(playPause, false, false, 0)
	container.PackStart(next, false, false, 0)
	container.SetHAlign(gtk.ALIGN_CENTER)

	return container
}
