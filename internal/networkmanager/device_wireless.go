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

// gets the ObjectPath of the access point currently used
// by the wireless device
func GetActiveAccessPoint(
	c *dbus.Client,
	wirelessDevicePath godbus.ObjectPath,
) (*godbus.ObjectPath, error) {
	obj := c.Conn.Object(BaseServiceName, wirelessDevicePath)
	rawActiveAP, err := obj.GetProperty(WirelessActiveAccessPoint)

	if err != nil {
		return nil, fmt.Errorf("GetActiveAccessPoint: %w", err)
	}

	activeAP := rawActiveAP.Value().(godbus.ObjectPath)

	return &activeAP, nil
}

// gets the properties of an access point and builds an AccessPoint struct/object
func GetAccessPointProperties(
	c *dbus.Client,
	wifiDevicePath godbus.ObjectPath,
	apPath godbus.ObjectPath,
) (*AccessPoint, error) {
	obj := c.Conn.Object(BaseServiceName, apPath)

	flags, err := obj.GetProperty(ApFlags)
	if err != nil {
		return nil, fmt.Errorf("GetAccessPointProperties: Flags property: %w", err)
	}
	wpaFlags, err := obj.GetProperty(ApWpaFlags)
	if err != nil {
		return nil, fmt.Errorf("GetAccessPointProperties: WpaFlags property: %w", err)
	}
	rsnFlags, err := obj.GetProperty(ApRsnFlags)
	if err != nil {
		return nil, fmt.Errorf("GetAccessPointProperties: RsnFlags property: %w", err)
	}

	ssid, err := obj.GetProperty(ApSsid)
	if err != nil {
		return nil, fmt.Errorf("GetAccessPointProperties: SSID property: %w", err)
	}

	freq, err := obj.GetProperty(ApFrequency)
	if err != nil {
		return nil, fmt.Errorf("GetAccessPointProperties: Frequency property: %w", err)
	}

	bssid, err := obj.GetProperty(ApHwAddress)
	if err != nil {
		return nil, fmt.Errorf("GetAccessPointProperties: HwAddress property: %w", err)
	}

	mode, err := obj.GetProperty(ApMode)
	if err != nil {
		return nil, fmt.Errorf("GetAccessPointProperties: Mode property: %w", err)
	}

	maxBitate, err := obj.GetProperty(ApMaxBitrate)
	if err != nil {
		return nil, fmt.Errorf("GetAccessPointProperties: MaxBitrate property: %w", err)
	}

	bandwidth, err := obj.GetProperty(ApBandwidth)
	if err != nil {
		return nil, fmt.Errorf("GetAccessPointProperties: Bandwidth property: %w", err)
	}

	strength, err := obj.GetProperty(ApStrength)
	if err != nil {
		return nil, fmt.Errorf("GetAccessPointProperties: Strength property: %w", err)
	}

	lastSeen, err := obj.GetProperty(ApLastSeen)
	if err != nil {
		return nil, fmt.Errorf("GetAccessPointProperties: LastSeen property: %w", err)
	}

	ap := &AccessPoint{
		Flags:    flags.Value().(uint32),
		WpaFlags: wpaFlags.Value().(uint32),
		RsnFlags: rsnFlags.Value().(uint32),

		SSID:      ssid.Value().([]byte),
		Frequency: freq.Value().(uint32),
		BSSID:     bssid.Value().(string),

		Mode:       mode.Value().(uint32),
		MaxBitrate: maxBitate.Value().(uint32),
		Bandwidth:  bandwidth.Value().(uint32),
		Strength:   strength.Value().(uint8),
		LastSeen:   lastSeen.Value().(int32),

		DevicePath: wifiDevicePath,
		APPath:     apPath,
	}

	return ap, nil
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
