package window

import (
	"fmt"

	"github.com/gotk3/gotk3/gtk"
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

	// TODO: add components

	win.SetDefaultSize(450, 150)
	win.ShowAll()
	gtk.Main()
}
