package dbus

import (
	"fmt"

	"github.com/godbus/dbus/v5"
)

type Client struct {
	Conn *dbus.Conn
}

func NewClient() (*Client, error) {
	conn, err := dbus.SystemBus()
	if err != nil {
		return nil, fmt.Errorf("DBus Client: failed to connect to system bus: %w", err)
	}
	return &Client{Conn: conn}, nil
}

func (c *Client) Close() error {
	return c.Conn.Close()
}
