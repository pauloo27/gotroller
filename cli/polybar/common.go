package polybar

import (
	"fmt"
	"strings"

	"github.com/Pauloo27/go-mpris"
	"github.com/Pauloo27/gotroller"
)

func handleError(err error, message string) {
	if err != nil {
		fmt.Println(Span{UNDERLINE, "#ff0000", message}.String())
	}
}

func plyctl(identity, command string) string {
	return fmt.Sprintf("playerctl %s -p %s", command, identity)
}

func printToPolybar(preferedPlayerSelectorCommand string) {
	player, err := gotroller.GetBestPlayer()
	handleError(err, "Cannot get best player")

	if player == nil {
		fmt.Println("--")
		return
	}

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

	identity, err := player.GetIdentity()
	handleError(err, "Cannot get player identity")

	shortIdentity := strings.TrimPrefix(identity, "org.mpris.MediaPlayer2.")

	var title string
	if rawTitle, ok := metadata["xesam:title"]; ok {
		title = rawTitle.Value().(string)
	}

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
		// TODO
		fmt.Printf("%s %s", "...", icon)
	} else {
		// Print everything
		fmt.Printf("%s %s %s %s %s",
			title,
			// restart contains previous
			previous.String(),
			playPause.String(),
			next.String(),
			volumeAction.String(),
		)
	}
}
