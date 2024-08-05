package terminal

import (
	"fmt"
	"io"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/franciscolkdo/guntar/tar"
)

const AppName = "Guntar"

type TerminalModel struct {
	textBox         TextBoxModel
	directoryLister ListerModel
	CurrentView     setViewTypeMsg
	KeyMap          KeyMap
	quitting        bool
	err             error
}

func New(tarFile io.Reader) (TerminalModel, error) {
	tb, _ := NewTextBox()

	root, err := tar.Scan(tarFile, OnNewNode)
	if err != nil {
		return TerminalModel{}, fmt.Errorf("error on scanning tar file: %s", err)
	}
	return TerminalModel{
		textBox:         tb,
		directoryLister: NewLister(root),
		CurrentView:     directoryLister,
		KeyMap:          DefaultKeyMap(),
		quitting:        false,
		err:             nil,
	}, nil
}

// Init initializes the file picker model.
func (m TerminalModel) Init() tea.Cmd {
	return tea.SetWindowTitle(AppName)
}

func (m TerminalModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.QuitMsg:
		m.quitting = true
		return m, tea.Quit
	case errMsg:
		m.err = msg
		m.quitting = true
		return m, tea.Quit
	case tea.WindowSizeMsg:
		m.directoryLister.SetSize(msg)
		m.textBox.SetSize(msg)
		return m, nil

	case setViewTypeMsg:
		m.CurrentView = msg
		if m.CurrentView == fileReader {
			return m, ReadData(m.directoryLister.GetSelectedFile().GetData())
		}
	}

	var cmd tea.Cmd
	switch m.CurrentView {
	case directoryLister:
		m.directoryLister, cmd = m.directoryLister.Update(msg)
		return m, cmd
	case fileReader:
		m.textBox, cmd = m.textBox.Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m TerminalModel) View() string {
	if m.quitting {
		if m.err != nil {
			return fmt.Sprint(m.err)
		}
		return "Good Bye!"
	}

	var s string
	switch m.CurrentView {
	case directoryLister:
		s = m.directoryLister.View()
	case fileReader:
		s = m.textBox.View()
	}
	return s
}
