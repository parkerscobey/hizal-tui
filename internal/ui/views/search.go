package views

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/XferOps/hizal-tui/internal/api"
	"github.com/XferOps/hizal-tui/internal/ui/styles"
)

type SearchResultsMsg struct {
	Chunks []api.Chunk
	Err    error
}

type SearchView struct {
	input    textinput.Model
	spinner  spinner.Model
	chunks   []api.Chunk
	cursor   int
	loading  bool
	err      error
	client   *api.Client
	scope    string
	scopes   []string
	scopeIdx int
	width    int
	height   int
}

func NewSearchView(client *api.Client) SearchView {
	ti := textinput.New()
	ti.Placeholder = "Search chunks..."
	ti.Focus()
	ti.CharLimit = 200

	sp := spinner.New()
	sp.Spinner = spinner.Dot
	sp.Style = lipgloss.NewStyle().Foreground(styles.ColorAccent)

	return SearchView{
		input:   ti,
		spinner: sp,
		client:  client,
		scopes:  []string{"all", "PROJECT", "AGENT", "ORG"},
		scope:   "all",
	}
}

func (v SearchView) Init() tea.Cmd {
	return textinput.Blink
}

type searchTickMsg struct{ query string }

func doSearch(client *api.Client, query, scope string) tea.Cmd {
	return func() tea.Msg {
		// Debounce: small sleep so rapid keystrokes collapse
		time.Sleep(300 * time.Millisecond)
		chunks, err := client.SearchChunks(query, scope)
		return SearchResultsMsg{Chunks: chunks, Err: err}
	}
}

func (v SearchView) Update(msg tea.Msg) (SearchView, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		v.width = msg.Width
		v.height = msg.Height

	case tea.KeyMsg:
		switch msg.String() {
		case "tab":
			v.scopeIdx = (v.scopeIdx + 1) % len(v.scopes)
			v.scope = v.scopes[v.scopeIdx]
			if v.input.Value() != "" {
				v.loading = true
				cmds = append(cmds, doSearch(v.client, v.input.Value(), v.scope))
			}
		case "up", "k":
			if v.cursor > 0 {
				v.cursor--
			}
		case "down", "j":
			if v.cursor < len(v.chunks)-1 {
				v.cursor++
			}
		}

	case SearchResultsMsg:
		v.loading = false
		v.err = msg.Err
		v.chunks = msg.Chunks
		v.cursor = 0

	case spinner.TickMsg:
		if v.loading {
			var cmd tea.Cmd
			v.spinner, cmd = v.spinner.Update(msg)
			cmds = append(cmds, cmd)
		}
	}

	prevQuery := v.input.Value()
	var inputCmd tea.Cmd
	v.input, inputCmd = v.input.Update(msg)
	cmds = append(cmds, inputCmd)

	// Trigger search on query change
	if v.input.Value() != prevQuery && v.input.Value() != "" {
		v.loading = true
		cmds = append(cmds, v.spinner.Tick, doSearch(v.client, v.input.Value(), v.scope))
	}

	return v, tea.Batch(cmds...)
}

func (v SearchView) SelectedChunk() *api.Chunk {
	if len(v.chunks) == 0 || v.cursor >= len(v.chunks) {
		return nil
	}
	return &v.chunks[v.cursor]
}

func (v SearchView) View() string {
	var b strings.Builder

	// Search bar
	scopeLabel := styles.Muted.Render(fmt.Sprintf("[%s]", v.scope))
	searchRow := lipgloss.JoinHorizontal(lipgloss.Top,
		v.input.View(),
		"  ",
		scopeLabel,
		"  Tab to cycle scope",
	)
	b.WriteString(styles.Border.Width(v.width - styles.SidebarWidth - 6).Render(searchRow))
	b.WriteString("\n\n")

	// Results
	if v.loading {
		b.WriteString(fmt.Sprintf("  %s Searching...\n", v.spinner.View()))
		return b.String()
	}

	if v.err != nil {
		b.WriteString(styles.Muted.Render(fmt.Sprintf("  Error: %s", v.err)))
		return b.String()
	}

	if v.input.Value() == "" {
		b.WriteString(styles.Muted.Render("  Start typing to search your context..."))
		return b.String()
	}

	if len(v.chunks) == 0 {
		b.WriteString(styles.Muted.Render("  No chunks found."))
		return b.String()
	}

	for i, chunk := range v.chunks {
		cursor := "  "
		titleStyle := lipgloss.NewStyle().Foreground(styles.ColorText)
		if i == v.cursor {
			cursor = styles.Selected.Render("▸ ")
			titleStyle = styles.Selected
		}

		typeLabel := styles.ChunkTypeLabel(chunk.ChunkType)

		inject := ""
		if chunk.InjectAudience != nil {
			inject = styles.InjectBadge.Render(" ↓")
		}

		title := chunk.QueryKey
		if title == "" {
			// Truncate content as title
			content := chunk.Content
			if len(content) > 60 {
				content = content[:60] + "..."
			}
			title = content
		}

		row := fmt.Sprintf("%s%s  %s%s",
			cursor,
			typeLabel,
			titleStyle.Render(title),
			inject,
		)
		b.WriteString(row + "\n")
	}

	b.WriteString(styles.Muted.Render(fmt.Sprintf("\n  %d results  ·  Enter to open  ·  Tab to change scope", len(v.chunks))))
	return b.String()
}
