package terminal

import "github.com/charmbracelet/lipgloss"

const (
	fileSizeWidth = 7
	paddingLeft   = 2
)

// Styles defines the possible customizations for styles in the file picker.
type Styles struct {
	DisabledCursor        lipgloss.Style
	Cursor                lipgloss.Style
	Directory             lipgloss.Style
	File                  lipgloss.Style
	Permission            lipgloss.Style
	CurrentSelected       lipgloss.Style
	SelectedStatus        lipgloss.Style
	PartialSelectedStatus lipgloss.Style
	FileSize              lipgloss.Style
	EmptyDirectory        lipgloss.Style
}

// DefaultStyles defines the default styling for the file picker.
func DefaultStyles() Styles {
	return DefaultStylesWithRenderer(lipgloss.DefaultRenderer())
}

// DefaultStylesWithRenderer defines the default styling for the file picker,
// with a given Lip Gloss renderer.
func DefaultStylesWithRenderer(r *lipgloss.Renderer) Styles {
	return Styles{
		DisabledCursor:        r.NewStyle().Foreground(lipgloss.Color("247")),
		Cursor:                r.NewStyle().Foreground(lipgloss.Color("212")),
		Directory:             r.NewStyle().Foreground(lipgloss.Color("33")).Bold(true),
		File:                  r.NewStyle(),
		Permission:            r.NewStyle().Foreground(lipgloss.Color("244")),
		CurrentSelected:       r.NewStyle().Foreground(lipgloss.Color("212")).Bold(true),
		SelectedStatus:        r.NewStyle().Foreground(lipgloss.Color("42")),
		PartialSelectedStatus: r.NewStyle().Foreground(lipgloss.Color("172")),
		FileSize:              r.NewStyle().Foreground(lipgloss.Color("240")).Width(fileSizeWidth).Align(lipgloss.Right),
		EmptyDirectory:        r.NewStyle().Foreground(lipgloss.Color("240")).PaddingLeft(paddingLeft).SetString("Bummer. No Files Found."),
	}
}

var defaultStyle Styles = DefaultStyles()
