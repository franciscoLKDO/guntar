package terminal

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/franciscolkdo/guntar/tar"
)

const marginBottom = 5

const checkMark = "✓"
const cursor = ">"

// ListerModel represents a file picker.
type ListerModel struct {
	exportPath      string
	KeyMap          KeyMap
	currentNode     *listerNode
	ShowPermissions bool
	ShowSize        bool
	selected        int
	selectedStack   stack
	min             int
	max             int
	maxStack        stack
	minStack        stack
	Height          int
	flatten         bool
	enterFileView   setViewTypeMsg
}

// NewLister return a Node lister with default styling and key bindings.
func NewLister(n *listerNode) ListerModel {
	pwd, err := os.Getwd()
	if err != nil {
		tea.Printf("%s", err)
	}

	return ListerModel{
		exportPath:      filepath.Join(pwd, "extracted"),
		selected:        0,
		currentNode:     n,
		ShowPermissions: true,
		ShowSize:        true,
		Height:          0,
		max:             0,
		min:             0,
		selectedStack:   newStack(),
		minStack:        newStack(),
		maxStack:        newStack(),
		KeyMap:          DefaultKeyMap(),
		flatten:         false,
		enterFileView:   fileReader,
	}
}

type stack struct {
	Push   func(int)
	Pop    func() int
	Length func() int
}

func newStack() stack {
	slice := make([]int, 0)
	return stack{
		Push: func(i int) {
			slice = append(slice, i)
		},
		Pop: func() int {
			res := slice[len(slice)-1]
			slice = slice[:len(slice)-1]
			return res
		},
		Length: func() int {
			return len(slice)
		},
	}
}

func (m *ListerModel) pushView(selected, min, max int) {
	m.selectedStack.Push(selected)
	m.minStack.Push(min)
	m.maxStack.Push(max)
}

func (m *ListerModel) popView() (int, int, int) {
	return m.selectedStack.Pop(), m.minStack.Pop(), m.maxStack.Pop()
}

type DirMsg struct {
	node *listerNode
}

func readDirNode(n *listerNode) tea.Cmd {
	return func() tea.Msg {
		return DirMsg{node: n}
	}
}

func (m ListerModel) GetSelectedFile() *listerNode {
	return m.currentNode.GetChildren()[m.selected]
}

func (m *ListerModel) up() {
	m.selected--
	if m.selected < 0 {
		m.selected = 0
	}
	if m.selected < m.min {
		m.min--
		m.max--
	}
}

func (m *ListerModel) down() {
	m.selected++
	if m.selected >= m.currentNode.LenChildren() {
		m.selected = m.currentNode.LenChildren() - 1
	}
	if m.selected > m.max {
		m.min++
		m.max++
	}
}

func (m ListerModel) open() (ListerModel, tea.Cmd) {
	if m.currentNode.LenChildren() == 0 {
		return m, nil
	}

	f := m.currentNode.GetChildren()[m.selected]

	if f.IsDir() {
		m.pushView(m.selected, m.min, m.max)
		m.selected = 0
		m.min = 0
		m.max = m.Height - 1
		return m, readDirNode(f)
	} else if f.Mode().IsRegular() {
		return m, setView(m.enterFileView)
	}
	return m, nil
}

func (m ListerModel) back() (ListerModel, tea.Cmd) {
	if m.selectedStack.Length() > 0 {
		m.selected, m.min, m.max = m.popView()
	} else {
		m.selected = 0
		m.min = 0
		m.max = m.Height - 1
	}

	return m, readDirNode(m.currentNode.GetParent())
}

// setSelectionNode run top to bot to select or not all children from current node
func setSelectionNode(node *listerNode, sel SelectedState) {
	node.Spec.selectionStatus = sel
	node.ForAllChildren(func(n *listerNode) error {
		n.Spec.selectionStatus = sel
		return nil
	})
}

// setSelectionParentNode run bot to top to set all parents from current node as partially selected
func setSelectionParentNode(node *listerNode) {
	if node.IsRoot() {
		return
	}

	p := node.GetParent()
	p.Spec.selectionStatus = PartialSelected
	setSelectionParentNode(p)
}

func (m *ListerModel) extract(node *listerNode) tea.Cmd {
	if err := tar.Extract(node, m.exportPath, func(n *listerNode) bool {
		return n.Spec.selectionStatus == NotSelected
	}); err != nil {
		return func() tea.Msg { return errMsg(err) }
	}
	return nil
}

func (m *ListerModel) SetSize(msg tea.WindowSizeMsg) {
	m.Height = msg.Height - marginBottom
	m.max = m.Height - 1
}

// Init initializes the file picker model.
func (m ListerModel) Init() tea.Cmd {
	return nil
}

// Update handles user interactions within the file picker model.
func (m ListerModel) Update(msg tea.Msg) (ListerModel, tea.Cmd) {
	switch msg := msg.(type) {
	case DirMsg:
		m.currentNode = msg.node
		m.max = max(m.max, m.Height-1)
	case tea.WindowSizeMsg:
		m.SetSize(msg)
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.KeyMap.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.KeyMap.GoToTop):
			m.selected = 0
			m.min = 0
			m.max = m.Height - 1
		case key.Matches(msg, m.KeyMap.GoToLast):
			m.selected = m.currentNode.LenChildren() - 1
			m.min = m.currentNode.LenChildren() - m.Height
			m.max = m.currentNode.LenChildren() - 1
		case key.Matches(msg, m.KeyMap.Down):
			m.down()
		case key.Matches(msg, m.KeyMap.Up):
			m.up()
		case key.Matches(msg, m.KeyMap.PageDown):
			m.selected += m.Height
			if m.selected >= m.currentNode.LenChildren() {
				m.selected = m.currentNode.LenChildren() - 1
			}
			m.min += m.Height
			m.max += m.Height

			if m.max >= m.currentNode.LenChildren() {
				m.max = m.currentNode.LenChildren() - 1
				m.min = m.max - m.Height
			}
		case key.Matches(msg, m.KeyMap.PageUp):
			m.selected -= m.Height
			if m.selected < 0 {
				m.selected = 0
			}
			m.min -= m.Height
			m.max -= m.Height

			if m.min < 0 {
				m.min = 0
				m.max = m.min + m.Height
			}
		case key.Matches(msg, m.KeyMap.Back):
			return m.back()
		case key.Matches(msg, m.KeyMap.Open):
			return m.open()
		case key.Matches(msg, m.KeyMap.Select):
			sf := m.currentNode.GetChildren()[m.selected]
			if getSelectionStatus(*sf) == NotSelected {
				setSelectionNode(sf, Selected)
			} else {
				setSelectionNode(sf, NotSelected)
			}
			setSelectionParentNode(sf)

		case key.Matches(msg, m.KeyMap.Extract):
			return m, m.extract(m.currentNode.GetRoot())
		}
	case tea.MouseMsg:
		if msg.Action == tea.MouseActionPress {
			switch msg.Button {
			case tea.MouseButtonWheelDown:
				m.down()
			case tea.MouseButtonWheelUp:
				m.up()
			}
		}
	}
	return m, nil
}

// View returns the view of the file picker.
func (m ListerModel) View() string {
	var s strings.Builder
	s.WriteString(fmt.Sprintf("[%s]\n", m.currentNode.GetPath()))

	if m.currentNode.LenChildren() == 0 {
		return defaultStyle.EmptyDirectory.Height(m.Height).MaxHeight(m.Height).String()
	}

	for i, n := range m.currentNode.GetChildren() {
		if i < m.min || i > m.max {
			continue
		}
		prefix := " "
		if getSelectionStatus(*n) != NotSelected {
			prefix = checkMark
		}
		style := getSelectionStyle(*n)
		if m.selected == i {
			prefix = cursor
		}

		// Render line
		s.WriteString(style.Render(prefix, formatNode(*n)))
		s.WriteRune('\n')
	}

	for i := lipgloss.Height(s.String()); i <= m.Height; i++ {
		s.WriteRune('\n')
	}
	return s.String()
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
