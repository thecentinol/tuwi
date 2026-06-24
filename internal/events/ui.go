package events

import (
	"github.com/thecentinol/tuwi/internal/models"
)

type ShowPasswordModalMsg struct {
	Network *models.AccessPoint
}

type PasswordResultMsg struct {
	Cancelled bool
	Password  string
}
