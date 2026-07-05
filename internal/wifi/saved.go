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

// get the wifi network that's currently connected to
func GetActiveNetwork(c *dbus.Client) (*nm.AccessPoint, error) {
	var activeNetwork nm.AccessPoint

	devicePath, err := nm.GetWifiDevice(c)
	if err != nil {
		return nil, fmt.Errorf("GetActiveNetwork: %w", err)
	}
	device := c.Conn.Object(nm.BaseServiceName, devicePath)
	active, err := device.GetProperty(nm.WirelessActiveAccessPoint)
	if err != nil {
		return nil, fmt.Errorf("GetActiveNetwork: ActiveAccessPoint property: %w", err)
	}
	activePath := active.Value().(godbus.ObjectPath)
	apObject := c.Conn.Object(nm.BaseServiceName, activePath)

	ssid, err := apObject.GetProperty(nm.ApSsid)
	if err != nil {
		return nil, fmt.Errorf("GetActiveNetwork: ActivePointSsid property: %w", err)
	}

	strength, err := apObject.GetProperty(nm.ApStrength)
	if err != nil {
		return nil, fmt.Errorf("GetActiveNetwork: AccessPointStrength property: %w", err)
	}

	flags, err := apObject.GetProperty(nm.ApFlags)
	if err != nil {
		return nil, fmt.Errorf("GetActiveNetwork: AccessPointFlags property: %w", err)
	}

	ssidBytes := string(ssid.Value().([]byte))
	flagsVal := flags.Value().(uint32)
	activeNetwork = nm.AccessPoint{
		SSID:     ssidBytes,
		Strength: strength.Value().(uint8),
		Secured:  flagsVal&0x00000001 != 0,
		HasWps:   flagsVal&0x00000002 != 0,
	}

	return &activeNetwork, nil
}
