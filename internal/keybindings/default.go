package keybindings

var Default = Keybindings{
	// Global
	FocusedWifiSaved: Binding{
		ViewName:        globalView,
		Group:           globalGroup,
		Keys:            []string{"1"},
		HelpKey:         "1",
		HelpDescription: "",
		Description:     "Switch focus to saved wifi connections pane",
	},
	FocusedWifiAvailable: Binding{
		ViewName:        globalView,
		Group:           globalGroup,
		Keys:            []string{"2"},
		HelpKey:         "2",
		HelpDescription: "",
		Description:     "Switch focus to available wifi networks pane",
	},
	Cancel: Binding{
		ViewName:        globalView,
		Group:           globalGroup,
		Keys:            []string{keyEsc},
		HelpKey:         keyEsc,
		HelpDescription: "cancel",
		Description:     "Cancel",
	},
	OpenKeybindingsModal: Binding{
		ViewName:        globalView,
		Group:           globalGroup,
		Keys:            []string{"?"},
		HelpKey:         "?",
		HelpDescription: "keybindings",
		Description:     "Open keybindings menu",
	},
	Quit: Binding{
		ViewName:        globalView,
		Group:           globalGroup,
		Keys:            []string{"q", "ctrl+c"},
		HelpKey:         "q",
		HelpDescription: "quit",
		Description:     "Quit the application",
	},

	// WiFi saved
	WifiConnectSaved: Binding{
		ViewName:        wifiSavedView,
		Group:           wifiSavedGroup,
		Keys:            []string{keyEnter},
		HelpKey:         keyEnter,
		HelpDescription: "connect",
		Description:     "Connect",
	},
	WifiDisconnect: Binding{
		ViewName:        wifiSavedView,
		Group:           wifiSavedGroup,
		Keys:            []string{"d"},
		HelpKey:         "d",
		HelpDescription: "disconnect",
		Description:     "Disconnect from current WiFi network",
	},
	WifiForget: Binding{
		ViewName:        wifiSavedView,
		Group:           wifiSavedGroup,
		Keys:            []string{"f"},
		HelpKey:         "f",
		HelpDescription: "forget",
		Description:     "Remove the selected WiFi connection",
	},
	WifiAutoConnect: Binding{
		ViewName:        wifiSavedView,
		Group:           wifiSavedGroup,
		Keys:            []string{"a"},
		HelpKey:         "a",
		HelpDescription: "auto-connect",
		Description:     "Toggle auto-connect on/off for the selected WiFi connection",
	},

	// WiFi available
	WifiScan: Binding{
		ViewName:        wifiAvailableView,
		Group:           wifiAvailableGroup,
		Keys:            []string{"s"},
		HelpKey:         "s",
		HelpDescription: "Scan",
		Description:     "Scan for nearby WiFi networks",
	},
	WifiConnectAvailable: Binding{
		ViewName:        wifiAvailableView,
		Group:           wifiAvailableGroup,
		Keys:            []string{keyEnter},
		HelpKey:         keyEnter,
		HelpDescription: "connect",
		Description:     "Connect",
	},

	// Navigation
	Up: Binding{
		Group:           navGroup,
		Keys:            []string{"k", "up"},
		HelpKey:         "↑",
		HelpDescription: "up",
		Description:     "Move up by one line",
	},
	Down: Binding{
		Group:           navGroup,
		Keys:            []string{"j", "down"},
		HelpKey:         "",
		HelpDescription: "down",
		Description:     "Move down by one line",
	},
	GotoTop: Binding{
		Group:           navGroup,
		Keys:            []string{"t"},
		HelpKey:         "t",
		HelpDescription: "top",
		Description:     "Move/Scroll to top",
	},
	GotoBottom: Binding{
		Group:           navGroup,
		Keys:            []string{"b"},
		HelpKey:         "b",
		HelpDescription: "bottom",
		Description:     "Move/Scroll to bottom",
	},

	// Error modal
	ErrorDismiss: Binding{
		ViewName:        errModalView,
		Group:           errModalGroup,
		Keys:            []string{keyEsc, keyEnter},
		HelpKey:         keyEsc,
		HelpDescription: "dismiss",
		Description:     "Dismiss the error modal",
	},

	// Password modal
	PasswordVisibility: Binding{
		ViewName:        passwordView,
		Group:           passwordGroup,
		Keys:            []string{"ctrl+t"},
		HelpKey:         "ctrl+t",
		HelpDescription: "hide/show password",
		Description:     "Toggle the password input's hidden state",
	},
	PasswordSubmit: Binding{
		ViewName:        passwordView,
		Group:           passwordGroup,
		Keys:            []string{keyEnter},
		HelpKey:         keyEnter,
		HelpDescription: "submit",
		Description:     "Submit password",
	},
}
