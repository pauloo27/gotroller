package dmenu

import (
	"io"
	"log"
	"os/exec"
	"strings"

	"github.com/pauloo27/gotroller"
)

func Select() {
	cmd := exec.Command("dmenu")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Fatalf("Cannot pipe stdin: %v", err)
	}

	names, err := gotroller.ListPlayersName()
	if err != nil {
		log.Fatalf("Cannot list players: %v", err)
	}

	shortNames := []string{"disabled", "auto"}

	for _, name := range names {
		shortNames = append(shortNames, strings.TrimPrefix(name, "org.mpris.MediaPlayer2."))
	}

	io.WriteString(stdin, strings.Join(shortNames, "\n"))
	err = stdin.Close()

	out, err := cmd.Output()
	if err != nil {
		log.Fatalf("Cannot run dmenu: %v", err)
	}

	name := strings.TrimSuffix(string(out), "\n")
	switch name {
	case "auto":
		gotroller.RemovePreferedPlayerName()
	case "disabled":
		gotroller.HideGotroller()
	default:
		gotroller.SetPreferedPlayerName("org.mpris.MediaPlayer2." + name)
	}

}
