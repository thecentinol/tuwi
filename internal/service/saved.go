package service

import (
	"fmt"

	godbus "github.com/godbus/dbus/v5"

	"github.com/thecentinol/tuwi/internal/dbus"
	nm "github.com/thecentinol/tuwi/internal/networkmanager"
)

// get a slice of the saved connection profile's
func GetSavedNetworks(c *dbus.Client) ([]nm.ConnectionSettings, error) {
	connPaths, err := nm.Settings.ListConnections(c)
	if err != nil {
		return nil, fmt.Errorf("GetSavedNetworks: list connections: %w", err)
	}

	var savedNetworks []nm.ConnectionSettings

	for _, path := range connPaths {
		saved, err := nm.SettingsConnection.GetSettings(c, path)
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
func GetNearbySavedNetworks(c *dbus.Client) ([]nm.ConnectionSettings, error) {
	available, err := nm.Device.AvailableConnections(c)
	if err != nil {
		return nil, fmt.Errorf("GetNearbySavedNetworks: %w", err)
	}

	var nearby []nm.ConnectionSettings

	for _, a := range available {
		saved, err := nm.SettingsConnection.GetSettings(c, a)
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
func BuildNewConnection(c *dbus.Client, connectionPath godbus.ObjectPath) (*nm.ConnectionSettings, error) {
	saved, err := nm.SettingsConnection.GetSettings(c, connectionPath)
	if err != nil {
		return nil, fmt.Errorf("AddNewConnection: %w", err)
	}

	return saved, nil
}

// get the wifi network that we're currently connected to
func GetActiveNetwork(c *dbus.Client) (*nm.ActiveConnectionProperties, error) {
	activeConnections, err := nm.ConnectionManager.ActiveConnections(c)
	if err != nil {
		return nil, fmt.Errorf("GetActiveNetwork: %w", err)
	}

	for _, ac := range activeConnections {
		connection, err := nm.ActiveConnection.Properties(c, ac)
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
	saved []nm.ConnectionSettings,
	available []nm.AccessPointProperties,
) []nm.NearbyConnection {
	nearbyConnections := make([]nm.NearbyConnection, 0, len(saved))

	visibleAPsByBSSID := make(map[string]*nm.AccessPointProperties)
	visibleAPsByConnection := make(map[string]*nm.AccessPointProperties)
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
