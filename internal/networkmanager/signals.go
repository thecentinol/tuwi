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
	ap, err := GetAccessPointProperties(c, signal.Path, apPath)
	if err != nil {
		return nil, fmt.Errorf("handleAccessPointAdded: %w", err)
	}

	return ap, nil
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
