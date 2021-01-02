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

	var icon string
	switch status {
	case mpris.PlaybackPaused:
		icon = gotroller.PAUSED
	case mpris.PlaybackStopped:
		icon = gotroller.STOPPED
	default:
		icon = gotroller.PLAYING
	}

	identity, err := player.GetIdentity()
	handleError(err, "Cannot get player identity")

	shortIdentity := strings.TrimPrefix(identity, "org.mpris.MediaPlayer2.")

	title := metadata["xesam:title"].Value()

	playPause := ActionButton{LEFT_CLICK, icon, plyctl(shortIdentity, "play-pause")}

	// one is "inside" of another
	previous := ActionButton{LEFT_CLICK, gotroller.PREVIOUS, plyctl(shortIdentity, "previous")}
	restart := ActionOver(previous, RIGHT_CLICK, plyctl(shortIdentity, "position 0"))

	next := ActionButton{LEFT_CLICK, gotroller.NEXT, plyctl(shortIdentity, "next")}

	// Print everything
	fmt.Printf("%s %s %s %s",
		title,
		// restart contains previous
		restart.String(),
		playPause.String(),
		next.String(),
	)
}
