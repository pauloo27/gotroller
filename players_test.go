package gotroller

import (
	"strings"
	"testing"
)

func TestListPlayers(t *testing.T) {
	names, err := ListPlayersName()
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("Found %d player(s): %s", len(names), strings.Join(names, ", "))
}
