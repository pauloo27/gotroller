package polybar

import "fmt"

type FormatType string

// from https://github.com/polybar/polybar/wiki/Formatting#format-tags
const (
	FOREGROUND = FormatType("F")
	BACKGROUND = FormatType("B")
	REVERSE    = FormatType("R")
	UNDERLINE  = FormatType("u")
	OVERLINE   = FormatType("o")
	FONT       = FormatType("T")
	OFFSET     = FormatType("o")
)

type Span struct {
	Format FormatType
	Extra  string
	Text   string
}

func (s Span) String() string {
	return fmt.Sprintf("%%{%s%s}%s%%{-%s}", s.Format, s.Extra, s.Text, s.Format)
}
