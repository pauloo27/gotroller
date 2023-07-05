package window

import (
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/pauloo27/go-mpris"
	"github.com/pauloo27/gotroller"
	"github.com/pauloo27/gotroller/cli/utils"
	"github.com/godbus/dbus/v5"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/joho/godotenv"
)

const (
	HEIGHT = 150
	WIDTH  = HEIGHT + 100
)

var (
	playerInstance              *mpris.Player
	maxTitleSize, maxArtistSize int
)

func loadMaxSizes() {
	home, err := os.UserHomeDir()
	if err == nil {
		godotenv.Load(path.Join(home, ".config", "gotroller.env"))
	}
	maxTitleSize = utils.AtoiOrDefault(os.Getenv("GOTROLLER_GUI_MAX_TITLE_SIZE"), 30)
	maxArtistSize = utils.AtoiOrDefault(os.Getenv("GOTROLLER_GUI_MAX_ARTIST_SIZE"), 20)
}

func StartGUI() {
	gtk.Init(nil)
	win, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	handleError(err)

	loadMaxSizes()

	win.SetTitle("Gotroller")

	cssProvider, err := gtk.CssProviderNew()
	handleError(err)
	cssProvider.LoadFromData(cssData)
	screen, err := gdk.ScreenGetDefault()
	handleError(err)
	gtk.AddProviderForScreen(screen, cssProvider, gtk.STYLE_PROVIDER_PRIORITY_APPLICATION)

	_, err = win.Connect("destroy", func() {
		fmt.Println("Closed")
		gtk.MainQuit()
	})
	handleError(err)

	mainContainer, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 1)
	handleError(err)

	infoContainer, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 5)
	handleError(err)

	mainContainer.PackStart(createAlbumArt(), false, true, 0)
	mainContainer.PackStart(infoContainer, true, true, 10)

	infoContainer.PackStart(createTitle(), false, false, 0)
	infoContainer.PackStart(createArtist(), false, false, 0)
	infoContainer.PackStart(createProgressBar(), false, false, 0)
	infoContainer.PackStart(createControllers(), false, false, 0)

	win.Add(mainContainer)
	win.SetResizable(false)

	go func() {
		var err error
		playerInstance, err = gotroller.GetBestPlayer()
		if err != nil {
			if errors.Is(err, gotroller.ErrDisabled{}) {
				gotroller.RemovePreferedPlayerName()
				os.Exit(0)
			}
			handleError(err)
		}

		if playerInstance == nil {
			fmt.Println("No player found")
			os.Exit(-1)
		}

		ch := make(chan *dbus.Signal)
		err = playerInstance.OnSignal(ch)
		handleError(err)

		callUpdate := func() { glib.IdleAdd(func() { updateAll() }) }

		callUpdate()

		for range ch {
			callUpdate()
		}
	}()

	win.SetPosition(gtk.WIN_POS_MOUSE)
	win.SetDefaultSize(WIDTH, HEIGHT)
	win.SetRealized(false)
	win.ShowAll()
	gtk.Main()
}
