package networkmanager

import (
	"errors"
	"fmt"
	"time"

	godbus "github.com/godbus/dbus/v5"

	"github.com/thecentinol/tuwi/internal/dbus"
)

const (
	baseIface   = "org.freedesktop.NetworkManager"
	baseObjPath = "/org/freedesktop/NetworkManager"

	// Methods
	cmGetAllDevices             = baseIface + ".GetAllDevices"
	cmActivateConnection        = baseIface + ".ActivateConnection"
	cmAddAndActivateConnection2 = baseIface + ".AddAndActivateConnection2"
	cmDeactivateConnection      = baseIface + ".DeactivateConnection"

	// Properties
	cmActiveConnections = baseIface + ".ActiveConnections"
)

var ConnectionManager connectionManager

type connectionManager struct{}

// GetAllDevices gets a slice of all network devices
func (connectionManager) GetAllDevices(c *dbus.Client) ([]godbus.ObjectPath, error) {
	obj := c.Conn.Object(baseIface, baseObjPath)
	call := obj.Call(cmGetAllDevices, 0)
	if call.Err != nil {
		return nil, fmt.Errorf("GetDevices: %w", call.Err)
	}
	var devices []godbus.ObjectPath
	if err := call.Store(&devices); err != nil {
		return nil, fmt.Errorf("GetDevices store: %w", err)
	}
	return devices, nil
}

func (connectionManager) ActivateConnection(
	c *dbus.Client,
	connectionPath godbus.ObjectPath,
	devicePath godbus.ObjectPath,
	apPath godbus.ObjectPath,
) (_ godbus.ObjectPath, err error) {
	obj := c.Conn.Object(baseIface, baseObjPath)
	call := obj.Call(
		cmActivateConnection,
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
	errAdd := c.Conn.AddMatchSignal(
		godbus.WithMatchInterface("org.freedesktop.DBus.Properties"),
		godbus.WithMatchMember("PropertiesChanged"),
		godbus.WithMatchObjectPath(activeConnection),
	)
	if errAdd != nil {
		return "", fmt.Errorf("ActivateConnection AddMatchSignal: %w", errAdd)
	}

	defer func() {
		errRm := c.Conn.RemoveMatchSignal(
			godbus.WithMatchInterface("org.freedesktop.DBus.Properties"),
			godbus.WithMatchMember("PropertiesChanged"),
			godbus.WithMatchObjectPath(activeConnection),
		)
		if errRm != nil {
			err = errors.Join(err, fmt.Errorf("ActivateConnection RemoveMatchSignal: %w", err))
		}
		c.Conn.RemoveSignal(ch)
	}()

	for {
		select {
		case sig := <-ch:

			iface := sig.Body[0]
			if iface != baseIface+".Connection.Active" {
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

// AddAndActivateConnection2 adds a new connection using the given details (if any) as a template
// (automatically filling in missing settings with the capabilities of the given device and specific
// object), then activates the new connection
func (connectionManager) AddAndActivateConnection2(
	c *dbus.Client,
	network AccessPointProperties,
	password string,
	securityType string,
) error {
	settings := baseSettings(network)
	if securityType != "open" && securityType != "unknown" {
		settings[ConnectionSettingsWirelessBlock]["security"] = godbus.MakeVariant(
			ConnectionSettingsWirelessSecurityBlock,
		)
		settings[ConnectionSettingsWirelessSecurityBlock] = securitySettings(securityType, password)
	}

	options := map[string]godbus.Variant{}

	obj := c.Conn.Object(baseIface, baseObjPath)
	call := obj.Call(cmAddAndActivateConnection2, 0, settings, network.DevicePath, network.APPath,
		options,
	)
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

// DeactivateConnection disconnects from the currently active connection
func (connectionManager) DeactivateConnection(c *dbus.Client,
	activeConnections []godbus.ObjectPath,
) error {
	obj := c.Conn.Object(baseIface, baseObjPath)

	var activeConnection godbus.ObjectPath
	for _, activePath := range activeConnections {
		conn := c.Conn.Object(baseIface, activePath)

		typeVariant, err := conn.GetProperty(activeConnectionType)
		if err != nil {
			return fmt.Errorf("DeactivateConnection: Type property: %w", err)
		}
		connType, ok := typeVariant.Value().(string)
		if !ok {
			return fmt.Errorf("DeactivateConnection: unexpected type for Type property")
		}

		stateVariant, err := conn.GetProperty(activeConnectionState)
		if err != nil {
			return fmt.Errorf("DeactivateConnection: State property: %w", err)
		}
		connState, ok := stateVariant.Value().(uint32)
		if !ok {
			return fmt.Errorf("DeactivateConnection: unexpected type for State property")
		}

		if connType == ConnectionSettingsWirelessBlock && connState == 2 {
			activeConnection = activePath
			break
		}
		continue
	}

	if activeConnection == "" {
		return fmt.Errorf("DeactivateConnection: active connection not found")
	}

	call := obj.Call(cmDeactivateConnection, 0, activeConnection)
	if call.Err != nil {
		return fmt.Errorf("DeactivateConnection: %w", call.Err)
	}
	return nil
}

// ActiveConnections gets a slice of active connection object paths
func (connectionManager) ActiveConnections(c *dbus.Client) ([]godbus.ObjectPath, error) {
	obj := c.Conn.Object(baseIface, baseObjPath)
	variant, err := obj.GetProperty(cmActiveConnections)
	if err != nil {
		return nil, fmt.Errorf("GetActiveConnections: %w", err)
	}

	activeConns, ok := variant.Value().([]godbus.ObjectPath)
	if !ok {
		return nil, fmt.Errorf("GetActiveConnections: error extracting variant info")
	}

	return activeConns, nil
}

func baseSettings(network AccessPointProperties) map[string]map[string]godbus.Variant {
	hidden := len(network.SSID) == 0
	settings := map[string]map[string]godbus.Variant{
		ConnectionSettingsConnectionBlock: {
			"id":   godbus.MakeVariant(string(network.SSID)),
			"type": godbus.MakeVariant("802-11-wireless"),
		},
		ConnectionSettingsWirelessBlock: {
			"ssid":   godbus.MakeVariant(network.SSID),
			"hidden": godbus.MakeVariant(hidden),
		},
		ConnectionSettingsIpv4Block: {
			"method": godbus.MakeVariant("auto"),
		},
		ConnectionSettingsIpv6Block: {
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
