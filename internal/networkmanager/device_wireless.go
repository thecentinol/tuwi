package networkmanager

import (
	"fmt"
	"time"

	godbus "github.com/godbus/dbus/v5"
	"github.com/thecentinol/tuwi/internal/dbus"
)

func RequestScan(c *dbus.Client, devicePath godbus.ObjectPath) error {
	device := c.Conn.Object(
		BaseServiceName,
		devicePath,
	)
	call := device.Call(WirelessRequestScan, 0, map[string]godbus.Variant{})
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
		BaseServiceName,
		devicePath,
	)

	call := obj.Call(WirelessGetAllAccessPoints, 0)
	if call.Err != nil {
		return nil, fmt.Errorf("GetAccessPoints: %w", call.Err)
	}

	var accessPoints []godbus.ObjectPath
	if err := call.Store(&accessPoints); err != nil {
		return nil, fmt.Errorf("GetAccessPoints store: %w", err)
	}
	return accessPoints, nil
}

func ActivateConnection(
	c *dbus.Client,
	connectionPath godbus.ObjectPath,
	devicePath godbus.ObjectPath,
	apPath godbus.ObjectPath,
) (godbus.ObjectPath, error) {
	obj := c.Conn.Object(BaseServiceName, BaseObjPath)
	call := obj.Call(
		CmActivateConnection,
		0,
		connectionPath,
		devicePath,
		apPath,
	)
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

			iface := sig.Body[0]
			if iface != BaseServiceName+".Connection.Active" {
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
	obj := c.Conn.Object(BaseServiceName, BaseObjPath)

	var activeConnection godbus.ObjectPath
	for _, activePath := range activeConnections {
		conn := c.Conn.Object(BaseServiceName, activePath)

		typeVariant, err := conn.GetProperty(ActConType)
		if err != nil {
			return fmt.Errorf("DeactivateConnection: Type property: %w", err)
		}
		connType, ok := typeVariant.Value().(string)
		if !ok {
			return fmt.Errorf("DeactivateConnection: unexpected type for Type property")
		}

		stateVariant, err := conn.GetProperty(ActConState)
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

	call := obj.Call(CmDeactivateConnection, 0, activeConnection)
	if call.Err != nil {
		return fmt.Errorf("DeactivateConnection: %w", call.Err)
	}
	return nil
}
