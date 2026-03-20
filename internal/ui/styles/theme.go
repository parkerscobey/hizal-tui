package styles

import "github.com/charmbracelet/lipgloss"

// Hizal color palette
var (
	ColorBase      = lipgloss.Color("#1a1a2e")
	ColorSurface   = lipgloss.Color("#16213e")
	ColorOverlay   = lipgloss.Color("#0f3460")
	ColorMuted     = lipgloss.Color("#6272a4")
	ColorSubtle    = lipgloss.Color("#44475a")
	ColorText      = lipgloss.Color("#f8f8f2")
	ColorAccent    = lipgloss.Color("#8be9fd")
	ColorGreen     = lipgloss.Color("#50fa7b")
	ColorYellow    = lipgloss.Color("#f1fa8c")
	ColorOrange    = lipgloss.Color("#ffb86c")
	ColorPink      = lipgloss.Color("#ff79c6")
	ColorPurple    = lipgloss.Color("#bd93f9")
	ColorRed       = lipgloss.Color("#ff5555")
	ColorCyan      = lipgloss.Color("#8be9fd")
)

// Chunk type colors
var ChunkTypeColors = map[string]lipgloss.Color{
	"IDENTITY":       ColorPink,
	"MEMORY":         ColorPurple,
	"KNOWLEDGE":      ColorCyan,
	"CONVENTION":     ColorGreen,
	"PRINCIPLE":      ColorOrange,
	"DECISION":       ColorYellow,
	"RESEARCH":       ColorAccent,
	"PLAN":           ColorAccent,
	"SPEC":           ColorYellow,
	"IMPLEMENTATION": ColorGreen,
	"CONSTRAINT":     ColorRed,
	"LESSON":         ColorOrange,
}

// Layout
var (
	SidebarWidth = 28

	Sidebar = lipgloss.NewStyle().
		Width(SidebarWidth).
		BorderStyle(lipgloss.NormalBorder()).
		BorderRight(true).
		BorderForeground(ColorSubtle).
		Padding(1, 1)

	MainPanel = lipgloss.NewStyle().
		Padding(1, 2)

	StatusBar = lipgloss.NewStyle().
		Background(ColorSurface).
		Foreground(ColorMuted).
		Padding(0, 1).
		Width(0) // set dynamically

	Title = lipgloss.NewStyle().
		Foreground(ColorAccent).
		Bold(true)

	Muted = lipgloss.NewStyle().
		Foreground(ColorMuted)

	Selected = lipgloss.NewStyle().
		Foreground(ColorAccent).
		Bold(true)

	NavItem = lipgloss.NewStyle().
		Foreground(ColorText).
		Padding(0, 1)

	NavItemActive = lipgloss.NewStyle().
		Foreground(ColorAccent).
		Bold(true).
		Padding(0, 1)

	Border = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(ColorSubtle)

	InjectBadge = lipgloss.NewStyle().
		Foreground(ColorGreen).
		Bold(true)
)

// ChunkTypeLabel returns a styled chunk type label
func ChunkTypeLabel(chunkType string) string {
	color, ok := ChunkTypeColors[chunkType]
	if !ok {
		color = ColorMuted
	}
	return lipgloss.NewStyle().
		Foreground(color).
		Bold(true).
		Render(chunkType)
}
