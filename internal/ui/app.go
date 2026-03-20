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

type panel int

const (
	panelSidebar panel = iota
	panelMain
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
	cfg           *config.Config
	client        *api.Client
	activeView    view
	focusedPanel  panel
	sidebarCursor int
	search        views.SearchView
	width         int
	height        int
	showHelp      bool
}

func New(cfg *config.Config) *App {
	client := api.New(cfg.APIURL, cfg.APIKey)
	return &App{
		cfg:           cfg,
		client:        client,
		activeView:    viewSearch,
		focusedPanel:  panelSidebar,
		sidebarCursor: 0,
		search:        views.NewSearchView(client),
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
		case "tab":
			if a.focusedPanel == panelSidebar {
				a.focusedPanel = panelMain
			} else {
				a.focusedPanel = panelSidebar
			}
			return a, nil
		case "shift+tab":
			if a.focusedPanel == panelSidebar {
				a.focusedPanel = panelMain
			} else {
				a.focusedPanel = panelSidebar
			}
			return a, nil
		case "esc":
			if a.focusedPanel == panelMain {
				a.focusedPanel = panelSidebar
			}
			return a, nil
		case "1":
			a.activeView = viewSearch
			a.sidebarCursor = 0
		case "2":
			a.activeView = viewSessions
			a.sidebarCursor = 1
		case "3":
			a.activeView = viewProjects
			a.sidebarCursor = 2
		case "4":
			a.activeView = viewAgents
			a.sidebarCursor = 3
		}

		// Panel-specific keybindings
		if a.focusedPanel == panelSidebar {
			switch msg.String() {
			case "j", "down":
				if a.sidebarCursor < len(navItems)-1 {
					a.sidebarCursor++
					a.activeView = navItems[a.sidebarCursor].view
				}
				return a, nil
			case "k", "up":
				if a.sidebarCursor > 0 {
					a.sidebarCursor--
					a.activeView = navItems[a.sidebarCursor].view
				}
				return a, nil
			case "enter", " ":
				a.focusedPanel = panelMain
				return a, nil
			}
		}
	}

	// Route to active view (only when main panel is focused)
	if a.focusedPanel == panelMain {
		switch a.activeView {
		case viewSearch:
			var cmd tea.Cmd
			a.search, cmd = a.search.Update(msg)
			cmds = append(cmds, cmd)
		}
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

	for i, item := range navItems {
		isCursor := i == a.sidebarCursor
		isActive := a.focusedPanel == panelSidebar && isCursor
		if isActive {
			b.WriteString(styles.NavItemActive.Render("▸ "+item.label) + "\n")
		} else if isCursor {
			b.WriteString(styles.NavItem.Render("  "+item.label) + "\n")
		} else {
			b.WriteString(styles.NavItem.Render("  "+item.label) + "\n")
		}
	}

	b.WriteString("\n" + styles.Muted.Render("─────────────────────") + "\n")
	focusHint := "Tab to main  ·  ? help"
	if a.focusedPanel == panelMain {
		focusHint = "Tab to sidebar  ·  Esc back  ·  ? help"
	}
	b.WriteString(styles.Muted.Render(focusHint) + "\n")

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
	panelIndicator := "SIDEBAR"
	if a.focusedPanel == panelMain {
		panelIndicator = "MAIN"
	}
	left := fmt.Sprintf("  %s  [%s]", a.cfg.APIURL, panelIndicator)
	right := "Tab focus  j/k nav  Enter select  Esc back  ? help  q quit  "
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
		"  Tab / Shift+Tab   Switch panel focus",
		"  1-4               Switch view directly",
		"  j / k / ↑ / ↓    Navigate lists (sidebar or main)",
		"  Enter / Space     Select / open",
		"  Esc               Back to sidebar",
		"",
		styles.ChunkTypeLabel("SEARCH (in main panel)"),
		"  Tab               Cycle scope (All / Project / Agent / Org)",
		"",
		styles.ChunkTypeLabel("GLOBAL"),
		"  ?                 Toggle this help",
		"  q / Ctrl+C        Quit",
	}, "\n"))

	return lipgloss.Place(a.width, a.height, lipgloss.Center, lipgloss.Center, help)
}
