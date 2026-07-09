package networkmanager

import (
	"fmt"
	godbus "github.com/godbus/dbus/v5"

	"github.com/thecentinol/tuwi/internal/dbus"
)

func handleAccessPointAdded(c *dbus.Client, signal *godbus.Signal) (*AccessPoint, error) {
	body := signal.Body
	if len(body) != 1 {
		return nil, fmt.Errorf("handleAccessPointAdded: signal body missing!")
	}
	apPath := body[0].(godbus.ObjectPath)

	obj := c.Conn.Object(BaseServiceName, apPath)

	ssidBytes, err := obj.GetProperty(ApSsid)
	if err != nil {
		return nil, fmt.Errorf("handleAccessPointAdded AP property: SSID: %w", err)
	}
	rawSsid := ssidBytes.Value().([]byte)

	hidden, ssid := IsHidden(rawSsid)

	bssid, err := obj.GetProperty(ApHwAddress)
	if err != nil {
		return nil, fmt.Errorf("handleAccessPointAdded AP property: BSSID/HwAddress: %w", err)
	}

	strength, err := obj.GetProperty(ApStrength)
	if err != nil {
		return nil, fmt.Errorf("handleAccessPointAdded AP property: strength: %w", err)
	}

	rawFlags, err := obj.GetProperty(ApFlags)
	if err != nil {
		return nil, fmt.Errorf("handleAccessPointAdded AP property: flags: %w", err)
	}
	rawWpaFlags, err := obj.GetProperty(ApWpaFlags)
	if err != nil {
		return nil, fmt.Errorf("handleAccessPointAdded AP property: wpa flags: %w", err)
	}
	rawRsnFlags, err := obj.GetProperty(ApRsnFlags)
	if err != nil {
		return nil, fmt.Errorf("handleAccessPointAdded AP property: rsn flags: %w", err)
	}

	flags := rawFlags.Value().(uint32)
	wpaFlags := rawWpaFlags.Value().(uint32)
	rsnFlags := rawRsnFlags.Value().(uint32)
	securityType := DetermineSecurityType(flags, wpaFlags, rsnFlags)

	ap := AccessPoint{
		SSID:         ssid,
		BSSID:        bssid.Value().(string),
		SecurityType: securityType,

		DevicePath: signal.Path,
		APPath:     apPath,

		Strength: strength.Value().(uint8),

		IsSaved: false,
		Secured: securityType != "open",
		Hidden:  hidden,
		HasWps:  flags&0x00000002 != 0,
	}
	return &ap, nil
}

func handleAccessPointRemoved(signal *godbus.Signal) (godbus.ObjectPath, error) {
	body := signal.Body
	if len(body) != 1 {
		return "", fmt.Errorf("handleAccessPointRemoved: signal body missing!")
	}

	apPath, ok := body[0].(godbus.ObjectPath)
	if !ok {
		return "", fmt.Errorf("handleAccessPointRemoved: failed to get AP path!")
	}
	return apPath, nil
}

func handleNewConnection(signal *godbus.Signal) (godbus.ObjectPath, error) {
	body := signal.Body
	if len(body) != 1 {
		return "", fmt.Errorf("handleNewConnection: signal body missing!")
	}

	conPath, ok := body[0].(godbus.ObjectPath)
	if !ok {
		return "", fmt.Errorf("handleNewConnection: failed to get connection path!")
	}

	return conPath, nil
}

func handleConnectionRemoved(signal *godbus.Signal) (godbus.ObjectPath, error) {
	body := signal.Body
	if len(body) != 1 {
		return "", fmt.Errorf("handleConnectionRemoved: signal body missing!")
	}

	conPath, ok := body[0].(godbus.ObjectPath)
	if !ok {
		return "", fmt.Errorf("handleConnectionRemoved: failed to get connection path!")
	}

	return conPath, nil
}

type SignalCallbacks struct {
	OnAccessPointAdded   func(ap AccessPoint)
	OnAccessPointRemoved func(apPath godbus.ObjectPath)

	OnNewConnection     func(path godbus.ObjectPath)
	OnConnectionRemoved func(path godbus.ObjectPath)

	OnError func(err error)
}

func ListenForSignals(c *dbus.Client, cb SignalCallbacks) {
	ch := make(chan *godbus.Signal, 10)
	c.Conn.Signal(ch)
	c.Conn.AddMatchSignal(
		godbus.WithMatchInterface(baseWirelessServiceName),
		godbus.WithMatchMember("AccessPointAdded"),
	)
	c.Conn.AddMatchSignal(
		godbus.WithMatchInterface(baseWirelessServiceName),
		godbus.WithMatchMember("AccessPointRemoved"),
	)
	c.Conn.AddMatchSignal(
		godbus.WithMatchInterface(BaseServiceName+".Settings"),
		godbus.WithMatchMember("NewConnection"),
	)
	c.Conn.AddMatchSignal(
		godbus.WithMatchInterface(BaseServiceName+".Settings"),
		godbus.WithMatchMember("ConnectionRemoved"),
	)

	for sig := range ch {
		switch sig.Name {
		case WirelessAccessPointAdded:
			ap, err := handleAccessPointAdded(c, sig)
			if err != nil {
				cb.OnError(err)
			} else {
				cb.OnAccessPointAdded(*ap)
			}

		case WirelessAccessPointRemoved:
			apPath, err := handleAccessPointRemoved(sig)
			if err != nil {
				cb.OnError(err)
			} else {
				cb.OnAccessPointRemoved(apPath)
			}

		case CspmNewConnection:
			conPath, err := handleNewConnection(sig)
			if err != nil {
				cb.OnError(err)
			} else {
				cb.OnNewConnection(conPath)
			}

		case CspmConnectionRemoved:
			conPath, err := handleConnectionRemoved(sig)
			if err != nil {
				cb.OnError(err)
			} else {
				cb.OnConnectionRemoved(conPath)
			}
		}
	}
}
