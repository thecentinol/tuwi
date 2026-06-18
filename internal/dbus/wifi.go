package dbus

import (
	"fmt"
	godbus "github.com/godbus/dbus/v5"
	"github.com/thecentinol/tuwi/internal/models"
	"time"
)

func GetDevices(c *Client) ([]godbus.ObjectPath, error) {
	nm := c.Conn.Object(BaseServiceName, BaseObjPath)
	call := nm.Call(BaseServiceName+".GetDevices", 0)
	if call.Err != nil {
		return nil, call.Err
	}
	var devices []godbus.ObjectPath
	if err := call.Store(&devices); err != nil {
		return nil, err
	}
	return devices, nil
}

func GetWifiDevice(c *Client) (godbus.ObjectPath, error) {
	devices, err := GetDevices(c)
	if err != nil {
		return "", err
	}

	// finding the wifi device
	for _, v := range devices {
		device := c.Conn.Object(BaseServiceName, v)
		variant, err := device.GetProperty(DeviceType)
		if err != nil {
			return "", err
		}
		if deviceTypeVal, ok := variant.Value().(uint32); ok {
			if deviceTypeVal == 2 {
				return v, nil
			}
		}
	}
	return "", fmt.Errorf("No wifi device found")
}

func RequestScan(c *Client, devicePath godbus.ObjectPath) error {
	device := c.Conn.Object(
		BaseServiceName,
		devicePath,
	)
	call := device.Call(DeviceWireless+".RequestScan", 0, map[string]godbus.Variant{})
	if call.Err != nil {
		return call.Err
	}

	// channel for listening to LastScan property
	ch := make(chan *godbus.Signal, 1)
	c.Conn.Signal(ch) // register the channel to receive all signal msgs

	c.Conn.AddMatchSignal(
		godbus.WithMatchInterface("org.freedesktop.DBus.Properties"),
		godbus.WithMatchMember("PropertiesChanged"),
		godbus.WithMatchObjectPath(devicePath),
	)

	select {
	case <-ch:
		// scan complete, continue
	case <-time.After(10 * time.Second):
		return fmt.Errorf("scan timed out")
	}

	c.Conn.RemoveMatchSignal(
		godbus.WithMatchInterface("org.freedesktop.DBus.Properties"),
		godbus.WithMatchMember("PropertiesChanged"),
		godbus.WithMatchObjectPath(devicePath),
	)
	c.Conn.RemoveSignal(ch) // deregister the channel
	return nil
}

// this gets the raw paths for the access points
func GetAccessPoints(c *Client, devicePath godbus.ObjectPath) ([]godbus.ObjectPath, error) {
	nm := c.Conn.Object(
		BaseServiceName,
		devicePath,
	)

	call := nm.Call(DeviceWireless+".GetAllAccessPoints", 0)
	if call.Err != nil {
		return nil, call.Err
	}

	var accessPoints []godbus.ObjectPath
	if err := call.Store(&accessPoints); err != nil {
		return nil, err
	}
	return accessPoints, nil
}

func AddAndActivateConnection(
	c *Client,
	network models.AccessPoint,
	password string,
) error {
	// keyMgmt := network.SecurityType
	// if len(keyMgmt) > 0 {
	// } else {
	//     keyMgmt = "wpa-psk"
	// }
	settings := map[string]map[string]godbus.Variant{
		"connection": {
			"id":   godbus.MakeVariant(network.SSID),
			"type": godbus.MakeVariant("802-11-wireless"),
			"uuid": godbus.MakeVariant(network.ConnectionUUID),
		},
		"802-11-wireless": {
			"ssid":     godbus.MakeVariant([]byte(network.SSID)),
			"security": godbus.MakeVariant("802-11-wireless-security"),
			"hidden":   godbus.MakeVariant(network.Hidden),
		},
		"802-11-wireless-security": {
			"key-mgmt": godbus.MakeVariant(network.SecurityType),
			"psk":      godbus.MakeVariant(password),
		},
		"ipv4": {
			"method": godbus.MakeVariant("auto"),
		},
		"ipv6": {
			"method": godbus.MakeVariant("auto"),
		},
	}

	options := map[string]godbus.Variant{}

	obj := c.Conn.Object(BaseServiceName, BaseObjPath)
	call := obj.Call(AddAndConnect, 0, settings, network.DevicePath, network.APPath, options)
	if call.Err != nil {
		return call.Err
	}

	// these variables are the responses from calling
	// the AddAndActivateConnection2() method.
	var connPath godbus.ObjectPath
	var activeConnPath godbus.ObjectPath
	var result map[string]godbus.Variant
	if err := call.Store(&connPath, &activeConnPath, &result); err != nil {
		return err
	}
	return nil
}
