package gotroller

import (
	"os"
	"testing"
)

func assertNameIs(t *testing.T, expected string) {
	name, err := GetPreferedPlayerName()
	if err != nil {
		t.Fatal(err)
	}

	if name != expected {
		t.Fatalf("Name was expected to be `%s`, but it is `%s` instead", expected, name)
	}
}

func TestPreferedPlayer(t *testing.T) {
	var err error

	// Used to restore the player after the test
	previous, err := GetPreferedPlayerName()
	if err != nil {
		previous = ""
	}

	// Remove
	err = RemovePreferedPlayerName()
	if err != nil && !os.IsNotExist(err) {
		t.Fatal(err)
	}
	assertNameIs(t, "")

	// Set
	err = SetPreferedPlayerName("testing")
	if err != nil {
		t.Fatal(err)
	}
	assertNameIs(t, "testing")

	// Hide
	err = HideGotroller()
	if err != nil {
		t.Fatal(err)
	}
	assertNameIs(t, "Disabled")

	err = SetPreferedPlayerName(previous)

	if err != nil {
		t.Fatal(err)
	}
}
