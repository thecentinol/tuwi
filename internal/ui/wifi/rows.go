package wifi

import (
	"strconv"

	"charm.land/bubbles/v2/table"

	nm "github.com/thecentinol/tuwi/internal/networkmanager"
	"github.com/thecentinol/tuwi/internal/service"
)

func setSavedRows(nc []nm.NearbyConnection, active *nm.ActiveConnection) []table.Row {
	rows := make([]table.Row, 0, len(nc))

	for _, n := range nc {
		isNearby := n.AP != nil
		isActive := false
		var security string
		var strength uint8

		if active != nil {
			isActive = active.Connection == n.Connection.ConnectionPath
		}

		if isNearby {
			security = service.DetermineSecurityType(n.AP.Flags, n.AP.WpaFlags, n.AP.RsnFlags)
			strength = n.AP.Strength
		} else {
			security = n.Connection.KeyMgmt
			strength = 0
		}

		rows = append(rows, table.Row{
			n.Connection.SSID,
			service.DetermineStatus(isNearby, isActive),
			security,
			strconv.FormatBool(n.Connection.Hidden),
			service.FormatStrength(strength),
		})
	}

	return rows
}

func setAvailableRows(networks []nm.AccessPoint) []table.Row {
	rows := make([]table.Row, 0, len(networks))

	for _, n := range networks {
		strength := strconv.FormatInt(int64(n.Strength), 10)
		securityType := service.DetermineSecurityType(n.Flags, n.WpaFlags, n.RsnFlags)
		channel := service.DetermineChannel(n.Frequency)

		rows = append(rows, table.Row{
			service.FormatSSID(n.SSID),
			securityType,
			service.DetermineFrequency(n.Frequency),
			strength + "%",
			service.FormatChannel(channel),
		})
	}

	return rows
}
