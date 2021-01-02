package polybar

import (
	"fmt"

	"github.com/Pauloo27/gotroller"
)

func handleError(err error, message string) {
	if err != nil {
		fmt.Println(Span{UNDERLINE, "#ff0000", message}.String())
	}
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

	title := metadata["xesam:title"].Value()

	// Print everything
	fmt.Printf("%s",
		title,
	)
}
