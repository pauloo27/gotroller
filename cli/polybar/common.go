package polybar

import (
	"fmt"
	"html"
	"os"
	"strings"

	"github.com/Pauloo27/go-mpris"
	"github.com/Pauloo27/gotroller"
	"github.com/Pauloo27/gotroller/cli/utils"
)

var (
	maxTitleSize, maxArtistSize int
)

func handleError(err error, message string) {
	if err != nil {
		fmt.Println(Span{UNDERLINE, "#ff0000", message}.String())
		os.Exit(-1)
	}
}

func gotrollerCLI(command string) string {
	return fmt.Sprintf("gotroller %s", command)
}

func startMainLoop(playerSelectCommand string) {
	maxTitleSize, maxArtistSize = utils.LoadMaxSizes()
	utils.StartMainLoop(Polybar{})
}

func printToPolybar(playerSelectCommand string, player *mpris.Player) {
	metadata, err := player.GetMetadata()
	handleError(err, "Cannot get player metadata")

	status, err := player.GetPlaybackStatus()
	handleError(err, "Cannot get playback status")

	volume, _ := player.GetVolume()
	// ok so some apps (like firefox) do not return the volume...
	//handleError(err, "Cannot get volume")

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

	fullTitle := utils.EnforceSize(title, maxTitleSize)
	if artist != "" {
		fullTitle += " from " + utils.EnforceSize(artist, maxArtistSize)
	}
	// since lainon.life radios' uses HTML notation in the "japanese" chars
	// we need to decode them
	fullTitle = html.UnescapeString(fullTitle)

	playerSelectorAction := ActionButton{LEFT_CLICK, gotroller.MENU, playerSelectCommand}

	playPause := ActionButton{LEFT_CLICK, icon, gotrollerCLI("play-pause")}

	// previous + restart
	previous := ActionOver(
		ActionButton{LEFT_CLICK, gotroller.PREVIOUS, gotrollerCLI("prev")},
		RIGHT_CLICK, gotrollerCLI("position 0"), // TODO:
	)

	next := ActionButton{LEFT_CLICK, gotroller.NEXT, gotrollerCLI("next")}

	volumeAction := ActionOver(
		ActionButton{SCROLL_UP, fmt.Sprintf("%s %.f%%", gotroller.VOLUME, volume*100), gotrollerCLI("volume +0.05")},
		SCROLL_DOWN,
		gotrollerCLI("volume -0.05"),
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
