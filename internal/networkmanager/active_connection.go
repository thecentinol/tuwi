package networkmanager

import (
	"fmt"

	godbus "github.com/godbus/dbus/v5"

	"github.com/thecentinol/tuwi/internal/dbus"
)

func GetActiveConnectionProperties(c *dbus.Client, path godbus.ObjectPath) (*ActiveConnection, error) {
	obj := c.Conn.Object(BaseServiceName, path)

	conn, err := obj.GetProperty(ActConConnection)
	if err != nil {
		return nil, fmt.Errorf("GetActiveConnectionProperties connection property: %w", err)
	}

	specificObj, err := obj.GetProperty(ActConSpecificObject)
	if err != nil {
		return nil, fmt.Errorf("GetActiveConnectionProperties SpecificObject property: %w", err)
	}

	id, err := obj.GetProperty(ActConId)
	if err != nil {
		return nil, fmt.Errorf("GetActiveConnectionProperties ID property: %w", err)
	}

	uuid, err := obj.GetProperty(ActConUuid)
	if err != nil {
		return nil, fmt.Errorf("GetActiveConnectionProperties UUID property: %w", err)
	}

	connType, err := obj.GetProperty(ActConType)
	if err != nil {
		return nil, fmt.Errorf("GetActiveConnectionProperties Type property: %w", err)
	}
	devices, err := obj.GetProperty(ActConDevices)
	if err != nil {
		return nil, fmt.Errorf("GetActiveConnectionProperties Devices property: %w", err)
	}
	state, err := obj.GetProperty(ActConState)
	if err != nil {
		return nil, fmt.Errorf("GetActiveConnectionProperties State property: %w", err)
	}
	stateFlags, err := obj.GetProperty(ActConStateFlags)
	if err != nil {
		return nil, fmt.Errorf("GetActiveConnectionProperties StateFlags property: %w", err)
	}
	connDefault, err := obj.GetProperty(ActConDefault)
	if err != nil {
		return nil, fmt.Errorf("GetActiveConnectionProperties Default property: %w", err)
	}
	ip4config, err := obj.GetProperty(ActConIp4Config)
	if err != nil {
		return nil, fmt.Errorf("GetActiveConnectionProperties Ip4Config property: %w", err)
	}
	dhcp4, err := obj.GetProperty(ActConDhcp4Config)
	if err != nil {
		return nil, fmt.Errorf("GetActiveConnectionProperties Dhcp4Config property: %w", err)
	}
	default6, err := obj.GetProperty(ActConDefault6)
	if err != nil {
		return nil, fmt.Errorf("GetActiveConnectionProperties Default6 property: %w", err)
	}
	ip6config, err := obj.GetProperty(ActConIp6Config)
	if err != nil {
		return nil, fmt.Errorf("GetActiveConnectionProperties Ip6Config property: %w", err)
	}
	dhcp6, err := obj.GetProperty(ActConDhcp6Config)
	if err != nil {
		return nil, fmt.Errorf("GetActiveConnectionProperties Dhcp6Config property: %w", err)
	}
	vpn, err := obj.GetProperty(ActConVpn)
	if err != nil {
		return nil, fmt.Errorf("GetActiveConnectionProperties Vpn property: %w", err)
	}
	controller, err := obj.GetProperty(ActConController)
	if err != nil {
		return nil, fmt.Errorf("GetActiveConnectionProperties Controller property: %w", err)
	}
	master, err := obj.GetProperty(ActConMaster)
	if err != nil {
		return nil, fmt.Errorf("GetActiveConnectionProperties Master property: %w", err)
	}

	ac := ActiveConnection{
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
		Master:         master.Value().(godbus.ObjectPath),
	}
	return &ac, nil
}
