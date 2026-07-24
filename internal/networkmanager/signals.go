package networkmanager

import (
	"fmt"
	godbus "github.com/godbus/dbus/v5"

	"github.com/thecentinol/tuwi/internal/dbus"
)

func handleAccessPointAdded(c *dbus.Client, signal *godbus.Signal) (*AccessPoint, error) {
	body := signal.Body
	if len(body) != 1 {
		return nil, fmt.Errorf("handleAccessPointAdded: signal body missing")
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
		return "", fmt.Errorf("handleAccessPointRemoved: signal body missing")
	}

	apPath, ok := body[0].(godbus.ObjectPath)
	if !ok {
		return "", fmt.Errorf("handleAccessPointRemoved: failed to get AP path")
	}
	return apPath, nil
}

func handleNewConnection(signal *godbus.Signal) (godbus.ObjectPath, error) {
	body := signal.Body
	if len(body) != 1 {
		return "", fmt.Errorf("handleNewConnection: signal body missing")
	}

	conPath, ok := body[0].(godbus.ObjectPath)
	if !ok {
		return "", fmt.Errorf("handleNewConnection: failed to get connection path")
	}

	return conPath, nil
}

func handleConnectionRemoved(signal *godbus.Signal) (godbus.ObjectPath, error) {
	body := signal.Body
	if len(body) != 1 {
		return "", fmt.Errorf("handleConnectionRemoved: signal body missing")
	}

	conPath, ok := body[0].(godbus.ObjectPath)
	if !ok {
		return "", fmt.Errorf("handleConnectionRemoved: failed to get connection path")
	}

	return conPath, nil
}

// org.freedesktop.NetworkManager.Connection.Active
func handleActiveConnectionState(signal *godbus.Signal) (uint32, uint32, error) {
	body := signal.Body
	if len(body) != 2 {
		return 0, 0, fmt.Errorf("handleActiveConnectionState: signal body missing")
	}

	state, ok := body[0].(uint32)
	if !ok {
		return 0, 0, fmt.Errorf("handleActiveConnectionState: failed to get state")
	}

	reason, ok := body[1].(uint32)
	if !ok {
		return 0, 0, fmt.Errorf("handleActiveConnectionState: failed to get reason")
	}

	return state, reason, nil
}

type SignalCallbacks struct {
	OnAccessPointAdded   func(AccessPoint)
	OnAccessPointRemoved func(godbus.ObjectPath)

	OnNewConnection          func(godbus.ObjectPath)
	OnConnectionRemoved      func(godbus.ObjectPath)
	OnActiveConnStateChanged func(uint32, uint32)

	OnError func(error)
}

func ListenForSignals(c *dbus.Client, cb SignalCallbacks) {
	ch := make(chan *godbus.Signal, 10)
	c.Conn.Signal(ch)
	errApAdded := c.Conn.AddMatchSignal(
		godbus.WithMatchInterface(baseWirelessServiceName),
		godbus.WithMatchMember("AccessPointAdded"),
	)
	if errApAdded != nil {
		cb.OnError(fmt.Errorf("ListenForSignals: AccessPointAdded: %w", errApAdded))
	}

	errApRm := c.Conn.AddMatchSignal(
		godbus.WithMatchInterface(baseWirelessServiceName),
		godbus.WithMatchMember("AccessPointRemoved"),
	)
	if errApRm != nil {
		cb.OnError(fmt.Errorf("ListenForSignals: AccessPointRemoved: %w", errApRm))
	}

	errNewConn := c.Conn.AddMatchSignal(
		godbus.WithMatchInterface(BaseServiceName+".Settings"),
		godbus.WithMatchMember("NewConnection"),
	)
	if errNewConn != nil {
		cb.OnError(fmt.Errorf("ListenForSignals: NewConnection: %w", errNewConn))
	}

	errConnRm := c.Conn.AddMatchSignal(
		godbus.WithMatchInterface(BaseServiceName+".Settings"),
		godbus.WithMatchMember("ConnectionRemoved"),
	)
	if errConnRm != nil {
		cb.OnError(fmt.Errorf("ListenForSignals: ConnectionRemoved: %w", errConnRm))
	}

	errStateChanged := c.Conn.AddMatchSignal(
		godbus.WithMatchInterface(baseActiveConnServiceName),
		godbus.WithMatchMember("StateChanged"),
	)
	if errStateChanged != nil {
		cb.OnError(fmt.Errorf("ListenForSignals: StateChanged: %w", errStateChanged))
	}

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

		case AcStateChanged:
			state, reason, err := handleActiveConnectionState(sig)
			if err != nil {
				cb.OnError(err)
			} else {
				cb.OnActiveConnStateChanged(state, reason)
			}
		}
	}
}
