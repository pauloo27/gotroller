package dmenu

import (
	"io"
	"log"
	"os/exec"
	"strings"

	"github.com/Pauloo27/gotroller"
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

	io.WriteString(stdin, strings.Join(names, "\n"))
	err = stdin.Close()

	out, err := cmd.Output()
	if err != nil {
		log.Fatalf("Cannot run dmenu: %v", err)
	}

	gotroller.SetPreferedPlayerName(strings.TrimSuffix(string(out), "\n"))
}
