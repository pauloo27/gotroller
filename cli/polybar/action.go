package polybar

import "fmt"

type MouseIndex uint

// From https://github.com/polybar/polybar/wiki/Formatting#action-a
const (
	LEFT_CLICK = MouseIndex(iota)
	MIDDLE_CLICK
	RIGHT_CLICK
	SCROLL_UP
	SCROLL_DOWN
	// Double click is kinda "meh", so avoid it
	DOUBLE_LEFT_CLICK
	DOUBLE_MIDDLE_CLICK
	DOUBLE_RIGHT_CLICK
)

type PolybarActionButton struct {
	Index            MouseIndex
	Display, Command string
}

func (a *PolybarActionButton) String() string {
	return fmt.Sprintf("%%{A%d:%s:}%s%%{A}", a.Index, a.Command, a.Display)
}
