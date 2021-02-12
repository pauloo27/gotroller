package window

import (
	"regexp"
	"strings"

	"github.com/Pauloo27/gotroller/gui/downloader"
	"github.com/godbus/dbus/v5"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

var albumImg *gtk.Image
var httpRe = regexp.MustCompile(`^https?:\/\/`)

func setAlbumImage(path string) {
	imagePix, err := gdk.PixbufNewFromFileAtSize(path, -1, HEIGHT)
	handleError(err)

	albumImg.SetFromPixbuf(imagePix)
}

func createAlbumArt() *gtk.Image {
	var err error
	albumImg, err = gtk.ImageNew()
	handleError(err)

	onUpdate(func(metadata map[string]dbus.Variant) {
		rawArtURL, ok := metadata["mpris:artUrl"]
		if !ok {
			return
		}
		artURL := rawArtURL.Value().(string)
		if httpRe.MatchString(artURL) {
			go func() {
				downloadedPath, err := downloader.DownloadRemoteArt(artURL)
				if err != nil {
					return
				}
				glib.IdleAdd(func() { setAlbumImage(downloadedPath) })
			}()
		} else if strings.HasPrefix(artURL, "file://") {
			// TODO: fix
			setAlbumImage(
				strings.ReplaceAll(strings.ReplaceAll(strings.TrimPrefix(artURL, "file://"), "%20", " "),
					"%7C", "|",
				))
		}
	})

	return albumImg
}
