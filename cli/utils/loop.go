package utils

import (
	"errors"
	"os"
	"sync"
	"time"

	"github.com/Pauloo27/go-mpris"
	"github.com/Pauloo27/gotroller"
	"github.com/fsnotify/fsnotify"
	"github.com/godbus/dbus/v5"
)

type BarAdapter interface {
	PrintDisabled()
	HandleError(error, string)
	HandleNothingPlaying() (shouldExit bool)
	Update(*mpris.Player)
}

func StartMainLoop(bar BarAdapter) {
	player := mustGetPlayer(bar)

	lastUpdateRequest := 0
	updateRequestLock := sync.Mutex{}

	scheduleUpdate := func() {
		updateRequestLock.Lock()
		defer updateRequestLock.Unlock()
		lastUpdateRequest++
		curRequest := lastUpdateRequest
		go func() {
			time.Sleep(100 * time.Millisecond)
			if lastUpdateRequest != curRequest {
				return
			}
			bar.Update(player)
		}()
	}

	if player == nil {
		shouldExit := bar.HandleNothingPlaying()
		if shouldExit {
			os.Exit(0)
		}
		for {
			time.Sleep(time.Second)
			player = mustGetPlayer(bar)
			if player != nil {
				break
			}
		}
	}

	bar.Update(player)
	mprisCh := make(chan *dbus.Signal)
	err := player.OnSignal(mprisCh)
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
		scheduleUpdate()
	}
}

func mustGetPlayer(bar BarAdapter) *mpris.Player {
	player, err := gotroller.GetBestPlayer()
	if err != nil {
		if errors.Is(err, gotroller.ErrDisabled{}) {
			bar.PrintDisabled()
			os.Exit(0)
		}
		bar.HandleError(err, "Cannot get best player")
	}
	return player
}
