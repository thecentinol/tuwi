package dbus

const (
	baseServiceName = "org.freedesktop.NetworkManager"
	baseObjPath     = "/org/freedesktop/NetworkManager"
	deviceType      = baseServiceName + ".Device.DeviceType"
	deviceWireless  = baseServiceName + ".Device.Wireless"
	deviceWired     = baseServiceName + ".Device.Wired"
)
