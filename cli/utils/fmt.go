package utils

import (
	"os"
	"path"
	"strconv"

	"github.com/joho/godotenv"
)

func EnforceSize(text string, maxLen int) string {
	if maxLen <= 0 || len(text) <= maxLen {
		return text
	}

	if maxLen-3 <= 0 {
		return ""
	}

	return text[0:maxLen-5] + "…"
}

func AtoiOrDefault(str string, defaultValue int) int {
	i, err := strconv.Atoi(str)
	if err != nil {
		return defaultValue
	}
	return i
}

func LoadMaxSizes() (maxTitleSize, maxArtistSize int) {
	home, err := os.UserHomeDir()
	if err == nil {
		godotenv.Load(path.Join(home, ".config", "gotroller.env"))
	}
	maxTitleSize = AtoiOrDefault(os.Getenv("GOTROLLER_MAX_TITLE_SIZE"), 30)
	maxArtistSize = AtoiOrDefault(os.Getenv("GOTROLLER_MAX_ARTIST_SIZE"), 20)
	return
}
