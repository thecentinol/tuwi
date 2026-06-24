package events

import (
	"github.com/thecentinol/tuwi/internal/models"
)

type SavedWifiMsg struct{ Networks []models.AccessPoint }
type AvailableWifiMsg struct{ Networks []models.AccessPoint }
type WifiConnectReqMsg struct {
	Network  *models.AccessPoint
	Password string
}
