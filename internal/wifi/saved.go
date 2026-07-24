package wifi

import (
	"fmt"
	godbus "github.com/godbus/dbus/v5"

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

// uses the ObjectPaths of saved connection profiles that are available
// to connect to and builds a SavedConnection struct for each one.
func GetNearbySavedNetworks(c *dbus.Client) ([]nm.SavedConnection, error) {
	available, err := nm.AvailableConnections(c)
	if err != nil {
		return nil, fmt.Errorf("GetNearbySavedNetworks: %w", err)
	}

	var nearby []nm.SavedConnection

	for _, a := range available {
		saved, err := nm.GetSettings(c, a)
		if err != nil {
			return nil, fmt.Errorf("GetNearbySavedNetworks: %w", err)
		}

		nearby = append(nearby, *saved)
	}
	return nearby, nil
}

// when a NewConnectionMsg is received, this function is used to
// build the SavedConnection struct from the returned connection-
// path from NewConnectionMsg
func BuildNewConnection(c *dbus.Client, connectionPath godbus.ObjectPath) (*nm.SavedConnection, error) {
	saved, err := nm.GetSettings(c, connectionPath)
	if err != nil {
		return nil, fmt.Errorf("AddNewConnection: %w", err)
	}

	return saved, nil
}

// get the wifi network that we're currently connected to
func GetActiveNetwork(c *dbus.Client) (*nm.ActiveConnection, error) {
	activeConnections, err := nm.GetActiveConnections(c)
	if err != nil {
		return nil, fmt.Errorf("GetActiveNetwork: %w", err)
	}

	for _, ac := range activeConnections {
		connection, err := nm.GetActiveConnectionProperties(c, ac)
		if err != nil {
			return nil, fmt.Errorf("GetActiveNetwork: %w", err)
		}

		if connection.Type == "802-11-wireless" {
			return connection, nil
		}
	}

	return nil, nil
}

// if a saved connection is in range we will match the access point
// to the saved connection profile by BSSID and return the NearbyConnection
// slice so we can display the properties of the network that are only
// available to access point paths. If a saved connection is not in range we will
// default to show info of the network available via its saved profile settings.
func DisplaySavedConnections(
	saved []nm.SavedConnection,
	available []nm.AccessPoint,
) []nm.NearbyConnection {
	nearbyConnections := make([]nm.NearbyConnection, 0, len(saved))

	visibleAPsByBSSID := make(map[string]*nm.AccessPoint)
	visibleAPsByConnection := make(map[string]*nm.AccessPoint)
	for i := range available {
		ap := &available[i]
		if ap.BSSID != "" {
			visibleAPsByBSSID[ap.BSSID] = ap
		}

		if len(ap.SSID) != 0 {
			visibleAPsByConnection[string(ap.SSID)] = ap
		}
	}

	for _, sc := range saved {
		nearbyConnection := nm.NearbyConnection{
			Connection: sc,
			AP:         nil,
		}

		if sc.BSSIDs != nil {
			for _, seenBssid := range sc.BSSIDs {
				if ap, found := visibleAPsByBSSID[seenBssid]; found {
					nearbyConnection.AP = ap
					break
				}
			}

		}
		if sc.SSID != "" {
			if ap, found := visibleAPsByConnection[sc.SSID]; found {
				nearbyConnection.AP = ap
			}
		}
		nearbyConnections = append(nearbyConnections, nearbyConnection)
	}

	return nearbyConnections
}
