package volume

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/Pauloo27/gotroller"
)

type VolumeOperation int

const (
	SET      = VolumeOperation(0)
	INCREASE = VolumeOperation(1)
	DECREASE = VolumeOperation(-1)

	MAX_VOLUME = 1.5
)

func SetVolume() {
	if len(os.Args) < 3 {
		fmt.Println("Missing volume (valid values: 0.5 +0.5 -0.5)")
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

	parsedVolume, err := strconv.ParseFloat(raw, 64)
	if err != nil {
		fmt.Println("Invalid volume (valid values: 0.5 +0.5 -0.5)")
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

	currentVolume, err := player.GetVolume()
	if err != nil {
		panic(err)
	}

	var newVolume float64

	switch operation {
	case SET:
		newVolume = parsedVolume
	case INCREASE:
		newVolume = currentVolume + parsedVolume
	case DECREASE:
		newVolume = currentVolume - parsedVolume
	}

	if newVolume < 0 {
		fmt.Println("Volume cannot be negative")
		os.Exit(-1)
	}

	if newVolume > MAX_VOLUME {
		fmt.Printf("Volume cannot be more than %.2f\n", MAX_VOLUME)
		os.Exit(-1)
	}

	err = player.SetVolume(newVolume)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Volume set to %.2f", newVolume)
}
