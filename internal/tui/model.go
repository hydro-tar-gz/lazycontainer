package tui

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"lazycontainer/internal/core"
)

const (
	modeUI = iota
	modeCLI
)

type instancesLoadedMsg struct {
	instances []core.Instance
	err       error
}

type tabLoadedMsg struct {
	tab     int
	name    string
	content string
	err     error
}

type actionDoneMsg struct {
	status string
	err    error
}

type execDoneMsg struct {
	out string
	err error
}

type shellDoneMsg struct {
	err   error
	shell string
}

type cliCmdDoneMsg struct {
	out string
	err error
}

type model struct {
	backend core.Backend
	theme   theme

	width  int
	height int
	leftW  int
	rightW int
	bodyH  int

	instances []core.Instance
	filtered  []core.Instance

	table table.Model
	vp    viewport.Model

	search        textinput.Model
	searchFocused bool

	execInput textinput.Model
	showExec  bool

	cliInput textinput.Model

	showDeleteConfirm bool

	modes      []string
	activeMode int

	tabs      []string
	activeTab int

	rightContent string
	status       string
	loading      bool
}

func NewModel(backend core.Backend) model {
	uiTheme := newTheme()
	cols := []table.Column{
		{Title: "NAME", Width: 24},
		{Title: "STATE", Width: 10},
		{Title: "IP", Width: 16},
		{Title: "IMAGE", Width: 28},
	}
	t := table.New(table.WithColumns(cols), table.WithFocused(true), table.WithHeight(12))
	t.SetRows([]table.Row{})
	uiTheme.applyTableStyles(&t)

	s := textinput.New()
	s.Placeholder = "Filter by name"
	s.Blur()

	e := textinput.New()
	e.Placeholder = "Command (example: uname -a)"
	e.Blur()

	cli := textinput.New()
	cli.Placeholder = "Type lazy command args (example: ls --json)"
	cli.Blur()

	vp := viewport.New(20, 10)
	vp.SetContent("Select an instance")

	return model{
		backend:   backend,
		theme:     uiTheme,
		table:     t,
		vp:        vp,
		search:    s,
		execInput: e,
		cliInput:  cli,
		modes:     []string{"UI", "CLI"},
		tabs:      []string{"Info", "Logs", "Snapshots"},
		status:    "Loading instances...",
		loading:   true,
	}
}

func (m model) Init() tea.Cmd {
	return refreshInstancesCmd(m.backend)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.applyLayout()
		return m, nil

	case instancesLoadedMsg:
		m.loading = false
		if msg.err != nil {
			m.status = "Refresh failed: " + msg.err.Error()
			return m, nil
		}
		m.instances = msg.instances
		m.applyFilter()
		m.status = fmt.Sprintf("Loaded %d instance(s)", len(msg.instances))
		if name := m.selectedName(); name != "" {
			return m, loadTabCmd(m.backend, name, m.activeTab)
		}
		m.setRightContent("No instance selected")
		return m, nil

	case tabLoadedMsg:
		if msg.tab != m.activeTab || msg.name != m.selectedName() {
			return m, nil
		}
		if msg.err != nil {
			m.setRightContent(msg.err.Error())
			return m, nil
		}
		m.setRightContent(msg.content)
		return m, nil

	case actionDoneMsg:
		if msg.err != nil {
			m.status = "Error: " + msg.err.Error()
			m.loading = false
			return m, nil
		}
		m.status = msg.status
		m.loading = true
		return m, refreshInstancesCmd(m.backend)

	case execDoneMsg:
		if msg.err != nil {
			m.status = "Exec failed: " + msg.err.Error()
			return m, nil
		}
		m.status = "Command complete"
		m.setRightContent(msg.out)
		return m, nil

	case cliCmdDoneMsg:
		if msg.err != nil {
			m.status = "CLI failed: " + msg.err.Error()
			if strings.TrimSpace(msg.out) != "" {
				m.setRightContent(strings.TrimSpace(msg.out))
			}
			return m, nil
		}
		m.status = "CLI command complete"
		if strings.TrimSpace(msg.out) == "" {
			m.setRightContent("(no output)")
		} else {
			m.setRightContent(strings.TrimSpace(msg.out))
		}
		return m, nil

	case shellDoneMsg:
		if msg.err != nil && msg.shell == "bash" {
			name := m.selectedName()
			if name != "" {
				return m, shellCmd(name, "sh")
			}
		}
		if msg.err != nil {
			m.status = "Shell failed: " + msg.err.Error()
			return m, nil
		}
		m.status = "Shell exited"
		return m, nil

	case tea.KeyMsg:
		if msg.String() == "tab" {
			m.switchMode()
			return m, nil
		}

		if msg.String() == "q" || msg.String() == "ctrl+c" {
			return m, tea.Quit
		}

		if m.activeMode == modeCLI {
			switch msg.String() {
			case "enter":
				cmdLine := strings.TrimSpace(m.cliInput.Value())
				if cmdLine == "" {
					return m, nil
				}
				m.status = "Running: lazy " + cmdLine
				return m, runLocalLazyCmd(cmdLine)
			case "esc":
				m.cliInput.SetValue("")
				return m, nil
			}
			var cmd tea.Cmd
			m.cliInput, cmd = m.cliInput.Update(msg)
			return m, cmd
		}

		if m.showExec {
			switch msg.String() {
			case "esc":
				m.showExec = false
				m.execInput.Blur()
				return m, nil
			case "enter":
				name := m.selectedName()
				cmdLine := strings.TrimSpace(m.execInput.Value())
				m.showExec = false
				m.execInput.Blur()
				if name == "" || cmdLine == "" {
					return m, nil
				}
				m.status = "Running command..."
				return m, runExecCmd(m.backend, name, strings.Fields(cmdLine))
			}
			var cmd tea.Cmd
			m.execInput, cmd = m.execInput.Update(msg)
			return m, cmd
		}

		if m.showDeleteConfirm {
			switch msg.String() {
			case "y", "Y", "enter":
				name := m.selectedName()
				m.showDeleteConfirm = false
				if name == "" {
					return m, nil
				}
				m.status = "Deleting " + name + "..."
				return m, deleteCmd(m.backend, name)
			case "n", "N", "esc":
				m.showDeleteConfirm = false
				return m, nil
			}
		}

		if m.searchFocused {
			switch msg.String() {
			case "esc", "enter":
				m.searchFocused = false
				m.search.Blur()
				return m, nil
			}
			var cmd tea.Cmd
			m.search, cmd = m.search.Update(msg)
			m.applyFilter()
			return m, cmd
		}

		switch msg.String() {
		case "/":
			m.searchFocused = true
			m.search.Focus()
			return m, nil
		case "r":
			m.loading = true
			m.status = "Refreshing..."
			return m, refreshInstancesCmd(m.backend)
		case "t":
			m.activeTab = (m.activeTab + 1) % len(m.tabs)
			if name := m.selectedName(); name != "" {
				m.setRightContent("Loading " + m.tabs[m.activeTab] + "...")
				return m, loadTabCmd(m.backend, name, m.activeTab)
			}
			return m, nil
		case "enter":
			if name := m.selectedName(); name != "" {
				m.status = "Opening shell..."
				return m, shellCmd(name, "bash")
			}
			return m, nil
		case "e":
			if m.selectedName() == "" {
				return m, nil
			}
			m.showExec = true
			m.execInput.SetValue("")
			m.execInput.Focus()
			return m, nil
		case "s":
			inst := m.selectedInstance()
			if inst == nil {
				return m, nil
			}
			m.status = "Toggling state..."
			return m, toggleStateCmd(m.backend, *inst)
		case "d":
			if m.selectedName() == "" {
				return m, nil
			}
			m.showDeleteConfirm = true
			return m, nil
		}

		prev := m.selectedName()
		var cmd tea.Cmd
		m.table, cmd = m.table.Update(msg)
		if now := m.selectedName(); now != "" && now != prev {
			return m, tea.Batch(cmd, loadTabCmd(m.backend, now, m.activeTab))
		}
		return m, cmd
	}

	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m model) View() string {
	modeParts := make([]string, 0, len(m.modes))
	for i, name := range m.modes {
		modeParts = append(modeParts, m.theme.pickItem(name, i == m.activeMode))
	}
	header := m.theme.title.Render("Mode: ") + strings.Join(modeParts, " | ")
	status := m.theme.status(m.status).Render(m.status)

	if m.activeMode == modeCLI {
		body := m.theme.panel.Render(
			m.theme.title.Render("CLI") + "\n" +
				m.theme.inactiveItem.Render("Run command args after `lazy` (example: ls --json)") + "\n" +
				m.cliInput.View() + "\n\n" + m.vp.View(),
		)
		help := "enter: run command • esc: clear input • tab: switch mode • q: quit"
		footer := m.theme.footer.Render(help)
		return header + "\n" + body + "\n" + status + "\n" + footer
	}

	searchLabel := "Search"
	if m.searchFocused {
		searchLabel = "Search (focused)"
	}
	top := header + "\n" + m.theme.searchLabel.Render(searchLabel+": ") + m.search.View()

	leftTitle := "Instances"
	if m.loading {
		leftTitle += " (loading...)"
	}
	leftPane := m.theme.panel.Render(
		m.theme.title.Render(leftTitle) + "\n" + m.table.View(),
	)

	tabParts := make([]string, 0, len(m.tabs))
	for i, t := range m.tabs {
		tabParts = append(tabParts, m.theme.pickItem(t, i == m.activeTab))
	}
	rightPane := m.theme.panel.Render(
		strings.Join(tabParts, " | ") + "\n" + m.vp.View(),
	)

	body := lipgloss.JoinHorizontal(lipgloss.Top, leftPane, rightPane)

	help := "j/k or arrows: move • /: search • r: refresh • enter: shell • e: exec • s: start/stop • d: delete • t: detail tabs • tab: mode • q: quit"
	if m.showDeleteConfirm {
		help = "Delete selected instance? (y/N)"
	}
	if m.showExec {
		help = "Exec command: " + m.execInput.View()
	}

	footer := m.theme.footer.Render(help)
	return top + "\n" + body + "\n" + status + "\n" + footer
}

func (m *model) setRightContent(s string) {
	m.rightContent = s
	m.vp.SetContent(s)
}

func (m *model) applyLayout() {
	m.leftW = max(40, m.width*55/100)
	m.rightW = max(30, m.width-m.leftW-4)
	m.bodyH = max(8, m.height-8)
	m.table.SetWidth(m.leftW - 4)
	m.table.SetHeight(m.bodyH)
	if m.activeMode == modeCLI {
		m.vp.Width = max(24, m.width-8)
		m.vp.Height = max(8, m.height-10)
		return
	}
	m.vp.Width = m.rightW - 4
	m.vp.Height = m.bodyH
}

func (m *model) switchMode() {
	m.activeMode = (m.activeMode + 1) % len(m.modes)
	if m.activeMode == modeCLI {
		m.searchFocused = false
		m.search.Blur()
		m.showExec = false
		m.showDeleteConfirm = false
		m.cliInput.Focus()
		m.status = "CLI mode"
	} else {
		m.cliInput.Blur()
		m.status = "UI mode"
	}
	m.applyLayout()
}

func (m *model) applyFilter() {
	q := strings.ToLower(strings.TrimSpace(m.search.Value()))
	rows := make([]table.Row, 0, len(m.instances))
	m.filtered = m.filtered[:0]
	for _, inst := range m.instances {
		if q != "" && !strings.Contains(strings.ToLower(inst.Name), q) {
			continue
		}
		m.filtered = append(m.filtered, inst)
		rows = append(rows, table.Row{inst.Name, inst.State, inst.IP, inst.Image})
	}
	m.table.SetRows(rows)
	if len(rows) == 0 {
		m.setRightContent("No matching instances")
		return
	}
	if m.table.Cursor() >= len(rows) {
		m.table.SetCursor(len(rows) - 1)
	}
}

func (m model) selectedName() string {
	idx := m.table.Cursor()
	if idx < 0 || idx >= len(m.filtered) {
		return ""
	}
	return m.filtered[idx].Name
}

func (m model) selectedInstance() *core.Instance {
	idx := m.table.Cursor()
	if idx < 0 || idx >= len(m.filtered) {
		return nil
	}
	inst := m.filtered[idx]
	return &inst
}

func refreshInstancesCmd(backend core.Backend) tea.Cmd {
	return func() tea.Msg {
		instances, err := backend.ListInstances(context.Background())
		return instancesLoadedMsg{instances: instances, err: err}
	}
}

func loadTabCmd(backend core.Backend, name string, tab int) tea.Cmd {
	return func() tea.Msg {
		switch tab {
		case 0:
			instances, err := backend.ListInstances(context.Background())
			if err != nil {
				return tabLoadedMsg{tab: tab, name: name, err: err}
			}
			for _, inst := range instances {
				if inst.Name == name {
					info := fmt.Sprintf("Name: %s\nState: %s\nIP: %s\nImage: %s", inst.Name, inst.State, inst.IP, inst.Image)
					return tabLoadedMsg{tab: tab, name: name, content: info}
				}
			}
			return tabLoadedMsg{tab: tab, name: name, content: "Instance not found"}
		case 1:
			out, err := backend.Logs(context.Background(), name)
			return tabLoadedMsg{tab: tab, name: name, content: out, err: err}
		case 2:
			out, err := backend.Snapshots(context.Background(), name)
			return tabLoadedMsg{tab: tab, name: name, content: out, err: err}
		default:
			return tabLoadedMsg{tab: tab, name: name, content: "Unknown tab"}
		}
	}
}

func toggleStateCmd(backend core.Backend, inst core.Instance) tea.Cmd {
	return func() tea.Msg {
		var err error
		if strings.EqualFold(inst.State, "running") {
			err = backend.Stop(context.Background(), inst.Name)
			return actionDoneMsg{status: "Stopped " + inst.Name, err: err}
		}
		err = backend.Start(context.Background(), inst.Name)
		return actionDoneMsg{status: "Started " + inst.Name, err: err}
	}
}

func deleteCmd(backend core.Backend, name string) tea.Cmd {
	return func() tea.Msg {
		err := backend.Delete(context.Background(), name, false)
		return actionDoneMsg{status: "Deleted " + name, err: err}
	}
}

func runExecCmd(backend core.Backend, name string, cmd []string) tea.Cmd {
	return func() tea.Msg {
		out, err := backend.Exec(context.Background(), name, cmd)
		return execDoneMsg{out: out, err: err}
	}
}

func runLocalLazyCmd(cmdLine string) tea.Cmd {
	return func() tea.Msg {
		exe, err := os.Executable()
		if err != nil {
			return cliCmdDoneMsg{err: err}
		}
		args := strings.Fields(cmdLine)
		cmd := exec.Command(exe, args...)
		out, runErr := cmd.CombinedOutput()
		if runErr != nil {
			return cliCmdDoneMsg{out: string(out), err: runErr}
		}
		return cliCmdDoneMsg{out: string(out)}
	}
}

func shellCmd(name, shell string) tea.Cmd {
	c := exec.Command("lxc", "exec", name, "--", shell)
	return tea.ExecProcess(c, func(err error) tea.Msg {
		return shellDoneMsg{err: err, shell: shell}
	})
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
