package tui

import "github.com/charmbracelet/bubbles/key"

type keyMap struct {
	back   key.Binding
	editor key.Binding
	enter  key.Binding
	quit   key.Binding
}

var customKeys = keyMap{
	back: key.NewBinding(
		key.WithKeys("b"), key.WithHelp("b", "go back"),
	),
	editor: key.NewBinding(
		key.WithKeys("e"), key.WithHelp("e", "edit"),
	),
	enter: key.NewBinding(
		key.WithKeys("enter"), key.WithHelp("enter", "select"),
	),
	quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"), key.WithHelp("q/ctrl+c", "quit"),
	),
}
