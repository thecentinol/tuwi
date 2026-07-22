package keybindings

import (
	"charm.land/bubbles/v2/key"
)

func DefaultKeybindings() Keybindings {
	return Keybindings{
		// Global
		FocusedWifiSaved: Binding{
			ViewName:        "",
			Group:           "global",
			Keys:            []string{"1"},
			HelpKey:         "",
			HelpDescription: "",
			Description:     "Switch focus to saved wifi connections pane",
		},
		FocusedWifiAvailable: Binding{
			ViewName:        "",
			Group:           "global",
			Keys:            []string{"2"},
			HelpKey:         "",
			HelpDescription: "",
			Description:     "Switch focus to available wifi networks pane",
		},
		Quit: Binding{
			ViewName:        "",
			Group:           "global",
			Keys:            []string{"q", "ctrl+c"},
			HelpKey:         "q",
			HelpDescription: "quit",
			Description:     "Quit the application",
		},

		// WiFi
		WifiConnectSaved: Binding{
			ViewName:        "wifi-saved",
			Group:           "wifi",
			Keys:            []string{"enter"},
			HelpKey:         "Enter",
			HelpDescription: "connect",
			Description:     "Connect to a saved wifi connection",
		},
		WifiDisconnect: Binding{
			ViewName:        "wifi-saved",
			Group:           "wifi",
			Keys:            []string{"d"},
			HelpKey:         "d",
			HelpDescription: "diconnect",
			Description:     "Disconnect from current wifi network",
		},
		WifiForget: Binding{
			ViewName:        "wifi-saved",
			Group:           "wifi",
			Keys:            []string{"f"},
			HelpKey:         "f",
			HelpDescription: "forget",
			Description:     "Remove the selected wifi connection",
		},
		WifiAutoConnect: Binding{
			ViewName:        "wifi-saved",
			Group:           "wifi",
			Keys:            []string{"a"},
			HelpKey:         "a",
			HelpDescription: "toggle auto-connect",
			Description:     "Toggle auto-connect state for the selected wifi connection",
		},
		WifiScan: Binding{
			ViewName:        "wifi-available",
			Group:           "wifi",
			Keys:            []string{"s"},
			HelpKey:         "s",
			HelpDescription: "Scan",
			Description:     "Scan for nearby wifi networks",
		},
		WifiConnectAvailable: Binding{
			ViewName:        "wifi-available",
			Group:           "wifi",
			Keys:            []string{"enter"},
			HelpKey:         "enter",
			HelpDescription: "connect",
			Description:     "Connect to a new wifi network",
		},

		// table navigation
		// ViewName is left nil because multiple models will use a table.
		LineUp: Binding{
			Group:           "table",
			Keys:            []string{"k", "up"},
			HelpKey:         "↑",
			HelpDescription: "move up",
			Description:     "Move up by one line",
		},
		LineDown: Binding{
			Group:           "table",
			Keys:            []string{"j", "down"},
			HelpKey:         "",
			HelpDescription: "move down",
			Description:     "Move down by one line",
		},
		GotoTop: Binding{
			Group:           "table",
			Keys:            []string{"t"},
			HelpKey:         "t",
			HelpDescription: "move to top",
			Description:     "Go to the top of the table",
		},
		GotoBottom: Binding{
			Group:           "table",
			Keys:            []string{"b"},
			HelpKey:         "b",
			HelpDescription: "move to bottom",
			Description:     "Go to the bottom of the table",
		},

		// Error modal
		ErrorDismiss: Binding{
			ViewName:        "error-modal",
			Group:           "error",
			Keys:            []string{"esc", "enter"},
			HelpKey:         "esc",
			HelpDescription: "dismiss",
			Description:     "Dismiss the error modal",
		},

		// Password modal
		PasswordVisibility: Binding{
			ViewName:        "password-modal",
			Group:           "password",
			Keys:            []string{"ctrl+t"},
			HelpKey:         "ctrl+t",
			HelpDescription: "hide/show password",
			Description:     "Toggle the password inputs hidden state",
		},
		PasswordSubmit: Binding{
			ViewName:        "password-modal",
			Group:           "password",
			Keys:            []string{"enter"},
			HelpKey:         "enter",
			HelpDescription: "submit",
			Description:     "Submit password",
		},
		PasswordCancel: Binding{
			ViewName:        "password-modal",
			Group:           "password",
			Keys:            []string{"esc"},
			HelpKey:         "esc",
			HelpDescription: "cancel",
			Description:     "Cancel",
		},
	}
}

type Keybindings struct {
	// NOTE: if you change a bindings ViewName, be sure to update the Focused consts
	// in the root model

	// Global
	FocusedWifiSaved     Binding
	FocusedWifiAvailable Binding
	Quit                 Binding

	// WiFi
	WifiConnectSaved     Binding
	WifiDisconnect       Binding
	WifiForget           Binding
	WifiAutoConnect      Binding
	WifiScan             Binding
	WifiConnectAvailable Binding

	// Table navigation
	LineUp     Binding
	LineDown   Binding
	GotoTop    Binding
	GotoBottom Binding

	// Error modal
	ErrorDismiss Binding

	// Password modal
	PasswordVisibility Binding
	PasswordSubmit     Binding
	PasswordCancel     Binding
}

type Binding struct {
	// used to match against which View is focused so we can display keybindings for the
	// focused View at the top of the keybindings modal.
	ViewName string
	Group    string // used for grouping in the keybindings modal

	Keys []string

	HelpKey         string // which key to display in the help view
	HelpDescription string // the help text for bubbles key.Binding.Help
	Description     string // the longer full description displayed in the keybindings modal
}

// converts the keybind to a bubbles keybinding
func (b Binding) ToBubbles() key.Binding {
	return key.NewBinding(
		key.WithKeys(b.Keys...),
		key.WithHelp(b.HelpKey, b.HelpDescription),
	)
}

// returns all keybindings
func (k Keybindings) All() []Binding {
	return []Binding{
		k.FocusedWifiSaved, k.FocusedWifiAvailable, k.Quit, k.WifiConnectSaved, k.WifiDisconnect,
		k.WifiForget, k.WifiAutoConnect, k.WifiScan, k.WifiConnectAvailable, k.LineUp, k.LineDown,
		k.GotoTop, k.GotoBottom, k.ErrorDismiss, k.PasswordSubmit, k.PasswordCancel,
	}
}

// get keybindings by ViewName
func (k Keybindings) ByView(viewName string) []Binding {
	var binds []Binding
	for _, b := range k.All() {
		if b.ViewName == viewName {
			binds = append(binds, b)
		}
	}
	return binds
}
