package polybar

import "github.com/pauloo27/gotroller/cli/utils"

func WithGUI() {
	utils.LoadMaxSizes()
	startMainLoop("gotroller-gui")
}
