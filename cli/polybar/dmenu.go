package polybar

import "github.com/Pauloo27/gotroller/cli/utils"

func WithDmenu() {
	utils.LoadMaxSizes()
	startMainLoop("gotroller dmenu-select")
}
