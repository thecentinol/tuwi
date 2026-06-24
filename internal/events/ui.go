package events

import (
	nm "github.com/thecentinol/tuwi/internal/networkmanager"
)

type ShowPasswordModalMsg struct {
	Network *nm.AccessPoint
}

type PasswordResultMsg struct {
	Cancelled bool
	Password  string
}
