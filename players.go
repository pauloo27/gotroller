package gotroller

import (
	"fmt"

	"github.com/Pauloo27/go-mpris"
	"github.com/godbus/dbus/v5"
)

var dbusConn *dbus.Conn

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
		return nil, fmt.Errorf("Cannot get dbus connection: %v", err)
	}

	names, err := mpris.List(conn)
	if err != nil {
		return nil, fmt.Errorf("Cannot list players: %v", err)
	}

	return names, nil
}
