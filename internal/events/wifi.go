package events

import (
	nm "github.com/thecentinol/tuwi/internal/networkmanager"
)

type SavedWifiMsg struct{ Networks []nm.AccessPoint }
type AvailableWifiMsg struct{ Networks []nm.AccessPoint }
type WifiConnectReqMsg struct {
	Network  *nm.AccessPoint
	Password string
}
