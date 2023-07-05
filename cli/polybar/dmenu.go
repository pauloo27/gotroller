package polybar

import "github.com/pauloo27/gotroller/cli/utils"

func WithDmenu() {
	utils.LoadMaxSizes()
	startMainLoop("gotroller dmenu-select")
}
