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


	if err != nil {
	}

	}



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
