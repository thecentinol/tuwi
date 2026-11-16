package service

import (
	"fmt"

	"github.com/thecentinol/tuwi/internal/dbus"
	nm "github.com/thecentinol/tuwi/internal/networkmanager"
)

// get a slice of the available wifi networks in range
func GetAvailableNetworks(c *dbus.Client) ([]nm.AccessPointProperties, error) {
	wifiDevicePath, err := nm.Device.GetWifiDevice(c)
	if err != nil {
		return nil, fmt.Errorf("GetAvailableNetworks: GetWifiDevice: %w", err)
	}

	if err := nm.DeviceWireless.RequestScan(c, wifiDevicePath); err != nil {
		return nil, fmt.Errorf("GetAvailableNetworks: %w", err)
	}

	accessPoints, err := nm.DeviceWireless.GetAllAccessPoints(c, wifiDevicePath)
	if err != nil {
		return nil, fmt.Errorf("GetAvailableNetworks: GetAccessPoints: %w", err)
	}

	var networks []nm.AccessPointProperties

	for _, accessPointPath := range accessPoints {
		ap, err := nm.DeviceWireless.GetAccessPointProperties(c, wifiDevicePath, accessPointPath)
		if err != nil {
			return nil, fmt.Errorf("GetAvailableNetworks: %w", err)
		}

		networks = append(networks, *ap)
	}

	return networks, nil
}

// used to filter out any access points that match a saved connection
// profile that's available to connect to and only return a slice of
// access points that don't match any saved nearby connection profile.
func DisplayAvailableAPs(
	nearbySaved []nm.ConnectionSettings, available []nm.AccessPointProperties,
) []nm.AccessPointProperties {
	var accessPoints []nm.AccessPointProperties

	savedBSSIDs := make(map[string]struct{})
	savedSSIDs := make(map[string]struct{})
	for _, ns := range nearbySaved {
		for _, bssid := range ns.BSSIDs {
			savedBSSIDs[bssid] = struct{}{}
		}

		// fallback to SSID if no BSSIDs are known
		if len(ns.BSSIDs) == 0 {
			if ns.SSID != "" {
				savedSSIDs[ns.SSID] = struct{}{}
			}
		}
	}

	for _, ap := range available {
		if _, ok := savedBSSIDs[ap.BSSID]; ok {
			continue
		}

		if _, ok := savedSSIDs[string(ap.SSID)]; ok {
			continue
		}
		accessPoints = append(accessPoints, ap)
	}
	return accessPoints
}
