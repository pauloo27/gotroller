package polybar

import (
	"fmt"

	"github.com/pauloo27/go-mpris"
	"github.com/pauloo27/gotroller"
	"github.com/pauloo27/gotroller/cli/utils"
)

var _ utils.BarAdapter = Polybar{}

type Polybar struct {
	playerSelectCommand string
}

func (Polybar) HandleError(err error, message string) {
	handleError(err, message)
}

func (Polybar) HandleNothingPlaying() (shouldExit bool) {
	fmt.Println("Nothing playing")
	return true
}

func (p Polybar) PrintDisabled() {
	playerSelectorAction := ActionButton{LEFT_CLICK, gotroller.MENU, p.playerSelectCommand}
	fmt.Printf("%s\n", playerSelectorAction.String())
}

func (p Polybar) Update(player *mpris.Player) {
	printToPolybar(p.playerSelectCommand, player)
}
