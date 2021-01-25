package window

import (
	"github.com/Pauloo27/go-mpris"
	"github.com/godbus/dbus/v5"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

var albumImg *gtk.Image

func setAlbumImage(path string) {
	imagePix, err := gdk.PixbufNewFromFileAtSize(path, -1, HEIGHT)
	handleError(err)

	albumImg.SetFromPixbuf(imagePix)
}

func createAlbumArt() *gtk.Image {
	var err error
	albumImg, err = gtk.ImageNew()
	handleError(err)

	onUpdate(func(player *mpris.Player, metadata map[string]dbus.Variant) {
		// TODO
		// setAlbumImage(path)
	})

	return albumImg
}
