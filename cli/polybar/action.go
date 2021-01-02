package polybar

import "fmt"

type MouseIndex uint

// From https://github.com/polybar/polybar/wiki/Formatting#action-a
const (
	LEFT_CLICK = MouseIndex(iota + 1)
	MIDDLE_CLICK
	RIGHT_CLICK
	SCROLL_UP
	SCROLL_DOWN
	// Double click is kinda "meh", so avoid it
	DOUBLE_LEFT_CLICK
	DOUBLE_MIDDLE_CLICK
	DOUBLE_RIGHT_CLICK
)

type ActionButton struct {
	Index            MouseIndex
	Display, Command string
}

func (a ActionButton) String() string {
	return fmt.Sprintf("%%{A%d:%s:}%s%%{A}", a.Index, a.Command, a.Display)
}

func ActionOver(a ActionButton, index MouseIndex, command string) ActionButton {
	return ActionButton{index, a.String(), command}
}
