package networkmanager

import (
	"fmt"

	godbus "github.com/godbus/dbus/v5"

	"github.com/thecentinol/tuwi/internal/dbus"
)

const (
	activeConnectionIface = baseIface + ".Connection.Active"

	// Signals
	activeConnectionStateChanged = activeConnectionIface + ".StateChanged"

	// Properties
	activeConnectionConnection     = activeConnectionIface + ".Connection"
	activeConnectionSpecificObject = activeConnectionIface + ".SpecificObject"
	activeConnectionId             = activeConnectionIface + ".Id"
	activeConnectionUuid           = activeConnectionIface + ".Uuid"
	activeConnectionType           = activeConnectionIface + ".Type"
	activeConnectionDevices        = activeConnectionIface + ".Devices"
	activeConnectionState          = activeConnectionIface + ".State"
	activeConnectionStateFlags     = activeConnectionIface + ".StateFlags"
	activeConnectionDefault        = activeConnectionIface + ".Default"
	activeConnectionIp4Config      = activeConnectionIface + ".Ip4Config"
	activeConnectionDhcp4Config    = activeConnectionIface + ".Dhcp4Config"
	activeConnectionDefault6       = activeConnectionIface + ".Default6"
	activeConnectionIp6Config      = activeConnectionIface + ".Ip6Config"
	activeConnectionDhcp6Config    = activeConnectionIface + ".Dhcp6Config"
	activeConnectionVpn            = activeConnectionIface + ".Vpn"
	activeConnectionController     = activeConnectionIface + ".Controller"
)

var ActiveConnection activeConnection

type activeConnection struct{}

type ActiveConnectionProperties struct {
	Connection     godbus.ObjectPath
	SpecificObject godbus.ObjectPath
	ID             string
	UUID           string
	Type           string
	Devices        []godbus.ObjectPath
	State          uint32
	StateFlags     uint32
	Default        bool
	IP4Config      godbus.ObjectPath
	DHCP4Config    godbus.ObjectPath
	Default6       bool
	IP6Config      godbus.ObjectPath
	DHCP6Config    godbus.ObjectPath
	VPN            bool
	Controller     godbus.ObjectPath
}

// Properties gets all properties of an access point
func (activeConnection) Properties(
	c *dbus.Client, path godbus.ObjectPath,
) (*ActiveConnectionProperties, error) {
	obj := c.Conn.Object(baseIface, path)

	conn, err := obj.GetProperty(activeConnectionConnection)
	if err != nil {
		return nil, fmt.Errorf("GetActiveConnectionProperties connection property: %w", err)
	}

	specificObj, err := obj.GetProperty(activeConnectionSpecificObject)
	if err != nil {
		return nil, fmt.Errorf("GetActiveConnectionProperties SpecificObject property: %w", err)
	}

	id, err := obj.GetProperty(activeConnectionId)
	if err != nil {
		return nil, fmt.Errorf("GetActiveConnectionProperties ID property: %w", err)
	}

	uuid, err := obj.GetProperty(activeConnectionUuid)
	if err != nil {
		return nil, fmt.Errorf("GetActiveConnectionProperties UUID property: %w", err)
	}

	connType, err := obj.GetProperty(activeConnectionType)
	if err != nil {
		return nil, fmt.Errorf("GetActiveConnectionProperties Type property: %w", err)
	}
	devices, err := obj.GetProperty(activeConnectionDevices)
	if err != nil {
		return nil, fmt.Errorf("GetActiveConnectionProperties Devices property: %w", err)
	}
	state, err := obj.GetProperty(activeConnectionState)
	if err != nil {
		return nil, fmt.Errorf("GetActiveConnectionProperties State property: %w", err)
	}
	stateFlags, err := obj.GetProperty(activeConnectionStateFlags)
	if err != nil {
		return nil, fmt.Errorf("GetActiveConnectionProperties StateFlags property: %w", err)
	}
	connDefault, err := obj.GetProperty(activeConnectionDefault)
	if err != nil {
		return nil, fmt.Errorf("GetActiveConnectionProperties Default property: %w", err)
	}
	ip4config, err := obj.GetProperty(activeConnectionIp4Config)
	if err != nil {
		return nil, fmt.Errorf("GetActiveConnectionProperties Ip4Config property: %w", err)
	}
	dhcp4, err := obj.GetProperty(activeConnectionDhcp4Config)
	if err != nil {
		return nil, fmt.Errorf("GetActiveConnectionProperties Dhcp4Config property: %w", err)
	}
	default6, err := obj.GetProperty(activeConnectionDefault6)
	if err != nil {
		return nil, fmt.Errorf("GetActiveConnectionProperties Default6 property: %w", err)
	}
	ip6config, err := obj.GetProperty(activeConnectionIp6Config)
	if err != nil {
		return nil, fmt.Errorf("GetActiveConnectionProperties Ip6Config property: %w", err)
	}
	dhcp6, err := obj.GetProperty(activeConnectionDhcp6Config)
	if err != nil {
		return nil, fmt.Errorf("GetActiveConnectionProperties Dhcp6Config property: %w", err)
	}
	vpn, err := obj.GetProperty(activeConnectionVpn)
	if err != nil {
		return nil, fmt.Errorf("GetActiveConnectionProperties Vpn property: %w", err)
	}
	controller, err := obj.GetProperty(activeConnectionController)
	if err != nil {
		return nil, fmt.Errorf("GetActiveConnectionProperties Controller property: %w", err)
	}

	props := ActiveConnectionProperties{
		Connection:     conn.Value().(godbus.ObjectPath),
		SpecificObject: specificObj.Value().(godbus.ObjectPath),
		ID:             id.Value().(string),
		UUID:           uuid.Value().(string),
		Type:           connType.Value().(string),
		Devices:        devices.Value().([]godbus.ObjectPath),
		State:          state.Value().(uint32),
		StateFlags:     stateFlags.Value().(uint32),
		Default:        connDefault.Value().(bool),
		IP4Config:      ip4config.Value().(godbus.ObjectPath),
		DHCP4Config:    dhcp4.Value().(godbus.ObjectPath),
		Default6:       default6.Value().(bool),
		IP6Config:      ip6config.Value().(godbus.ObjectPath),
		DHCP6Config:    dhcp6.Value().(godbus.ObjectPath),
		VPN:            vpn.Value().(bool),
		Controller:     controller.Value().(godbus.ObjectPath),
	}
	return &props, nil
}
