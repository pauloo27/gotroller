package window

import (
	"github.com/Pauloo27/go-mpris"
	"github.com/godbus/dbus/v5"
)

type UpdateFunction func(player *mpris.Player, metadata map[string]dbus.Variant)

var updaters []UpdateFunction

func updateAll(player *mpris.Player) {
	metadata, err := player.GetMetadata()
	handleError(err)
	for _, update := range updaters {
		update(player, metadata)
	}
}

func onUpdate(update UpdateFunction) {
	updaters = append(updaters, update)
}
