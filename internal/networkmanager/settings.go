package networkmanager

import (
	"fmt"

	godbus "github.com/godbus/dbus/v5"

	"github.com/thecentinol/tuwi/internal/dbus"
)

const (
	settingsIface   = baseIface + ".Settings"
	settingsObjPath = baseObjPath + "/Settings"

	// Methods
	settingsListConnections = settingsIface + ".ListConnections"

	// Signals
	settingsNewConnection     = settingsIface + ".NewConnection"
	settingsConnectionRemoved = settingsIface + ".ConnectionRemoved"
)

var Settings settings

type settings struct{}

// ListConnections lists the saved network connections known to NetworkManager
func (settings) ListConnections(c *dbus.Client) ([]godbus.ObjectPath, error) {
	obj := c.Conn.Object(baseIface, settingsObjPath)

	var connPaths []godbus.ObjectPath
	err := obj.Call(settingsListConnections, 0).Store(&connPaths)
	if err != nil {
		return nil, fmt.Errorf("GetSavedNetworks: list connections: %w", err)
	}

	return connPaths, nil
}
