package window

import (
	"fmt"
	"time"

	"github.com/godbus/dbus/v5"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

func createProgressBar() *gtk.Scale {
	scale, err := gtk.ScaleNewWithRange(gtk.ORIENTATION_HORIZONTAL, 0, 1.0, 0.01)
	handleError(err)

	scale.SetDrawValue(false)

	scale.SetTooltipText(fmt.Sprintf(
		"%s | %s",
		"00:00", "00:00",
	))

	go func() {
		for {
			time.Sleep(1 * time.Second)

			positionInSec, err := playerInstance.GetPosition()
			if err != nil {
				continue
			}

			position := time.Duration(positionInSec) * time.Second

			positionFormatted := fmt.Sprintf(
				"%02d:%02d",
				int(position.Minutes()),
				int(position.Seconds())%60,
			)

			durationInSec, err := playerInstance.GetLength()
			if err != nil {
				continue
			}

			duration := time.Duration(durationInSec) * time.Second
			durationFormatted := fmt.Sprintf(
				"%02d:%02d",
				int(duration.Minutes()),
				int(duration.Seconds())%60,
			)

			glib.IdleAdd(func() {
				scale.SetTooltipText(fmt.Sprintf(
					"%s | %s",
					positionFormatted, durationFormatted,
				))
			})
		}
	}()

	duration := 0.0
	expectedValue := 0.0

	updateValue := func() {
		position, err := playerInstance.GetPosition()
		if err != nil {
			return
		}

		expectedValue = position / duration
		glib.IdleAdd(func() {
			scale.SetValue(expectedValue)
		})
	}

	onUpdate(func(metadata map[string]dbus.Variant) {
		duration, err = playerInstance.GetLength()
		if err == nil {
			updateValue()
		}
	})

	scale.Connect("value-changed", func() {
		value := scale.GetValue()
		if value != expectedValue {
			playerInstance.SetPosition(value * duration)
		}
	})

	go func() {
		for {
			time.Sleep(1 * time.Second)
			updateValue()
		}
	}()

	return scale
}
