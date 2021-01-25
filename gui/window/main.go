package window

import (
	"fmt"

	"github.com/Pauloo27/gotroller"
	"github.com/godbus/dbus/v5"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

const (
	HEIGHT = 150
	WIDTH  = HEIGHT + 450
)

func StartGUI() {
	gtk.Init(nil)
	win, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	handleError(err)

	win.SetTitle("Gotroller")
	win.SetPosition(gtk.WIN_POS_MOUSE)

	_, err = win.Connect("destroy", func() {
		fmt.Println("Closed")
		gtk.MainQuit()
	})
	handleError(err)

	mainContainer, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 1)
	handleError(err)

	controllerContainer, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 1)
	handleError(err)

	mainContainer.PackStart(createAlbumArt(), false, true, 0)
	mainContainer.PackStart(controllerContainer, true, true, 1)

	controllerContainer.PackStart(createArtist(), false, false, 0)
	controllerContainer.PackStart(createTitle(), false, false, 0)

	win.Add(mainContainer)

	go func() {
		player, err := gotroller.GetBestPlayer()
		handleError(err)

		ch := make(chan *dbus.Signal)
		err = player.OnSignal(ch)
		handleError(err)

		callUpdate := func() { glib.IdleAdd(func() { updateAll(player) }) }

		callUpdate()

		for range ch {
			callUpdate()
		}
	}()

	win.SetDefaultSize(WIDTH, HEIGHT)
	win.ShowAll()
	gtk.Main()
}
