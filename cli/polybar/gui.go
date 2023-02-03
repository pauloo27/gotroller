package polybar

import "github.com/Pauloo27/gotroller/cli/utils"

func WithGUI() {
	utils.LoadMaxSizes()
	startMainLoop("gotroller-gui")
}
