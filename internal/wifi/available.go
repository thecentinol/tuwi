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
		ap, err := nm.GetAccessPointProperties(c, wifiDevicePath, accessPointPath)
		if err != nil {
			return nil, fmt.Errorf("GetAvailableNetworks: %w", err)
		}

		networks = append(networks, *ap)
	}

	return networks, nil
}
