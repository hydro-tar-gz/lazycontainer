package tui

import (
	"strings"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
)

type theme struct {
	accent      lipgloss.AdaptiveColor
	panelBorder lipgloss.AdaptiveColor
	muted       lipgloss.AdaptiveColor
	success     lipgloss.AdaptiveColor
	danger      lipgloss.AdaptiveColor
	selectedFG  lipgloss.AdaptiveColor
	selectedBG  lipgloss.AdaptiveColor

	activeItem   lipgloss.Style
	inactiveItem lipgloss.Style
	title        lipgloss.Style
	searchLabel  lipgloss.Style
	panel        lipgloss.Style
	footer       lipgloss.Style
}

func newTheme() theme {
	t := theme{
		accent:      lipgloss.AdaptiveColor{Light: "#005f87", Dark: "#5fd7ff"},
		panelBorder: lipgloss.AdaptiveColor{Light: "#9ca3af", Dark: "#475569"},
		muted:       lipgloss.AdaptiveColor{Light: "#6b7280", Dark: "#94a3b8"},
		success:     lipgloss.AdaptiveColor{Light: "#047857", Dark: "#34d399"},
		danger:      lipgloss.AdaptiveColor{Light: "#b91c1c", Dark: "#f87171"},
		selectedFG:  lipgloss.AdaptiveColor{Light: "#ffffff", Dark: "#0f111a"},
		selectedBG:  lipgloss.AdaptiveColor{Light: "#0087af", Dark: "#87d7ff"},
	}

	t.activeItem = lipgloss.NewStyle().Bold(true).Foreground(t.accent).Underline(true)
	t.inactiveItem = lipgloss.NewStyle().Foreground(t.muted)
	t.title = lipgloss.NewStyle().Bold(true).Foreground(t.accent)
	t.searchLabel = lipgloss.NewStyle().Bold(true).Foreground(t.accent)
	t.panel = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(t.panelBorder).Padding(0, 1)
	t.footer = lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), true, false, false, false).
		BorderForeground(t.panelBorder).
		Foreground(t.muted)

	return t
}

func (t theme) applyTableStyles(tbl *table.Model) {
	s := table.DefaultStyles()
	s.Header = s.Header.
		Bold(true).
		Foreground(t.accent).
		BorderStyle(lipgloss.NormalBorder()).
		BorderBottom(true).
		BorderForeground(t.accent)
	s.Selected = s.Selected.Bold(true).Foreground(t.selectedFG).Background(t.selectedBG)
	s.Cell = s.Cell.Foreground(t.muted)
	tbl.SetStyles(s)
}

func (t theme) status(msg string) lipgloss.Style {
	msg = strings.ToLower(msg)
	switch {
	case strings.Contains(msg, "error"), strings.Contains(msg, "failed"):
		return lipgloss.NewStyle().Bold(true).Foreground(t.danger)
	case strings.Contains(msg, "loaded"),
		strings.Contains(msg, "complete"),
		strings.Contains(msg, "started"),
		strings.Contains(msg, "stopped"),
		strings.Contains(msg, "deleted"):
		return lipgloss.NewStyle().Foreground(t.success)
	default:
		return lipgloss.NewStyle().Foreground(t.muted)
	}
}

func (t theme) pickItem(label string, active bool) string {
	if active {
		return t.activeItem.Render(label)
	}
	return t.inactiveItem.Render(label)
}

