package waybar

import (
	"fmt"
	"html"
	"strings"

	"github.com/Pauloo27/go-mpris"
	"github.com/Pauloo27/gotroller"
	"github.com/Pauloo27/gotroller/cli/utils"
)

var (
	maxTitleSize, maxArtistSize int

	lastLine string
)

func Start() {
	utils.LoadMaxSizes()
	maxTitleSize, maxArtistSize = utils.LoadMaxSizes()

	utils.StartMainLoop(Waybar{})
}

var _ utils.BarAdapter = Waybar{}

type Waybar struct{}

func (Waybar) HandleError(err error, message string) {
	handleError(err, message)
}

func (Waybar) HandleNothingPlaying() (shouldExit bool) {
	fmt.Println("Nothing playing")
	return false
}

func (Waybar) PrintDisabled() {
	fmt.Println("Disabled")
}

func (Waybar) Update(player *mpris.Player) {
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

	if title == "" && artist == "" {
		fmt.Println("Nothing playing")
		return
	}

	fullTitle := utils.EnforceSize(title, maxTitleSize)
	if artist != "" {
		fullTitle += " from " + utils.EnforceSize(artist, maxArtistSize)
	}
	// since lainon.life radios' uses HTML notation in the "japanese" chars
	// we need to decode them
	fullTitle = html.UnescapeString(fullTitle)

	line := fmt.Sprintf("%s %s",
		icon,
		fullTitle,
	)

	if line != lastLine {
		fmt.Println(line)
	}
	lastLine = line
}

func handleError(err error, message string) {
	if err != nil {
		fmt.Println(message)
	}
}
