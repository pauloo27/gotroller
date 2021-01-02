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

func playerctl(identity, command string) string {
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

	// play/pause
	playPause := ActionButton{LEFT_CLICK, icon, playerctl(shortIdentity, "play-pause")}

	// Print everything
	fmt.Printf("%s %s",
		title,
		playPause.String(),
	)
}
