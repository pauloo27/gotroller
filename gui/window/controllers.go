package window

import (
	"github.com/pauloo27/go-mpris"
	"github.com/godbus/dbus/v5"
	"github.com/gotk3/gotk3/gtk"
)

func loadIcon(name string) *gtk.Image {
	img, err := gtk.ImageNewFromIconName(name, gtk.ICON_SIZE_BUTTON)
	handleError(err)
	return img
}

func createControllers() *gtk.Box {
	container, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 10)
	handleError(err)

	prev, err := gtk.ButtonNewFromIconName("media-seek-backward", gtk.ICON_SIZE_BUTTON)
	handleError(err)

	prev.Connect("clicked", func() {
		if playerInstance == nil {
			return
		}
		playerInstance.Previous()
	})

	next, err := gtk.ButtonNewFromIconName("media-seek-forward", gtk.ICON_SIZE_BUTTON)
	handleError(err)

	next.Connect("clicked", func() {
		if playerInstance == nil {
			return
		}
		playerInstance.Next()
	})

	playingIcon := loadIcon("media-playback-pause")
	pausedIcon := loadIcon("media-playback-start")

	playPause, err := gtk.ButtonNew()
	handleError(err)

	playPause.Connect("clicked", func() {
		if playerInstance == nil {
			return
		}
		playerInstance.PlayPause()
	})

	onUpdate(func(metadata map[string]dbus.Variant) {
		status, err := playerInstance.GetPlaybackStatus()
		handleError(err)

		if status == mpris.PlaybackPaused {
			playPause.SetImage(pausedIcon)
		} else {
			playPause.SetImage(playingIcon)
		}
	})

	volumeBtn, err := gtk.ButtonNewFromIconName("audio-volume-high", gtk.ICON_SIZE_BUTTON)
	handleError(err)

	createVolumePopover(volumeBtn)

	container.PackStart(volumeBtn, false, false, 0)

	moreBtn, err := gtk.ButtonNewFromIconName("view-more-symbolic", gtk.ICON_SIZE_BUTTON)
	handleError(err)

	createMorePopover(moreBtn)

	container.PackEnd(moreBtn, false, false, 0)

	container.PackStart(prev, false, false, 0)
	container.PackStart(playPause, false, false, 0)
	container.PackStart(next, false, false, 0)
	container.SetHAlign(gtk.ALIGN_CENTER)

	playPause.GrabFocus()

	return container
}

func createMorePopover(relative *gtk.Button) *gtk.Popover {
	popover, err := gtk.PopoverNew(relative)
	handleError(err)

	mainContainer, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 5)
	handleError(err)

	buttonContainer, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 5)
	handleError(err)

	buttonLoop, err := gtk.ButtonNewFromIconName("media-playlist-repeat", gtk.ICON_SIZE_BUTTON)
	handleError(err)

	buttonShuffle, err := gtk.ButtonNewFromIconName("media-playlist-shuffle", gtk.ICON_SIZE_BUTTON)
	handleError(err)

	buttonRaisePlayer, err := gtk.ButtonNewFromIconName("go-up", gtk.ICON_SIZE_BUTTON)
	handleError(err)

	buttonRaisePlayer.Connect("clicked", func() {
		if playerInstance == nil {
			return
		}
		_ = playerInstance.Raise()
	})

	buttonContainer.PackStart(buttonLoop, false, false, 0)
	buttonContainer.PackStart(buttonShuffle, false, false, 0)
	buttonContainer.PackStart(buttonRaisePlayer, false, false, 0)

	mainContainer.SetMarginEnd(5)
	mainContainer.SetMarginStart(5)
	mainContainer.SetMarginTop(5)
	mainContainer.SetMarginBottom(5)

	mainContainer.PackStart(buttonContainer, false, false, 0)
	mainContainer.PackStart(createPlayerSelector(), false, false, 0)

	popover.Add(mainContainer)
	mainContainer.ShowAll()

	popover.SetPosition(gtk.POS_TOP)

	relative.Connect("clicked", func() {
		popover.SetVisible(!popover.GetVisible())
	})

	return popover
}

func createVolumePopover(relative *gtk.Button) *gtk.Popover {
	popover, err := gtk.PopoverNew(relative)
	handleError(err)

	container := createVolumeController()
	popover.Add(container)
	container.ShowAll()

	width, height := container.GetSizeRequest()
	container.SetSizeRequest(width, height+100)

	popover.SetPosition(gtk.POS_TOP)

	relative.Connect("clicked", func() {
		popover.SetVisible(!popover.GetVisible())
	})

	return popover
}
