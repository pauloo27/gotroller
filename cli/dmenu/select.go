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

	shortNames := []string{"auto"}

	for _, name := range names {
		shortNames = append(shortNames, strings.TrimPrefix(name, "org.mpris.MediaPlayer2."))
	}

	io.WriteString(stdin, strings.Join(shortNames, "\n"))
	err = stdin.Close()

	out, err := cmd.Output()
	if err != nil {
		log.Fatalf("Cannot run dmenu: %v", err)
	}

	gotroller.SetPreferedPlayerName("org.mpris.MediaPlayer2." + strings.TrimSuffix(string(out), "\n"))
}
