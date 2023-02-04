package waybar

import (
	"encoding/json"
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
	printToWaybar("Silence", "")
	return false
}

func (Waybar) PrintDisabled() {
	printToWaybar("Disabled", "")
}

func (Waybar) Update(player *mpris.Player) {
	metadata, err := player.GetMetadata()
	if err != nil {
		return
	}

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
		printToWaybar("Nothing playing", "")
		return
	}

	fullTitle := utils.EnforceSize(title, maxTitleSize)
	if artist != "" {
		fullTitle += " from " + utils.EnforceSize(artist, maxArtistSize)
	}

	line := fmt.Sprintf("%s %s",
		icon,
		fullTitle,
	)

	printToWaybar(line, fmt.Sprintf("%s from %s", title, artist))
}

func handleError(err error, message string) {
	if err != nil {
		printToWaybar(message, "")
	}
}

func printToWaybar(line, tooltip string) {
	line = html.EscapeString(line)
	tooltip = html.EscapeString(tooltip)

	if line != lastLine {
		data := map[string]string{
			"text":    line,
			"tooltip": tooltip,
		}
		rawJSON, err := json.Marshal(data)
		if err != nil {
			fmt.Println("Cannot marshal data")
			os.Exit(-1)
		}
		fmt.Println(string(rawJSON))
	}
	lastLine = line
}
