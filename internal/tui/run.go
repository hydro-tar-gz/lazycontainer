package tui

import (
	tea "github.com/charmbracelet/bubbletea"

	"lazycontainer/internal/core"
)

func Run(backend core.Backend) error {
	p := tea.NewProgram(NewModel(backend), tea.WithAltScreen())
	_, err := p.Run()
	return err
}
