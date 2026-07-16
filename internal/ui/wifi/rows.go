package wifi

import (
	"charm.land/bubbles/v2/table"
	"strconv"

	nm "github.com/thecentinol/tuwi/internal/networkmanager"
	wifidomain "github.com/thecentinol/tuwi/internal/wifi"
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
			security = wifidomain.DetermineSecurityType(n.AP.Flags, n.AP.WpaFlags, n.AP.RsnFlags)
			strength = n.AP.Strength
		} else {
			security = n.Connection.KeyMgmt
			strength = 0
		}

		rows = append(rows, table.Row{
			n.Connection.SSID,
			wifidomain.DetermineStatus(isNearby, isActive),
			security,
			strconv.FormatBool(n.Connection.Hidden),
			wifidomain.FormatStrength(strength),
		})
	}

	return rows
}

func setAvailableRows(networks []nm.AccessPoint) []table.Row {
	rows := make([]table.Row, 0, len(networks))

	for _, n := range networks {
		strength := strconv.FormatInt(int64(n.Strength), 10)
		securityType := wifidomain.DetermineSecurityType(n.Flags, n.WpaFlags, n.RsnFlags)
		channel := wifidomain.DetermineChannel(n.Frequency)

		rows = append(rows, table.Row{
			wifidomain.FormatSSID(n.SSID),
			securityType,
			wifidomain.DetermineFrequency(n.Frequency),
			strength + "%",
			wifidomain.FormatChannel(channel),
		})
	}

	return rows
}
