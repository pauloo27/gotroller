package window

import (
	"net/url"
	"regexp"
	"strings"

	"github.com/godbus/dbus/v5"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/pauloo27/gotroller/gui/downloader"
)

var albumImg *gtk.Image
var httpRe = regexp.MustCompile(`^https?:\/\/`)

func setAlbumImage(path string) {
	imagePix, err := gdk.PixbufNewFromFileAtSize(path, -1, WIDTH)
	handleError(err)

	albumImg.SetFromPixbuf(imagePix)
}

func createAlbumArt() *gtk.Image {
	var err error
	albumImg, err = gtk.ImageNew()
	handleError(err)

	albumImg.SetSizeRequest(WIDTH, WIDTH)

	onUpdate(func(metadata map[string]dbus.Variant) {
		rawAlbumName, ok := metadata["xesam:album"]
		if ok {
			albumImg.SetTooltipText(rawAlbumName.Value().(string))
		}

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
			decodedURL, err := url.QueryUnescape(artURL)
			if err == nil {
				artURL = decodedURL
			}
			setAlbumImage(strings.TrimPrefix(artURL, "file://"))
		}
	})

	return albumImg
}
