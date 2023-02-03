package utils

import (
	"errors"
	"os"

	"github.com/Pauloo27/go-mpris"
	"github.com/Pauloo27/gotroller"
	"github.com/fsnotify/fsnotify"
	"github.com/godbus/dbus/v5"
)

type BarAdapter interface {
	PrintDisabled()
	HandleError(error, string)
	HandleNothingPlaying()
	Update(*mpris.Player)
}

func StartMainLoop(bar BarAdapter) {
	player, err := gotroller.GetBestPlayer()
	if err != nil {
		if errors.Is(err, gotroller.ErrDisabled{}) {
			bar.PrintDisabled()
			return
		}
		bar.HandleError(err, "Cannot get best player")
	}
	if player == nil {
		bar.HandleNothingPlaying()
		return
	}

	bar.Update(player)
	mprisCh := make(chan *dbus.Signal)
	err = player.OnSignal(mprisCh)
	bar.HandleError(err, "Cannot listen to mpris signals")

	preferedPlayerCh := make(chan fsnotify.Event)
	gotroller.ListenToChanges(preferedPlayerCh)

	go func() {
		for range preferedPlayerCh {
			// prefered player changed
			os.Exit(0)
		}
	}()

	for sig := range mprisCh {
		if sig.Name == "org.freedesktop.DBus.NameOwnerChanged" {
			// player exitted
			if len(sig.Body) == 3 && sig.Body[0] == "org.mpris.MediaPlayer2.mpv" {
				os.Exit(0)
			}
		}
		bar.Update(player)
	}
}
