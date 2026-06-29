package wifi

import (
	"fmt"
	godbus "github.com/godbus/dbus/v5"
	"github.com/thecentinol/tuwi/internal/dbus"
	nm "github.com/thecentinol/tuwi/internal/networkmanager"
	"time"
)

func GetDevices(c *dbus.Client) ([]godbus.ObjectPath, error) {
	obj := c.Conn.Object(nm.BaseServiceName, nm.BaseObjPath)
	call := obj.Call(nm.BaseServiceName+".GetDevices", 0)
	if call.Err != nil {
		return nil, fmt.Errorf("GetDevices: %w", call.Err)
	}
	var devices []godbus.ObjectPath
	if err := call.Store(&devices); err != nil {
		return nil, fmt.Errorf("GetDevices store: %w", err)
	}
	return devices, nil
}

func GetWifiDevice(c *dbus.Client) (godbus.ObjectPath, error) {
	devices, err := GetDevices(c)
	if err != nil {
		return "", fmt.Errorf("GetWifiDevice: %w", err)
	}

	// finding the wifi device
	for _, v := range devices {
		device := c.Conn.Object(nm.BaseServiceName, v)
		variant, err := device.GetProperty(nm.DeviceType)
		if err != nil {
			return "", fmt.Errorf("GetWifiDevice: failed to get DeviceType property %w", err)
		}
		if deviceTypeVal, ok := variant.Value().(uint32); ok {
			if deviceTypeVal == 2 {
				return v, nil
			}
		}
	}
	return "", fmt.Errorf("GetWifiDevice: no wifi device found")
}

func RequestScan(c *dbus.Client, devicePath godbus.ObjectPath) error {
	device := c.Conn.Object(
		nm.BaseServiceName,
		devicePath,
	)
	call := device.Call(nm.RequestScan, 0, map[string]godbus.Variant{})
	if call.Err != nil {
		return fmt.Errorf("RequestScan: %w", call.Err)
	}

	// channel for listening to LastScan property
	ch := make(chan *godbus.Signal, 1)
	c.Conn.Signal(ch) // register the channel to receive all signal msgs

	c.Conn.AddMatchSignal(
		godbus.WithMatchInterface("org.freedesktop.DBus.Properties"),
		godbus.WithMatchMember("PropertiesChanged"),
		godbus.WithMatchObjectPath(devicePath),
	)

	defer func() {
		c.Conn.RemoveMatchSignal(
			godbus.WithMatchInterface("org.freedesktop.DBus.Properties"),
			godbus.WithMatchMember("PropertiesChanged"),
			godbus.WithMatchObjectPath(devicePath),
		)
		c.Conn.RemoveSignal(ch) // deregister the channel
	}()

	select {
	case <-ch:
		// scan complete, continue
	case <-time.After(10 * time.Second):
		return fmt.Errorf("RequestScan: timed out waiting for scan completion")
	}

	return nil
}

// this gets the raw paths for the access points
func GetAccessPoints(c *dbus.Client, devicePath godbus.ObjectPath) ([]godbus.ObjectPath, error) {
	obj := c.Conn.Object(
		nm.BaseServiceName,
		devicePath,
	)

	call := obj.Call(nm.GetAllAccessPoints, 0)
	if call.Err != nil {
		return nil, fmt.Errorf("GetAccessPoints: %w", call.Err)
	}

	var accessPoints []godbus.ObjectPath
	if err := call.Store(&accessPoints); err != nil {
		return nil, fmt.Errorf("GetAccessPoints store: %w", err)
	}
	return accessPoints, nil
}

func baseSettings(network nm.AccessPoint) map[string]map[string]godbus.Variant {
	settings := map[string]map[string]godbus.Variant{
		"connection": {
			"id":   godbus.MakeVariant(network.SSID),
			"type": godbus.MakeVariant("802-11-wireless"),
		},
		"802-11-wireless": {
			"ssid":   godbus.MakeVariant([]byte(network.SSID)),
			"hidden": godbus.MakeVariant(network.Hidden),
		},
		"ipv4": {
			"method": godbus.MakeVariant("auto"),
		},
		"ipv6": {
			"method": godbus.MakeVariant("auto"),
		},
	}
	return settings
}

func securitySettings(
	network nm.AccessPoint,
	password string,
) map[string]godbus.Variant {
	settings := map[string]godbus.Variant{
		"key-mgmt": godbus.MakeVariant(network.SecurityType),
		"psk":      godbus.MakeVariant(password),
	}
	return settings
}

func AddAndActivateConnection(
	c *dbus.Client,
	network nm.AccessPoint,
	password string,
) error {
	settings := baseSettings(network)
	if network.Secured {
		settings["802-11-wireless"]["security"] = godbus.MakeVariant("802-11-wireless-security")
		settings["802-11-wireless-security"] = securitySettings(network, password)
	}

	options := map[string]godbus.Variant{}

	obj := c.Conn.Object(nm.BaseServiceName, nm.BaseObjPath)
	call := obj.Call(nm.AddAndActivateConnection2, 0, settings, network.DevicePath, network.APPath, options)
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

func ActiveConnections(c *dbus.Client) ([]godbus.ObjectPath, error) {
	obj := c.Conn.Object(nm.BaseServiceName, nm.BaseObjPath)
	variant, err := obj.GetProperty(nm.ActiveConnections)
	if err != nil {
		return nil, fmt.Errorf("ActiveConnections: %w", err)
	}

	activeConns, ok := variant.Value().([]godbus.ObjectPath)
	if !ok {
		return nil, fmt.Errorf("ActiveConnections: error extracting variant info")
	}

	return activeConns, nil
}

func ActivateConnection(c *dbus.Client, network nm.AccessPoint) (godbus.ObjectPath, error) {
	obj := c.Conn.Object(nm.BaseServiceName, nm.BaseObjPath)
	call := obj.Call(nm.ActivateConnection, 0, network.ConnectionPath, network.DevicePath, network.APPath)
	if call.Err != nil {
		return "", fmt.Errorf("ActivateConnection: %w", call.Err)
	}

	var activeConnection godbus.ObjectPath
	if err := call.Store(&activeConnection); err != nil {
		return "", fmt.Errorf("ActivateConnection: store: %w", err)
	}

	ch := make(chan *godbus.Signal, 1)
	c.Conn.Signal(ch)

	// Listen to NMActiveConnectionState:
	// https://networkmanager.dev/docs/libnm/latest/libnm-nm-dbus-interface.html#NMActiveConnectionState
	// for changes to the state of the connection so we can
	// act on the signal received.
	c.Conn.AddMatchSignal(
		godbus.WithMatchInterface("org.freedesktop.DBus.Properties"),
		godbus.WithMatchMember("PropertiesChanged"),
		godbus.WithMatchObjectPath(activeConnection),
	)

	defer func() {
		c.Conn.RemoveMatchSignal(
			godbus.WithMatchInterface("org.freedesktop.DBus.Properties"),
			godbus.WithMatchMember("PropertiesChanged"),
			godbus.WithMatchObjectPath(activeConnection),
		)
		c.Conn.RemoveSignal(ch)
	}()

	for {
		select {
		case sig := <-ch:

			iface := sig.Body[0].(string)
			if iface != nm.BaseServiceName+".Connection.Active" {
				continue
			}

			props := sig.Body[1].(map[string]godbus.Variant)
			stateVar, ok := props["State"]
			if !ok {
				continue
			}

			state := stateVar.Value().(uint32)

			switch state {
			case 1:
				continue
			case 2:
				return activeConnection, nil
			case 3:
				continue
			case 4:
				return "", fmt.Errorf("ActivateConnection: connection failed")
			}
		case <-time.After(10 * time.Second):
			return "", fmt.Errorf("ActivateConnection: timed out waiting for connection activation")
		}
	}
}

func DeactivateConnection(c *dbus.Client, activeConnections []godbus.ObjectPath) error {
	obj := c.Conn.Object(nm.BaseServiceName, nm.BaseObjPath)

	var activeConnection godbus.ObjectPath
	for _, activePath := range activeConnections {
		conn := c.Conn.Object(nm.BaseServiceName, activePath)

		typeVariant, err := conn.GetProperty(nm.TypeActiveConnection)
		if err != nil {
			return fmt.Errorf("DeactivateConnection: Type property: %w", err)
		}
		connType, ok := typeVariant.Value().(string)
		if !ok {
			return fmt.Errorf("DeactivateConnection: unexpected type for Type property")
		}

		stateVariant, err := conn.GetProperty(nm.StateActiveConnection)
		if err != nil {
			return fmt.Errorf("DeactivateConnection: State property: %w", err)
		}
		connState, ok := stateVariant.Value().(uint32)
		if !ok {
			return fmt.Errorf("DeactivateConnection: unexpected type for State property")
		}

		if connType == "802-11-wireless" && connState == 2 {
			activeConnection = activePath
			break
		}
		continue
	}

	if activeConnection == "" {
		return fmt.Errorf("DeactivateConnection: active connection not found")
	}

	call := obj.Call(nm.DeactivateConnection, 0, activeConnection)
	if call.Err != nil {
		return fmt.Errorf("DeactivateConnection: %w", call.Err)
	}
	return nil
}
