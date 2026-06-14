package ui

import (
	"context"
	"fmt"
	"github.com/Horryportier/openup/internal/config"
	"github.com/Horryportier/openup/internal/editor"
	"github.com/Horryportier/openup/internal/entry"
	"github.com/Horryportier/openup/internal/favorites"
	"github.com/Horryportier/openup/internal/filesystem"
	"github.com/Horryportier/openup/internal/history"
	"github.com/Horryportier/openup/internal/providers"
	"github.com/Horryportier/openup/internal/search"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"os/exec"
	"path/filepath"
	"strings"
)

type item struct{ entry.Entry }

func (i item) Title() string       { return i.Icon + " " + i.Name }
func (i item) Description() string { return i.Entry.Description }
func (i item) FilterValue() string { return i.Entry.FilterValue() }

type Model struct {
	ctx           context.Context
	cfg           config.Config
	hist          history.Service
	fav           favorites.Service
	list          list.Model
	search        textinput.Model
	searching     bool
	cwd           string
	crumbs        []string
	back, forward []string
	status        string
	width, height int
	all           []entry.Entry
}

var titleStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("212"))
var bar = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
var box = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Padding(0, 1)

func New(ctx context.Context, cfg config.Config, h history.Service, f favorites.Service, start string) (Model, error) {
	if start == "" {
		start = "."
	}
	abs, _ := filepath.Abs(start)
	ti := textinput.New()
	ti.Placeholder = "search everything"
	ti.Prompt = "/ "
	m := Model{ctx: ctx, cfg: cfg, hist: h, fav: f, search: ti, cwd: abs, crumbs: []string{abs}}
	m.list = list.New(nil, list.NewDefaultDelegate(), 0, 0)
	m.list.Title = "OpenedUp"
	return m, m.load(abs)
}
func (m Model) Init() tea.Cmd { return textinput.Blink }
func (m *Model) load(dir string) error {
	ents, err := filesystem.Entries(m.ctx, dir, 500)
	if err != nil {
		return err
	}
	home := []entry.Entry{{ID: "home:favorites", Type: entry.Group, Name: "Favorites", Description: "Pinned and smart favorites", Icon: "★", Action: entry.ActionNavigate, Metadata: map[string]string{"target": "favorites"}}, {ID: "home:history", Type: entry.Group, Name: "History", Description: "Recent files, directories, shell commands", Icon: "◴", Action: entry.ActionNavigate, Metadata: map[string]string{"target": "history"}}, {ID: "home:commands", Type: entry.Group, Name: "Commands", Description: "Configured commands and services", Icon: "⚙", Action: entry.ActionNavigate, Metadata: map[string]string{"target": "commands"}}}
	m.all = append(home, ents...)
	m.set(m.all)
	return nil
}
func (m *Model) set(es []entry.Entry) {
	its := make([]list.Item, len(es))
	for i, e := range es {
		its[i] = item{e}
	}
	m.list.SetItems(its)
}
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.list.SetSize(msg.Width, max(8, msg.Height-7))
	case statusMsg:
		m.status = string(msg)
		return m, nil
	case tea.KeyMsg:
		if m.searching {
			switch msg.String() {
			case "esc":
				m.searching = false
				m.search.Blur()
				m.set(m.all)
				return m, nil
			case "enter":
				m.searching = false
				m.search.Blur()
				return m.openSelected()
			default:
				var cmd tea.Cmd
				m.search, cmd = m.search.Update(msg)
				m.set(search.Rank(m.search.Value(), m.all, m.cfg.SearchLimit))
				return m, cmd
			}
		}
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "ctrl+f":
			m.searching = true
			m.search.Focus()
			return m, textinput.Blink
		case "backspace":
			return m.goBack()
		case "h":
			return m.goBack()
		case "l", "enter":
			return m.openSelected()
		case "ctrl+h":
			m.show("history")
			return m, nil
		case "ctrl+d":
			m.show("favorites")
			return m, nil
		case "ctrl+g":
			m.home()
			return m, nil
		case "ctrl+s":
			m.status = "settings: edit " + config.Path()
			return m, nil
		case "?":
			m.status = "enter open • backspace back • ctrl+f search • ctrl+h history • ctrl+d favorites • ctrl+g home • q quit"
			return m, nil
		}
	}
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}
func (m *Model) home() { m.cwd = "/"; _ = m.load("."); m.crumbs = []string{"home"} }
func (m Model) goBack() (tea.Model, tea.Cmd) {
	if len(m.back) == 0 {
		return m, nil
	}
	prev := m.back[len(m.back)-1]
	m.back = m.back[:len(m.back)-1]
	m.forward = append(m.forward, m.cwd)
	m.cwd = prev
	_ = m.load(prev)
	m.crumbs = append(m.crumbs, prev)
	return m, nil
}
func (m *Model) show(name string) {
	var es []entry.Entry
	switch name {
	case "favorites":
		es = m.fav.Entries(m.ctx)
	case "history":
		es = append(m.hist.Entries(m.ctx), history.ShellCommands(100)...)
	case "commands":
		for _, c := range m.cfg.FavoriteCommands {
			es = append(es, entry.Entry{ID: "cmd:" + c.Name, Type: entry.Command, Name: c.Name, Description: c.Description, Icon: "$", Action: entry.ActionRunCommand, Metadata: map[string]string{"command": c.Command}})
		}
		es = append(es, providers.Systemd(m.ctx, 50)...)
	}
	m.all = es
	m.set(es)
	m.crumbs = append(m.crumbs, name)
}
func (m Model) openSelected() (tea.Model, tea.Cmd) {
	it, ok := m.list.SelectedItem().(item)
	if !ok {
		return m, nil
	}
	e := it.Entry
	_ = m.hist.Record(m.ctx, e)
	switch e.Action {
	case entry.ActionViewFolder:
		p := e.Metadata["path"]
		m.back = append(m.back, m.cwd)
		m.cwd = p
		_ = m.load(p)
		m.crumbs = append(m.crumbs, p)
		return m, nil
	case entry.ActionOpenFile:
		ed := editor.Detect(m.cfg.Editor)
		p := e.Metadata["path"]
		return m, tea.ExecProcess(exec.Command(ed, p), func(error) tea.Msg { return nil })
	case entry.ActionNavigate:
		m.show(e.Metadata["target"])
		return m, nil
	case entry.ActionRunCommand:
		cmd := exec.Command("sh", "-lc", e.Metadata["command"])
		return m, tea.ExecProcess(cmd, func(err error) tea.Msg {
			if err != nil {
				return statusMsg(err.Error())
			}
			return statusMsg("command completed")
		})
	}
	return m, nil
}

type statusMsg string

func (m Model) View() string {
	preview := ""
	if m.cfg.Preview {
		if it, ok := m.list.SelectedItem().(item); ok {
			preview = box.Width(max(20, m.width/3)).Render(filesystem.Preview(it.Metadata["path"], 2000))
		}
	}
	if preview != "" {
		m.list.SetWidth(max(20, m.width-lipgloss.Width(preview)-4))
	}
	header := titleStyle.Render("OpenedUp v2") + " " + bar.Render(strings.Join(m.crumbs, " › "))
	searchLine := ""
	if m.searching {
		searchLine = m.search.View()
	} else {
		searchLine = bar.Render("ctrl+f search • ctrl+h history • ctrl+d favorites • ctrl+g home • ? help • q quit")
	}
	body := m.list.View()
	if preview != "" {
		body = lipgloss.JoinHorizontal(lipgloss.Top, body, "  ", preview)
	}
	return fmt.Sprintf("%s\n%s\n%s\n%s", header, searchLine, body, bar.Render(m.status))
}
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
