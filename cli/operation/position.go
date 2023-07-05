package operation

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/pauloo27/gotroller"
)

func SetPosition() {
	if len(os.Args) < 3 {
		fmt.Println("Missing position (valid values (in seconds): 1.5 +1.5 -1.5)")
		os.Exit(-1)
	}

	raw := os.Args[2]
	operation := SET
	if strings.HasPrefix(raw, "+") {
		raw = strings.TrimPrefix(raw, "+")
		operation = INCREASE
	} else if strings.HasPrefix(raw, "-") {
		raw = strings.TrimPrefix(raw, "-")
		operation = DECREASE
	}

	parsedPosition, err := strconv.ParseFloat(raw, 64)
	if err != nil {
		fmt.Println("Invalid position (valid values (in seconds): 1.5 +1.5 -1.5)")
		os.Exit(-1)
	}

	player, err := gotroller.GetBestPlayer()
	if err != nil {
		panic(err)
	}

	if player == nil {
		fmt.Printf("No players found")
		os.Exit(-1)
	}

	currentPosition, err := player.GetPosition()
	if err != nil {
		panic(err)
	}

	var newPosition float64

	switch operation {
	case SET:
		newPosition = parsedPosition
	case INCREASE:
		newPosition = currentPosition + parsedPosition
	case DECREASE:
		newPosition = currentPosition - parsedPosition
	}

	if newPosition < 0 {
		fmt.Println("Position cannot be negative")
		os.Exit(-1)
	}

	err = player.SetPosition(newPosition)
	if err != nil {
		panic(err)
	}
}
