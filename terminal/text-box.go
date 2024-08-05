package terminal

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

const helpHeight = 5
const textboxDefaultWidth = 78
const textboxDefaultWHeight = 30

var helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render

type TextBoxModel struct {
	viewport viewport.Model
	renderer *glamour.TermRenderer
	KeyMap   KeyMap
	exitView setViewTypeMsg
}

func NewTextBox() (TextBoxModel, error) {
	vp := viewport.New(textboxDefaultWidth, textboxDefaultWHeight)
	vp.Style = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		PaddingRight(2)

	renderer, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(textboxDefaultWidth),
	)
	if err != nil {
		return TextBoxModel{}, err
	}

	return TextBoxModel{
		viewport: vp,
		renderer: renderer,
		KeyMap:   DefaultKeyMap(),
		exitView: directoryLister,
	}, nil
}

func (t *TextBoxModel) SetSize(msg tea.WindowSizeMsg) {
	t.viewport.Height = msg.Height - helpHeight
	t.viewport.Width = msg.Width
}

func (t *TextBoxModel) renderData(msg readDataMsg) (TextBoxModel, tea.Cmd) {
	buf, err := t.renderer.RenderBytes(msg)
	if err != nil {
		return *t, func() tea.Msg { return errMsg(fmt.Errorf("error on render data: %s", err)) }
	}
	t.viewport.SetContent(string(buf))
	return t.updateViewport(msg)
}

func (t *TextBoxModel) updateViewport(msg tea.Msg) (TextBoxModel, tea.Cmd) {
	var cmd tea.Cmd
	t.viewport, cmd = t.viewport.Update(msg)
	return *t, cmd
}

func (t TextBoxModel) exitViewCmd() (TextBoxModel, tea.Cmd) {
	return t, setView(t.exitView)
}

func (t TextBoxModel) Init() tea.Cmd {
	return nil
}

func (t TextBoxModel) Update(msg tea.Msg) (TextBoxModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		t.SetSize(msg)
	case readDataMsg:
		return t.renderData(msg)
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, t.KeyMap.Quit):
			return t, tea.Quit
		case key.Matches(msg, t.KeyMap.Back):
			return t.exitViewCmd()
		default:
			return t.updateViewport(msg)
		}
	}
	return t, nil
}

func (t TextBoxModel) View() string {
	return t.viewport.View() + t.helpView()
}

func (t TextBoxModel) helpView() string {
	return helpStyle("\n  ↑/↓: Navigate • q: Quit\n")
}
