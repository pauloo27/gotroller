package polybar

import (
	"fmt"

	"github.com/Pauloo27/go-mpris"
	"github.com/Pauloo27/gotroller"
	"github.com/Pauloo27/gotroller/cli/utils"
)

var _ utils.BarAdapter = Polybar{}

type Polybar struct {
	playerSelectCommand string
}

func (Polybar) HandleError(err error, message string) {
	handleError(err, message)
}

func (Polybar) HandleNothingPlaying() {
	fmt.Println("Nothing playing")
}

func (p Polybar) PrintDisabled() {
	playerSelectorAction := ActionButton{LEFT_CLICK, gotroller.MENU, p.playerSelectCommand}
	fmt.Printf("%s\n", playerSelectorAction.String())
}

func (p Polybar) Update(player *mpris.Player) {
	printToPolybar(p.playerSelectCommand, player)
}
