package polybar

import (
	"errors"
	"fmt"
	"html"
	"os"
	"strings"

	"github.com/Pauloo27/go-mpris"
	"github.com/Pauloo27/gotroller"
	"github.com/fsnotify/fsnotify"
	"github.com/godbus/dbus/v5"
)

func handleError(err error, message string) {
	if err != nil {
		fmt.Println(Span{UNDERLINE, "#ff0000", message}.String())
	}
}

func plyctl(identity, command string) string {
	return fmt.Sprintf("playerctl %s -p %s", command, identity)
}

func startMainLoop(playerSelectCommand string) {
	player, err := gotroller.GetBestPlayer()
	if err != nil {
		if errors.Is(err, gotroller.ErrDisabled{}) {
			playerSelectorAction := ActionButton{LEFT_CLICK, gotroller.MENU, playerSelectCommand}
			fmt.Printf("%s\n", playerSelectorAction.String())
			return
		}
		handleError(err, "Cannot get best player")
	}
	if player == nil {
		fmt.Println("Nothing playing...")
		return
	}

	update := func() {
		printToPolybar(playerSelectCommand, player)
	}

	update()
	mprisCh := make(chan *dbus.Signal)
	err = player.OnSignal(mprisCh)
	handleError(err, "Cannot listen to mpris signals")

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
		update()
	}
}

func printToPolybar(playerSelectCommand string, player *mpris.Player) {
	metadata, err := player.GetMetadata()
	handleError(err, "Cannot get player metadata")

	status, err := player.GetPlaybackStatus()
	handleError(err, "Cannot get playback status")

	volume, err := player.GetVolume()
	handleError(err, "Cannot get volume")

	stopped := false

	var icon string
	switch status {
	case mpris.PlaybackPaused:
		icon = gotroller.PAUSED
	case mpris.PlaybackStopped:
		icon = gotroller.STOPPED
		stopped = true
	default:
		icon = gotroller.PLAYING
	}

	// GetName returns identity.instancePID while GetIdentity() returns only
	// the identity.
	identity := player.GetName()
	handleError(err, "Cannot get player identity")

	shortIdentity := strings.TrimPrefix(identity, "org.mpris.MediaPlayer2.")

	var title string
	if rawTitle, ok := metadata["xesam:title"]; ok {
		title = rawTitle.Value().(string)
	}

	var artist string
	if rawArtist, ok := metadata["xesam:artist"]; ok {
		switch rawArtist.Value().(type) {
		case string:
			artist = rawArtist.Value().(string)
		case []string:
			artist = strings.Join(rawArtist.Value().([]string), ", ")
		}
	}

	fullTitle := title
	if artist != "" {
		fullTitle += " from " + artist
	}
	// since lainon.life radios' uses HTML notation in the "japanese" chars
	// we need to decode them
	fullTitle = html.UnescapeString(fullTitle)

	playerSelectorAction := ActionButton{LEFT_CLICK, gotroller.MENU, playerSelectCommand}

	playPause := ActionButton{LEFT_CLICK, icon, plyctl(shortIdentity, "play-pause")}

	// previous + restart
	previous := ActionOver(
		ActionButton{LEFT_CLICK, gotroller.PREVIOUS, plyctl(shortIdentity, "previous")},
		RIGHT_CLICK, plyctl(shortIdentity, "position 0"),
	)

	next := ActionButton{LEFT_CLICK, gotroller.NEXT, plyctl(shortIdentity, "next")}

	volumeAction := ActionOver(
		ActionButton{SCROLL_UP, fmt.Sprintf("%s %.f%%", gotroller.VOLUME, volume*100), plyctl(shortIdentity, "volume 0.05+")},
		SCROLL_DOWN,
		plyctl(shortIdentity, "volume 0.05-"),
	)

	if stopped {
		fmt.Printf("%s %s\n", playerSelectorAction.String(), icon)
	} else {
		// Print everything
		fmt.Printf("%s %s %s %s %s %s\n",
			playerSelectorAction.String(),
			fullTitle,
			// restart contains previous
			previous.String(),
			playPause.String(),
			next.String(),
			volumeAction.String(),
		)
	}
}
