package window

import (
	"github.com/godbus/dbus/v5"
)

type UpdateFunction func(metadata map[string]dbus.Variant)

var updaters []UpdateFunction

func updateAll() {
	metadata, err := playerInstance.GetMetadata()
	handleError(err)
	for _, update := range updaters {
		update(metadata)
	}
}

func onUpdate(update UpdateFunction) {
	updaters = append(updaters, update)
}
