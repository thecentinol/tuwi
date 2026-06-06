package dbus

import (
	"fmt"
	"github.com/godbus/dbus/v5"
	"log"
	"time"
)

type NetworkManager struct {
	client *Client
}

func GetDevices(c *Client) ([]dbus.ObjectPath, error) {
	nm := c.conn.Object(baseServiceName, baseObjPath)
	call := nm.Call(baseServiceName+".GetDevices", 0)
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
		log.Fatalf("Error fetching devices: %v", err)
	}

	// finding the wifi device
	for _, v := range devices {
		device := c.conn.Object(baseServiceName, v)
		variant, err := device.GetProperty(deviceType)
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

func RequestScan(c *Client) error {
	devicePath, err := GetWifiDevice(c)
	if err != nil {
		return err
	}

	device := c.conn.Object(
		baseServiceName,
		devicePath,
	)
	call := device.Call(deviceWireless+".RequestScan", 0, map[string]dbus.Variant{})
	if call.Err != nil {
		return call.Err
	}

	// channel for listening to LastScan property
	ch := make(chan *dbus.Signal, 1)
	c.conn.Signal(ch) // register the channel to receive all signal msgs

	c.conn.AddMatchSignal(
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
	return nil
}

func GetAccessPoints(c *Client) ([]dbus.ObjectPath, error) {
	devicePath, err := GetWifiDevice(c)
	if err != nil {
		return nil, err
	}

	if err := RequestScan(c); err != nil {
		return nil, err
	}

	nm := c.conn.Object(
		baseServiceName,
		devicePath,
	)

	call := nm.Call(deviceWireless+".GetAllAccessPoints", 0)
	if call.Err != nil {
		return nil, call.Err
	}

	var accessPoints []dbus.ObjectPath
	if err := call.Store(&accessPoints); err != nil {
		return nil, err
	}
	return accessPoints, nil
}
