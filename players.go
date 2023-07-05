package gotroller

import (
	"fmt"

	"github.com/pauloo27/go-mpris"
	"github.com/godbus/dbus/v5"
)

var dbusConn *dbus.Conn

// GetDbusConn Creates and cache a dbus connection. The cache is used in
// subsequent calls.
func GetDbusConn() (*dbus.Conn, error) {
	if dbusConn != nil {
		return dbusConn, nil
	}
	var err error
	dbusConn, err = dbus.SessionBus()
	return dbusConn, err
}

func ListPlayersName() ([]string, error) {
	conn, err := GetDbusConn()
	if err != nil {
		return nil, fmt.Errorf("Cannot get dbus connection: %w", err)
	}

	names, err := mpris.List(conn)
	if err != nil {
		return nil, fmt.Errorf("Cannot list players: %w", err)
	}

	return names, nil
}

type ErrDisabled struct{}

func (e ErrDisabled) Error() string {
	return "Disabled"
}

// GetBestPlayer Returns the "best" player to be displayed. THIS IS NOT CACHED,
// so avoid calling it twice.
// The "best" is the one selected by the user. If not players were selected or
// the selected player isn't running then the best is the first player in the
// list that isn't "Stopped". If all players are stopped, fallback to the
// first one in the list.
func GetBestPlayer() (*mpris.Player, error) {
	names, err := ListPlayersName()
	if err != nil {
		return nil, fmt.Errorf("Cannot get best player: %w", err)
	}

	if len(names) == 0 {
		return nil, nil
	}

	conn, err := GetDbusConn()
	if err != nil {
		return nil, fmt.Errorf("Cannot get dbus connection: %w", err)
	}

	preferedPlayerName, err := GetPreferedPlayerName()

	if preferedPlayerName == "Disabled" {
		return nil, ErrDisabled{}
	}

	if preferedPlayerName != "" && err == nil {
		for _, name := range names {
			if name == preferedPlayerName {
				// N I C E: Founded! Now connect and return!
				return mpris.New(conn, name), nil
			}
		}
	}

	for _, name := range names {
		player := mpris.New(conn, name)
		if err != nil {
			return nil, fmt.Errorf("Cannot connect to player: %w", err)
		}
		status, err := player.GetPlaybackStatus()
		if err != nil && status != mpris.PlaybackStopped {
			return player, nil
		}
	}

	return mpris.New(conn, names[0]), nil
}
