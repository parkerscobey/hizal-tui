package views

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/XferOps/hizal-tui/internal/api"
	"github.com/XferOps/hizal-tui/internal/ui/styles"
)

type ChunkDetailLoadedMsg struct {
	Detail *api.ChunkDetail
	Err    error
}

type ChunkDeletedMsg struct {
	ID string
}

type BackMsg struct{}

type DetailView struct {
	spinner   spinner.Model
	detail    *api.ChunkDetail
	loading   bool
	err       error
	confirmDelete bool
	client    *api.Client
	chunk     api.Chunk
	width     int
	height    int
	scrollOff int
}

func NewDetailView(client *api.Client, chunk api.Chunk) DetailView {
	sp := spinner.New()
	sp.Spinner = spinner.Dot
	sp.Style = lipgloss.NewStyle().Foreground(styles.ColorAccent)

	return DetailView{
		spinner: sp,
		client:  client,
		chunk:   chunk,
		loading: true,
	}
}

func (v DetailView) Init() tea.Cmd {
	return tea.Batch(v.spinner.Tick, func() tea.Msg {
		detail, err := v.client.GetChunk(v.chunk.ID)
		return ChunkDetailLoadedMsg{Detail: detail, Err: err}
	})
}

func doGetChunk(client *api.Client, id string) tea.Cmd {
	return func() tea.Msg {
		detail, err := client.GetChunk(id)
		return ChunkDetailLoadedMsg{Detail: detail, Err: err}
	}
}

func doDeleteChunk(client *api.Client, id string) tea.Cmd {
	return func() tea.Msg {
		if err := client.DeleteChunk(id); err != nil {
			return ChunkDeletedMsg{ID: ""}
		}
		return ChunkDeletedMsg{ID: id}
	}
}

func (v DetailView) Update(msg tea.Msg) (DetailView, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		v.width = msg.Width
		v.height = msg.Height

	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return v, func() tea.Msg { return BackMsg{} }

		case "d":
			if !v.confirmDelete {
				v.confirmDelete = true
				return v, nil
			}

		case "y":
			if v.confirmDelete {
				return v, doDeleteChunk(v.client, v.chunk.ID)
			}

		case "n":
			if v.confirmDelete {
				v.confirmDelete = false
				return v, nil
			}

		case "up", "k":
			if v.scrollOff > 0 {
				v.scrollOff--
			}

		case "down", "j":
			v.scrollOff++
		}

	case ChunkDetailLoadedMsg:
		v.loading = false
		v.err = msg.Err
		v.detail = msg.Detail

	case spinner.TickMsg:
		if v.loading {
			var cmd tea.Cmd
			v.spinner, cmd = v.spinner.Update(msg)
			cmds = append(cmds, cmd)
		}
	}

	return v, tea.Batch(cmds...)
}

func (v DetailView) View() string {
	var b strings.Builder

	typeLabel := styles.ChunkTypeLabel(v.chunk.ChunkType)
	header := fmt.Sprintf("%s  %s  %s", typeLabel, v.chunk.QueryKey, styles.Muted.Render("[d]elete"))
	b.WriteString(styles.Border.Width(v.width - styles.SidebarWidth - 6).Render(header))
	b.WriteString("\n")

	if v.confirmDelete {
		b.WriteString(styles.Selected.Render("  Delete this chunk? [y/N]"))
		b.WriteString("\n")
		return b.String()
	}

	if v.loading {
		b.WriteString(fmt.Sprintf("  %s Loading...\n", v.spinner.View()))
		return b.String()
	}

	if v.err != nil {
		b.WriteString(styles.Muted.Render(fmt.Sprintf("  Error: %s", v.err)))
		b.WriteString("\n")
		return b.String()
	}

	display := v.chunk
	if v.detail != nil {
		display = v.detail.Chunk
	}

	scope := display.Scope
	if scope == "" {
		scope = "—"
	}
	chunkType := display.ChunkType
	if chunkType == "" {
		chunkType = "—"
	}

	autoInject := "off"
	if display.InjectAudience != nil {
		if len(display.InjectAudience.Rules) > 0 {
			rule := display.InjectAudience.Rules[0]
			if rule.All {
				autoInject = "all agents"
			} else {
				var parts []string
				if len(rule.LifecycleTypes) > 0 {
					parts = append(parts, rule.LifecycleTypes...)
				}
				if len(rule.AgentTypes) > 0 {
					parts = append(parts, rule.AgentTypes...)
				}
				if len(parts) > 0 {
					autoInject = strings.Join(parts, ", ")
				}
			}
		}
	}

	metaLine := fmt.Sprintf("Scope: %s  |  Type: %s  |  Auto-inject: %s",
		styles.Muted.Render(scope),
		styles.Muted.Render(chunkType),
		styles.InjectBadge.Render(autoInject),
	)
	b.WriteString(styles.Muted.Render("  " + metaLine))
	b.WriteString("\n\n")

	contentWidth := v.width - styles.SidebarWidth - 10
	if contentWidth < 1 {
		contentWidth = 80
	}

	lines := strings.Split(display.Content, "\n")
	for i := v.scrollOff; i < len(lines); i++ {
		line := lines[i]
		if len(line) > contentWidth {
			line = line[:contentWidth]
		}
		b.WriteString("  " + line + "\n")
	}

	if v.detail != nil && len(v.detail.Versions) > 0 {
		b.WriteString("\n")
		b.WriteString(styles.Muted.Render("  ── Version History ──"))
		b.WriteString("\n")
		for i, ver := range v.detail.Versions {
			isCurrent := ""
			if i == 0 {
				isCurrent = " " + styles.Selected.Render("(current)")
			}
			date := ver.CreatedAt.Format("Jan 02, 2006")
			b.WriteString(fmt.Sprintf("  v%d  %s%s", ver.Version, styles.Muted.Render(date), isCurrent))
			if ver.ChangeNote != "" {
				b.WriteString(" — " + ver.ChangeNote)
			}
			b.WriteString("\n")
		}
	}

	b.WriteString(styles.Muted.Render("  j/k scroll  ·  Esc back"))
	return b.String()
}
