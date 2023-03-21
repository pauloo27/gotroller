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
	for {
		handlePlayer(bar)
	}
}

func handlePlayer(bar BarAdapter) {
	player := mustGetPlayer(bar)

	lastUpdateRequest := 0
	updateRequestLock := sync.Mutex{}

	scheduleUpdate := func() {
		updateRequestLock.Lock()
		defer updateRequestLock.Unlock()
		lastUpdateRequest++
		curRequest := lastUpdateRequest
		go func() {
			time.Sleep(300 * time.Millisecond)
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

	for {
		select {
		case sig := <-mprisCh:
			if sig.Name == "org.freedesktop.DBus.NameOwnerChanged" {
				if len(sig.Body) == 3 {
					// player (maybe) exited
					return
				}
			}
			scheduleUpdate()
		case <-preferedPlayerCh:
			return
		}
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
