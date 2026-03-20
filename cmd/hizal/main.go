package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"

	"github.com/XferOps/hizal-tui/internal/config"
	"github.com/XferOps/hizal-tui/internal/ui"
)

const version = "0.1.0"

// ─── Setup wizard (first run) ────────────────────────────────────────────────

type setupStep int

const (
	stepAPIURL setupStep = iota
	stepAPIKey
	stepDone
)

type setupModel struct {
	step   setupStep
	url    textinput.Model
	key    textinput.Model
	err    string
}

func newSetupModel() setupModel {
	url := textinput.New()
	url.Placeholder = "https://winnow-api.xferops.dev"
	url.Focus()
	url.Width = 60

	key := textinput.New()
	key.Placeholder = "ctx_..."
	key.EchoMode = textinput.EchoPassword
	key.Width = 60

	return setupModel{step: stepAPIURL, url: url, key: key}
}

var (
	setupTitle  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#bd93f9"))
	setupLabel  = lipgloss.NewStyle().Foreground(lipgloss.Color("#8be9fd"))
	setupMuted  = lipgloss.NewStyle().Foreground(lipgloss.Color("#6272a4"))
	setupError  = lipgloss.NewStyle().Foreground(lipgloss.Color("#ff5555"))
)

func (m setupModel) Init() tea.Cmd { return textinput.Blink }

func (m setupModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "enter":
			if m.step == stepAPIURL {
				if m.url.Value() == "" {
					m.err = "API URL is required"
					return m, nil
				}
				m.err = ""
				m.step = stepAPIKey
				m.url.Blur()
				return m, m.key.Focus()
			}
			if m.step == stepAPIKey {
				if m.key.Value() == "" {
					m.err = "API key is required"
					return m, nil
				}
				cfg := &config.Config{
					APIURL: m.url.Value(),
					APIKey: m.key.Value(),
				}
				if err := config.Save(cfg); err != nil {
					m.err = fmt.Sprintf("Failed to save config: %v", err)
					return m, nil
				}
				m.step = stepDone
				return m, tea.Quit
			}
		}
	}

	var cmd tea.Cmd
	if m.step == stepAPIURL {
		m.url, cmd = m.url.Update(msg)
	} else {
		m.key, cmd = m.key.Update(msg)
	}
	return m, cmd
}

func (m setupModel) View() string {
	s := setupTitle.Render("  hizal — first run setup\n\n")
	s += setupMuted.Render("  Connect to your Hizal API instance.\n\n")

	if m.step == stepAPIURL || m.step == stepDone {
		s += setupLabel.Render("  API URL") + "\n"
		s += "  " + m.url.View() + "\n\n"
	}
	if m.step == stepAPIKey {
		s += setupLabel.Render("  API URL: ") + setupMuted.Render(m.url.Value()) + "\n\n"
		s += setupLabel.Render("  API Key") + "\n"
		s += "  " + m.key.View() + "\n\n"
		s += setupMuted.Render("  Your key starts with ctx_...") + "\n"
	}
	if m.err != "" {
		s += "\n  " + setupError.Render(m.err) + "\n"
	}
	s += "\n  " + setupMuted.Render("Enter to confirm • Ctrl+C to quit")
	return s
}

// ─── Entry point ─────────────────────────────────────────────────────────────

func main() {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "--version", "-v", "version":
			fmt.Printf("hizal %s\n", version)
			os.Exit(0)
		case "--help", "-h", "help":
			fmt.Println("hizal — terminal UI for Hizal structured memory")
			fmt.Printf("version %s\n\n", version)
			fmt.Println("Usage: hizal [flags]")
			fmt.Println()
			fmt.Println("Flags:")
			fmt.Println("  --version    Print version and exit")
			fmt.Println("  --config     Print config path and exit")
			fmt.Println("  --setup      Re-run the setup wizard")
			fmt.Println("  --help       Show this help")
			os.Exit(0)
		case "--config":
			fmt.Println(config.Path())
			os.Exit(0)
		case "--setup":
			runSetup()
			return
		}
	}

	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}

	if !cfg.IsConfigured() {
		cfg = runSetupAndReturn()
		if cfg == nil {
			return
		}
	}

	app := ui.New(cfg)
	p := tea.NewProgram(app, tea.WithAltScreen(), tea.WithMouseCellMotion())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func runSetup() {
	p := tea.NewProgram(newSetupModel())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Setup error: %v\n", err)
		os.Exit(1)
	}
}

func runSetupAndReturn() *config.Config {
	m := newSetupModel()
	p := tea.NewProgram(m)
	result, err := p.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Setup error: %v\n", err)
		return nil
	}
	final, ok := result.(setupModel)
	if !ok || final.step != stepDone {
		return nil
	}
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config after setup: %v\n", err)
		return nil
	}
	return cfg
}
