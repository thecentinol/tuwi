package keybindings

import (
	"charm.land/bubbles/v2/key"
)

// converts a single keybinding to a bubbles keybinding
func (b Binding) ToBubbles() key.Binding {
	return key.NewBinding(
		key.WithKeys(b.Keys...),
		key.WithHelp(b.HelpKey, b.HelpDescription),
	)
}

// converts multiple keybindings to a bubbles keybindings
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

// Returns all keybindings
func (k Keybindings) All() []Binding {
	return []Binding{
		k.FocusedWifiSaved, k.FocusedWifiAvailable, k.Quit, k.WifiConnectSaved, k.WifiDisconnect,
		k.WifiForget, k.WifiAutoConnect, k.WifiScan, k.WifiConnectAvailable, k.LineUp, k.LineDown,
		k.GotoTop, k.GotoBottom, k.OpenKeybindingsModal, k.CloseKeybindingsModal, k.ErrorDismiss,
		k.PasswordSubmit, k.PasswordCancel,
	}
}

// Get all Keybindings and sort them by Group
func (k Keybindings) AllGrouped() []BindingGroup {
	var groups []BindingGroup
	var globalGroup BindingGroup
	var wifiGroup BindingGroup
	var tableGroup BindingGroup
	var kbGroup BindingGroup
	var errGroup BindingGroup
	var passwordGroup BindingGroup

	for _, b := range k.All() {
		if b.Group == "global" {
			globalGroup = BindingGroup{
				Title:    b.Group,
				Bindings: k.ByGroup("global"),
			}
		}
		if b.Group == "wifi" {
			wifiGroup = BindingGroup{
				Title:    b.Group,
				Bindings: k.ByGroup("wifi"),
			}
		}
		if b.Group == "table" {
			tableGroup = BindingGroup{
				Title:    b.Group,
				Bindings: k.ByGroup("table"),
			}
		}
		if b.Group == "keybindings" {
			kbGroup = BindingGroup{
				Title:    b.Group,
				Bindings: k.ByGroup("keybindings"),
			}
		}
		if b.Group == "error" {
			errGroup = BindingGroup{
				Title:    b.Group,
				Bindings: k.ByGroup("error"),
			}
		}
		if b.Group == "password" {
			passwordGroup = BindingGroup{
				Title:    b.Group,
				Bindings: k.ByGroup("password"),
			}
		}
	}

	groups = append(groups, globalGroup)
	groups = append(groups, wifiGroup)
	groups = append(groups, tableGroup)
	groups = append(groups, kbGroup)
	groups = append(groups, errGroup)
	groups = append(groups, passwordGroup)
	return groups
}

// Get keybindings by ViewName
func (k Keybindings) ByView(viewName string) []Binding {
	var binds []Binding
	for _, b := range k.All() {
		if b.ViewName == viewName {
			binds = append(binds, b)
		}
	}
	return binds
}

// Get keybindings by Group
func (k Keybindings) ByGroup(group string) []Binding {
	var binds []Binding
	for _, kb := range k.All() {
		if kb.Group == group {
			binds = append(binds, kb)
		}
	}
	return binds
}
