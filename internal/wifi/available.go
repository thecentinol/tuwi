package wifi

import (
	"fmt"

	"github.com/thecentinol/tuwi/internal/dbus"
	nm "github.com/thecentinol/tuwi/internal/networkmanager"
)

// get a slice of the available wifi networks in range
func GetAvailableNetworks(c *dbus.Client) ([]nm.AccessPoint, error) {
	wifiDevicePath, err := nm.GetWifiDevice(c)
	if err != nil {
		return nil, fmt.Errorf("GetAvailableNetworks: GetWifiDevice: %w", err)
	}

	if err := nm.RequestScan(c, wifiDevicePath); err != nil {
		return nil, fmt.Errorf("GetAvailableNetworks: %w", err)
	}

	accessPoints, err := nm.GetAccessPoints(c, wifiDevicePath)
	if err != nil {
		return nil, fmt.Errorf("GetAvailableNetworks: GetAccessPoints: %w", err)
	}

	var networks []nm.AccessPoint

	for _, accessPointPath := range accessPoints {
		ap := c.Conn.Object(nm.BaseServiceName, accessPointPath)

		ssidBytes, err := ap.GetProperty(nm.ApSsid)
		if err != nil {
			return nil, fmt.Errorf("GetAvailableNetworks: AccessPointSsid property: %w", err)
		}
		rawSsid := ssidBytes.Value().([]byte)

		hidden, ssid := nm.IsHidden(rawSsid)

		bssid, err := ap.GetProperty(nm.ApHwAddress)
		if err != nil {
			return nil, fmt.Errorf("GetAvailableNetworks: AccessPointBssid property: %w", err)
		}

		strength, err := ap.GetProperty(nm.ApStrength)
		if err != nil {
			return nil, fmt.Errorf("GetAvailableNetworks: AccessPointStrength property: %w", err)
		}

		flagsRaw, err := ap.GetProperty(nm.ApFlags)
		if err != nil {
			return nil, fmt.Errorf("GetAvailableNetworks: AccessPointFlags property: %w", err)
		}
		wpaFlags, err := ap.GetProperty(nm.ApWpaFlags)
		if err != nil {
			return nil, fmt.Errorf("GetAvailableNetworks: AccessPointWpaFlags property: %w", err)
		}
		rsnFlags, err := ap.GetProperty(nm.ApRsnFlags)
		if err != nil {
			return nil, fmt.Errorf("GetAvailableNetworks: AccessPointRsnFlags property: %w", err)
		}

		flags := flagsRaw.Value().(uint32)
		wpa := wpaFlags.Value().(uint32)
		rsn := rsnFlags.Value().(uint32)
		securityType := nm.DetermineSecurityType(flags, wpa, rsn)

		networks = append(networks, nm.AccessPoint{
			SSID:         ssid,
			BSSID:        bssid.Value().(string),
			SecurityType: securityType,

			DevicePath: wifiDevicePath,
			APPath:     accessPointPath,

			Strength: strength.Value().(uint8),

			IsSaved: false,
			Secured: securityType != "open",
			Hidden:  hidden,
			HasWps:  flags&0x00000002 != 0,
		})
	}

	return networks, nil
}
