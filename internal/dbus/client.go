package dbus

import (
	"fmt"
	"log/slog"

	"github.com/godbus/dbus/v5"
)

type Client struct {
	Conn *dbus.Conn
}

func NewClient() (*Client, error) {
	conn, err := dbus.SystemBus()
	if err != nil {
		slog.Error("Connection to system bus failed", "error", err)
		return nil, fmt.Errorf("D-Bus: failed to connect to system bus: %w", err)
	}

	return &Client{Conn: conn}, nil
}

func (c *Client) Close() error {
	return c.Conn.Close()
}
