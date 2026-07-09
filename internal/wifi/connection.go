package wifi

import (
	"fmt"
	godbus "github.com/godbus/dbus/v5"

	"github.com/thecentinol/tuwi/internal/dbus"
	nm "github.com/thecentinol/tuwi/internal/networkmanager"
)

// connect to a saved wifi network
func ConnectSaved(client *dbus.Client, nc nm.NearbyConnection) (godbus.ObjectPath, error) {
	if nc.AP == nil {
		return "", fmt.Errorf("ConnectSaved: network is out of range")
	}
	obj, err := nm.ActivateConnection(
		client,
		nc.Connection.ConnectionPath,
		nc.AP.DevicePath,
		nc.AP.APPath,
	)

	if err != nil {
		return "", fmt.Errorf("ConnectSaved: %w", err)
	}

	return obj, nil
}

// connect to an available wifi network that is not open
func ConnectSecured(
	client *dbus.Client,
	network *nm.AccessPoint,
	password string,
) error {
	err := nm.AddAndActivateConnection(
		client,
		*network,
		password,
	)

	if err != nil {
		return fmt.Errorf("ConnectSecured: %w", err)
	}

	return nil
}

// connect to an available wifi network that is open
func ConnectOpen(
	client *dbus.Client,
	network *nm.AccessPoint,
) error {
	err := nm.AddAndActivateConnection(
		client,
		*network,
		"",
	)

	if err != nil {
		return fmt.Errorf("ConnectOpen: %w", err)
	}

	return nil
}

// disconnect from a wifi network
func Disconnect(client *dbus.Client) error {
	ACs, err := nm.GetActiveConnections(client)
	if err != nil {
		return fmt.Errorf("Disconnect: %w", err)
	}
	err = nm.DeactivateConnection(client, ACs)
	if err != nil {
		return fmt.Errorf("Disconnect: %w", err)
	}
	return nil
}

// forget a wifi network - deletes the saved connection profile if using NetworkManager backend
func Forget(client *dbus.Client, conPath godbus.ObjectPath) error {
	err := nm.DeleteConnection(client, conPath)
	if err != nil {
		return fmt.Errorf("Forget: %w", err)
	}
	return nil
}
