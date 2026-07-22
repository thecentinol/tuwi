package keybindings

import (
	"charm.land/bubbles/v2/key"
)

// converts a single keybind to a bubbles keybinding
func (b Binding) ToBubbles() key.Binding {
	return key.NewBinding(
		key.WithKeys(b.Keys...),
		key.WithHelp(b.HelpKey, b.HelpDescription),
	)
}

func ToBubblesBatch(keys []Binding) []key.Binding {
	var binds []key.Binding

	for _, k := range keys {
		bind := key.NewBinding(
			key.WithKeys(k.Keys...),
			key.WithHelp(k.HelpKey, k.Description),
		)
		binds = append(binds, bind)
	}

	return binds
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

func (k Keybindings) ByGroup(group string) []Binding {
	var binds []Binding
	for _, kb := range k.All() {
		if kb.Group == group {
			binds = append(binds, kb)
		}
	}
	return binds
}
