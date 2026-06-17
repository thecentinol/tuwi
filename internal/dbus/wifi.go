package dbus

import (
	"fmt"
	"github.com/godbus/dbus/v5"
	"time"
)

func GetDevices(c *Client) ([]dbus.ObjectPath, error) {
	nm := c.Conn.Object(BaseServiceName, BaseObjPath)
	call := nm.Call(BaseServiceName+".GetDevices", 0)
	if call.Err != nil {
		return nil, call.Err
	}
	var devices []dbus.ObjectPath
	if err := call.Store(&devices); err != nil {
		return nil, err
	}
	return devices, nil
}

func GetWifiDevice(c *Client) (dbus.ObjectPath, error) {
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

func RequestScan(c *Client, devicePath dbus.ObjectPath) error {
	device := c.Conn.Object(
		BaseServiceName,
		devicePath,
	)
	call := device.Call(DeviceWireless+".RequestScan", 0, map[string]dbus.Variant{})
	if call.Err != nil {
		return call.Err
	}

	// channel for listening to LastScan property
	ch := make(chan *dbus.Signal, 1)
	c.Conn.Signal(ch) // register the channel to receive all signal msgs

	c.Conn.AddMatchSignal(
		dbus.WithMatchInterface("org.freedesktop.DBus.Properties"),
		dbus.WithMatchMember("PropertiesChanged"),
		dbus.WithMatchObjectPath(devicePath),
	)

	select {
	case <-ch:
		// scan complete, continue
	case <-time.After(10 * time.Second):
		return fmt.Errorf("scan timed out")
	}

	c.Conn.RemoveMatchSignal(
		dbus.WithMatchInterface("org.freedesktop.DBus.Properties"),
		dbus.WithMatchMember("PropertiesChanged"),
		dbus.WithMatchObjectPath(devicePath),
	)
	c.Conn.RemoveSignal(ch) // deregister the channel
	return nil
}

// this gets the raw paths for the access points
func GetAccessPoints(c *Client, devicePath dbus.ObjectPath) ([]dbus.ObjectPath, error) {
	nm := c.Conn.Object(
		BaseServiceName,
		devicePath,
	)

	call := nm.Call(DeviceWireless+".GetAllAccessPoints", 0)
	if call.Err != nil {
		return nil, call.Err
	}

	var accessPoints []dbus.ObjectPath
	if err := call.Store(&accessPoints); err != nil {
		return nil, err
	}
	return accessPoints, nil
}
