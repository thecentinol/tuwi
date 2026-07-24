package networkmanager

import (
	"fmt"

	godbus "github.com/godbus/dbus/v5"

	"github.com/thecentinol/tuwi/internal/dbus"
)

func GetDevices(c *dbus.Client) ([]godbus.ObjectPath, error) {
	obj := c.Conn.Object(BaseServiceName, BaseObjPath)
	call := obj.Call(CmGetAllDevices, 0)
	if call.Err != nil {
		return nil, fmt.Errorf("GetDevices: %w", call.Err)
	}
	var devices []godbus.ObjectPath
	if err := call.Store(&devices); err != nil {
		return nil, fmt.Errorf("GetDevices store: %w", err)
	}
	return devices, nil
}

func baseSettings(network AccessPoint) map[string]map[string]godbus.Variant {
	hidden := len(network.SSID) == 0
	settings := map[string]map[string]godbus.Variant{
		SettingsConnection: {
			"id":   godbus.MakeVariant(string(network.SSID)),
			"type": godbus.MakeVariant("802-11-wireless"),
		},
		SettingsWireless: {
			"ssid":   godbus.MakeVariant(network.SSID),
			"hidden": godbus.MakeVariant(hidden),
		},
		SettingsIpv4: {
			"method": godbus.MakeVariant("auto"),
		},
		SettingsIpv6: {
			"method": godbus.MakeVariant("auto"),
		},
	}
	return settings
}

func securitySettings(
	securityType string,
	password string,
) map[string]godbus.Variant {
	settings := map[string]godbus.Variant{
		"key-mgmt": godbus.MakeVariant(securityType),
		"psk":      godbus.MakeVariant(password),
	}
	return settings
}

func AddAndActivateConnection(
	c *dbus.Client,
	network AccessPoint,
	password string,
	securityType string,
) error {
	settings := baseSettings(network)
	if securityType != "open" && securityType != "unknown" {
		settings[SettingsWireless]["security"] = godbus.MakeVariant(SettingsSecurity)
		settings[SettingsSecurity] = securitySettings(securityType, password)
	}

	options := map[string]godbus.Variant{}

	obj := c.Conn.Object(BaseServiceName, BaseObjPath)
	call := obj.Call(CmAddAndActivateConnection2, 0, settings, network.DevicePath, network.APPath, options)
	if call.Err != nil {
		return fmt.Errorf("AddAndActivateConnection: %w", call.Err)
	}

	// these variables are the responses from calling
	// the AddAndActivateConnection2() method.
	var connPath godbus.ObjectPath
	var activeConnPath godbus.ObjectPath
	var result map[string]godbus.Variant
	if err := call.Store(&connPath, &activeConnPath, &result); err != nil {
		return fmt.Errorf("AddAndActivateConnection store: %w", err)
	}
	return nil
}

func GetActiveConnections(c *dbus.Client) ([]godbus.ObjectPath, error) {
	obj := c.Conn.Object(BaseServiceName, BaseObjPath)
	variant, err := obj.GetProperty(CmActiveConnections)
	if err != nil {
		return nil, fmt.Errorf("GetActiveConnections: %w", err)
	}

	activeConns, ok := variant.Value().([]godbus.ObjectPath)
	if !ok {
		return nil, fmt.Errorf("GetActiveConnections: error extracting variant info")
	}

	return activeConns, nil
}
