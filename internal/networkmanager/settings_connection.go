package networkmanager

import (
	"fmt"
	godbus "github.com/godbus/dbus/v5"

	"github.com/thecentinol/tuwi/internal/dbus"
)

func DeleteConnection(c *dbus.Client, conPath string) error {
	obj := c.Conn.Object(BaseServiceName, godbus.ObjectPath(conPath))
	call := obj.Call(CspDelete, 0)

	if call.Err != nil {
		return fmt.Errorf("DeleteConnection: %w", call.Err)
	}

	return nil
}
