package cli

import (
	"io"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

type palette struct {
	header  lipgloss.Style
	name    lipgloss.Style
	muted   lipgloss.Style
	running lipgloss.Style
	stopped lipgloss.Style
}

func newPalette(cmd *cobra.Command, out io.Writer) palette {
	enabled := colorsEnabled(cmd, out)
	if !enabled {
		return palette{}
	}
	return palette{
		header:  lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("12")),
		name:    lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("14")),
		muted:   lipgloss.NewStyle().Foreground(lipgloss.Color("8")),
		running: lipgloss.NewStyle().Foreground(lipgloss.Color("10")),
		stopped: lipgloss.NewStyle().Foreground(lipgloss.Color("11")),
	}
}

func colorsEnabled(cmd *cobra.Command, out io.Writer) bool {
	noColor, _ := cmd.Flags().GetBool("no-color")
	if noColor {
		return false
	}
	if os.Getenv("NO_COLOR") != "" || os.Getenv("CLICOLOR") == "0" || strings.EqualFold(os.Getenv("TERM"), "dumb") {
		return false
	}
	if os.Getenv("FORCE_COLOR") != "" {
		return true
	}
	f, ok := out.(*os.File)
	if !ok {
		return false
	}
	info, err := f.Stat()
	if err != nil {
		return false
	}
	return (info.Mode() & os.ModeCharDevice) != 0
}

func (p palette) Header(s string) string {
	if p.header.String() == "" {
		return s
	}
	return p.header.Render(s)
}

func (p palette) Name(s string) string {
	if p.name.String() == "" {
		return s
	}
	return p.name.Render(s)
}

func (p palette) Muted(s string) string {
	if p.muted.String() == "" {
		return s
	}
	return p.muted.Render(s)
}

func (p palette) State(s string) string {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "running":
		if p.running.String() == "" {
			return s
		}
		return p.running.Render(s)
	case "stopped":
		if p.stopped.String() == "" {
			return s
		}
		return p.stopped.Render(s)
	default:
		return p.Muted(s)
	}
}

