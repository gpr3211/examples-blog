package model

import (
	tea "github.com/charmbracelet/bubbletea"
)

type RootModel struct {
	Current tea.Model
}

func (m RootModel) Init() tea.Cmd {
	return m.Current.Init()
}

func (m RootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.Current, cmd = m.Current.Update(msg)

	// If current model returns a special message, handle transition
	switch msg.(type) {
	case goToFilePickerMsg:
		m.Current = newFilePickerModel()
		return m, m.Current.Init()
	}

	return m, cmd
}

func (m RootModel) View() string {
	return m.Current.View()
}

type MainMenuModel struct{}

type goToFilePickerMsg struct{}

func (m MainMenuModel) Init() tea.Cmd {
	return nil
}

func (m MainMenuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case tea.KeyMsg:
		return m, func() tea.Msg { return goToFilePickerMsg{} }
	}
	return m, nil
}

func (m MainMenuModel) View() string {
	return "\n Press any key to open File Picker.\n"
}
