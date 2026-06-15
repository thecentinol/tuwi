// This file
package wifi

import (
	"fmt"
	godbus "github.com/godbus/dbus/v5"
	"github.com/thecentinol/tuwi/internal/dbus"
)

// this gets the actual APs + details (ssid, strength, etc)
// for available WiFi networks
func GetAvailableNetworks(c *dbus.Client) ([]AccessPoint, error) {
	accessPoints, err := dbus.GetAccessPoints(c)
	if err != nil {
		return nil, fmt.Errorf("Failed to get access points: %v", err)
	}

	var networks []AccessPoint

	for _, accessPointPath := range accessPoints {
		ap := c.Conn.Object(dbus.BaseServiceName, accessPointPath)

		ssid, err := ap.GetProperty(dbus.AccessPointSsid)
		if err != nil {
			return nil, err
		}

		bssid, err := ap.GetProperty(dbus.AccessPointBssid)
		if err != nil {
			return nil, err
		}

		strength, err := ap.GetProperty(dbus.AccessPointStrength)
		if err != nil {
			return nil, err
		}

		flags, err := ap.GetProperty(dbus.AccessPointFlags)
		if err != nil {
			return nil, err
		}

		ssidBytes := string(ssid.Value().([]byte))
		flagsVal := flags.Value().(uint32)
		networks = append(networks, AccessPoint{
			Hidden:   ssidBytes == "",
			SSID:     ssidBytes,
			BSSID:    bssid.Value().(string),
			Strength: strength.Value().(uint8),
			Secured:  flagsVal&0x00000001 != 0,
			HasWps:   flagsVal&0x00000002 != 0,
		})
	}

	return networks, nil
}

func GetSavedNetworks(c *dbus.Client, available []AccessPoint) ([]AccessPoint, error) {
	settings := c.Conn.Object(dbus.BaseServiceName, dbus.SettingsBaseObjPath)
	call := settings.Call(dbus.ListConnections, 0)
	if call.Err != nil {
		return nil, call.Err
	}

	var connPaths []godbus.ObjectPath
	if err := call.Store(&connPaths); err != nil {
		return nil, err
	}

	var savedNetworks []AccessPoint

	for _, path := range connPaths {
		conn := c.Conn.Object(dbus.BaseServiceName, path)
		pathCall := conn.Call(dbus.GetSettings, 0)

		// this is the response shape for calling GetSettings,
		// as seen in ...Settings.Connection in the docs
		var settings map[string]map[string]godbus.Variant
		if err := pathCall.Store(&settings); err != nil {
			return nil, err
		}

		// filter out non-WiFi connections
		if settings["connection"]["type"].Value().(string) != "802-11-wireless" {
			continue
		}

		// now we can get the SSID & other details
		// NOTE: strength is not available here
		ssid := string(settings["802-11-wireless"]["ssid"].Value().([]byte))
		_, secured := settings["802-11-wireless-security"] // checking if key exists
		// get the full slice of bssids to compare the
		bssids := settings["802-11-wireless"]["seen-bssids"].Value().([]string)

		var strength uint8
		for _, network := range available {
			for _, bssid := range bssids {
				if network.BSSID == bssid {
					strength = network.Strength
				}
			}
		}

		savedNetworks = append(savedNetworks, AccessPoint{
			Hidden:   ssid == "",
			SSID:     ssid,
			BSSID:    bssids[0],
			Strength: strength,
			Secured:  secured,
		})
	}

	return savedNetworks, nil
}

func GetActiveNetwork(c *dbus.Client) (*AccessPoint, error) {
	var activeNetwork AccessPoint

	devicePath, err := dbus.GetWifiDevice(c)
	if err != nil {
		return nil, err
	}
	device := c.Conn.Object(dbus.BaseServiceName, devicePath)
	active, err := device.GetProperty(dbus.ActiveAccessPoint)
	if err != nil {
		return nil, err
	}
	activePath := active.Value().(godbus.ObjectPath)
	apObject := c.Conn.Object(dbus.BaseServiceName, activePath)

	ssid, err := apObject.GetProperty(dbus.AccessPointSsid)
	if err != nil {
		return nil, err
	}

	strength, err := apObject.GetProperty(dbus.AccessPointStrength)
	if err != nil {
		return nil, err
	}

	flags, err := apObject.GetProperty(dbus.AccessPointFlags)
	if err != nil {
		return nil, err
	}

	ssidBytes := string(ssid.Value().([]byte))
	flagsVal := flags.Value().(uint32)
	activeNetwork = AccessPoint{
		SSID:     ssidBytes,
		Strength: strength.Value().(uint8),
		Secured:  flagsVal&0x00000001 != 0,
		HasWps:   flagsVal&0x00000002 != 0,
	}

	return &activeNetwork, nil
}
