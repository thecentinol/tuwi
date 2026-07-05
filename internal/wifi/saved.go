package wifi

import (
	"fmt"
	godbus "github.com/godbus/dbus/v5"
	"slices"

	"github.com/thecentinol/tuwi/internal/dbus"
	nm "github.com/thecentinol/tuwi/internal/networkmanager"
)

// get a slice of the saved connection profile's
func GetSavedNetworks(c *dbus.Client) ([]nm.SavedConnection, error) {
	settingsObj := c.Conn.Object(nm.BaseServiceName, nm.SettingsBaseObjPath)

	var connPaths []godbus.ObjectPath
	err := settingsObj.Call(nm.CspmListConnections, 0).Store(&connPaths)
	if err != nil {
		return nil, fmt.Errorf("GetSavedNetworks: list connections: %w", err)
	}

	var savedNetworks []nm.SavedConnection

	for _, path := range connPaths {
		saved, err := nm.GetSettings(c, path)
		if err != nil {
			return nil, fmt.Errorf("GetSavedNetworks: %w", err)
		}

		if saved != nil {
			savedNetworks = append(savedNetworks, *saved)
		}
	}

	return savedNetworks, nil
}

// when a NewConnectionMsg is received, this function is used to
// build the SavedConnection struct from the returned connection-
// path from NewConnectionMsg
func AppendNewConnection(c *dbus.Client, connectionPath godbus.ObjectPath) (*nm.NearbyConnection, error) {
	var savedConnection nm.NearbyConnection

	saved, err := nm.GetSettings(c, connectionPath)
	if err != nil {
		return nil, fmt.Errorf("AddNewConnection: %w", err)
	}

	savedConnection = nm.NearbyConnection{
		Connection: *saved,
		AP:         nil,
	}

	return &savedConnection, nil
}

// when a ConnectionRemovedMsg is received, this function is used to
// remove the selected connection from the SavedConnection (networks) slice.
func RemoveConnection(networks []nm.NearbyConnection, connectionPath godbus.ObjectPath) []nm.NearbyConnection {
	result := slices.DeleteFunc(networks, func(nc nm.NearbyConnection) bool {
		return connectionPath == nc.Connection.ConnectionPath
	})

	return result
}

// get the wifi network that we're currently connected to
func GetActiveNetwork(c *dbus.Client) (*nm.AccessPoint, error) {
	devicePath, err := nm.GetWifiDevice(c)
	if err != nil {
		return nil, fmt.Errorf("GetActiveNetwork: %w", err)
	}
	activeAP, err := nm.GetActiveAccessPoint(c, devicePath)

	ap, err := nm.GetAccessPointProperties(c, devicePath, *activeAP)
	if err != nil {
		return nil, fmt.Errorf("GetActiveNetwork: %w", err)
	}

	return ap, nil
}

	}

	}

}
