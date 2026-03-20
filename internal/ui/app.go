package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/XferOps/hizal-tui/internal/api"
	"github.com/XferOps/hizal-tui/internal/config"
	"github.com/XferOps/hizal-tui/internal/ui/styles"
	"github.com/XferOps/hizal-tui/internal/ui/views"
)

type view int

const (
	viewSearch view = iota
	viewSessions
	viewProjects
	viewAgents
)

type navItem struct {
	label string
	view  view
}

var navItems = []navItem{
	{"Search", viewSearch},
	{"Sessions", viewSessions},
	{"Projects", viewProjects},
	{"Agents", viewAgents},
}

type App struct {
	cfg        *config.Config
	client     *api.Client
	activeView view
	search     views.SearchView
	width      int
	height     int
	showHelp   bool
}

func New(cfg *config.Config) *App {
	client := api.New(cfg.APIURL, cfg.APIKey)
	return &App{
		cfg:        cfg,
		client:     client,
		activeView: viewSearch,
		search:     views.NewSearchView(client),
	}
}

func (a *App) Init() tea.Cmd {
	return a.search.Init()
}

func (a *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		a.width = msg.Width
		a.height = msg.Height

	case tea.KeyMsg:
		// Global keybindings
		switch msg.String() {
		case "ctrl+c", "q":
			return a, tea.Quit
		case "?":
			a.showHelp = !a.showHelp
			return a, nil
		case "1":
			a.activeView = viewSearch
		case "2":
			a.activeView = viewSessions
		case "3":
			a.activeView = viewProjects
		case "4":
			a.activeView = viewAgents
		}
	}

	// Route to active view
	switch a.activeView {
	case viewSearch:
		var cmd tea.Cmd
		a.search, cmd = a.search.Update(msg)
		cmds = append(cmds, cmd)
	}

	return a, tea.Batch(cmds...)
}

func (a *App) View() string {
	if a.showHelp {
		return a.helpView()
	}

	sidebar := a.sidebarView()
	main := a.mainView()

	body := lipgloss.JoinHorizontal(lipgloss.Top, sidebar, main)
	status := a.statusBarView()

	return lipgloss.JoinVertical(lipgloss.Left, body, status)
}

func (a *App) sidebarView() string {
	var b strings.Builder

	b.WriteString(styles.Title.Render("hizal") + "\n")
	b.WriteString(styles.Muted.Render("─────────────────────") + "\n\n")

	for _, item := range navItems {
		if item.view == a.activeView {
			b.WriteString(styles.NavItemActive.Render("▸ "+item.label) + "\n")
		} else {
			b.WriteString(styles.NavItem.Render("  "+item.label) + "\n")
		}
	}

	b.WriteString("\n" + styles.Muted.Render("─────────────────────") + "\n")
	b.WriteString(styles.Muted.Render("? help") + "\n")

	return styles.Sidebar.Height(a.height - 1).Render(b.String())
}

func (a *App) mainView() string {
	width := a.width - styles.SidebarWidth - 2
	if width < 0 {
		width = 0
	}

	var content string
	switch a.activeView {
	case viewSearch:
		content = a.search.View()
	case viewSessions:
		content = styles.Muted.Render("Sessions — coming soon")
	case viewProjects:
		content = styles.Muted.Render("Projects — coming soon")
	case viewAgents:
		content = styles.Muted.Render("Agents — coming soon")
	}

	return styles.MainPanel.Width(width).Height(a.height - 1).Render(content)
}

func (a *App) statusBarView() string {
	left := fmt.Sprintf("  %s", a.cfg.APIURL)
	right := "1 Search  2 Sessions  3 Projects  4 Agents  ? Help  q Quit  "
	gap := a.width - len(left) - len(right)
	if gap < 0 {
		gap = 0
	}
	bar := left + strings.Repeat(" ", gap) + right
	return styles.StatusBar.Width(a.width).Render(bar)
}

func (a *App) helpView() string {
	help := styles.Border.Padding(1, 3).Render(strings.Join([]string{
		styles.Title.Render("Hizal TUI — Keybindings"),
		"",
		styles.ChunkTypeLabel("NAVIGATION"),
		"  1-4        Switch views",
		"  j / ↓     Move down",
		"  k / ↑     Move up",
		"  Enter     Open / select",
		"  Esc       Back",
		"",
		styles.ChunkTypeLabel("SEARCH"),
		"  Tab       Cycle scope (All / Project / Agent / Org)",
		"",
		styles.ChunkTypeLabel("GLOBAL"),
		"  ?         Toggle this help",
		"  q         Quit",
	}, "\n"))

	return lipgloss.Place(a.width, a.height, lipgloss.Center, lipgloss.Center, help)
}
