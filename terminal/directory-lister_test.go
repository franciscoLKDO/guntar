package terminal

import (
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/franciscolkdo/guntar/tar"
	"github.com/franciscolkdo/guntar/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewLister(t *testing.T) {
	files := []test.File{
		{Name: "./test/", Mode: 493, Body: ""},
		{Name: "./test/readme.txt", Mode: 0600, Body: "This archive contains some text files."},
		{Name: "./test/hello.txt", Mode: 0600, Body: "world"},
		{Name: "./gopher.txt", Mode: 0600, Body: "Gopher names:\nGeorge\nGeoffrey\nGonzo"},
		{Name: "./todo.txt", Mode: 0600, Body: "Get animal handling license."},
	}
	buf := test.CreateArchive(t, files)
	root, err := tar.Scan(buf, OnNewNode)
	require.Nil(t, err)
	l := NewLister(root, "")
	t.Run("New Lister get default", func(t *testing.T) {
		assert.Equal(t, fileReader, l.enterFileView)
		assert.True(t, l.currentNode.IsDir())
		assert.True(t, l.currentNode.IsRoot())
		assert.Equal(t, l.currentNode.LenChildren(), 3) // test,goopher.txt,todo.txt
	})

	t.Run("Update Lister Size on windowSize cmd", func(t *testing.T) {
		newHeigh := 15
		var cmd tea.Cmd
		l, cmd = l.Update(tea.WindowSizeMsg{Height: newHeigh})
		assert.Nil(t, cmd)
		assert.Equal(t, newHeigh-marginBottom, l.Height)
	})

	t.Run("Go inside directory on key enter (directory selected)", func(t *testing.T) {
		var cmd tea.Cmd
		l, cmd = l.Update(tea.KeyMsg{Type: tea.KeyEnter}) // Run open function on KeyEnter
		assert.IsType(t, DirMsg{}, cmd())                 // New command with DirMsg is created
		l, cmd = l.Update(cmd())                          // Update lister with provided command
		assert.Nil(t, cmd)
		assert.True(t, l.currentNode.IsDir())   // Current Node is ./test
		assert.False(t, l.currentNode.IsRoot()) // Current Node is not root
		assert.Equal(t, "/test", l.currentNode.GetPath())
	})

	t.Run("Ask enterFileView on key enter (file selected)", func(t *testing.T) {
		var cmd tea.Cmd
		l, cmd = l.Update(tea.KeyMsg{Type: tea.KeyEnter}) // Run open function on keyEnter
		assert.Equal(t, l.enterFileView, cmd().(setViewTypeMsg))
	})

	t.Run("Go backward on key backspace", func(t *testing.T) {
		var cmd tea.Cmd
		l, cmd = l.Update(tea.KeyMsg{Type: tea.KeyBackspace}) // Run back function on KeyBackspace
		assert.IsType(t, DirMsg{}, cmd())                     // New command with DirMsg is created
		l, cmd = l.Update(cmd())                              // Update lister with provided command
		assert.Nil(t, cmd)
		assert.True(t, l.currentNode.IsDir())  // Current Node is .
		assert.True(t, l.currentNode.IsRoot()) // Current Node is root again
	})

	type test struct {
		name             string
		key              []rune
		expectedPosition int
	}

	// This is a sequencial test, can fail if order changed
	positionTests := []test{
		{
			name:             "Go to last children",
			key:              []rune{'G'},
			expectedPosition: l.currentNode.LenChildren() - 1,
		},
		{
			name:             "Go to first children",
			key:              []rune{'g'},
			expectedPosition: 0,
		},
		{
			name:             "Go down children",
			key:              []rune{'j'},
			expectedPosition: 1,
		},
		{
			name:             "Go up children",
			key:              []rune{'k'},
			expectedPosition: 0,
		},
		{
			name:             "Go to page down children",
			key:              []rune{'J'},
			expectedPosition: 2,
		},
		{
			name:             "Go to page up children",
			key:              []rune{'K'},
			expectedPosition: 0,
		},
	}
	for _, tt := range positionTests {
		t.Run(tt.name, func(t *testing.T) {
			var cmd tea.Cmd
			l, cmd = l.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: tt.key}) // Type go to bottom button
			assert.Nil(t, cmd)
			assert.True(t, l.currentNode.IsRoot())           // Current Node is still root
			assert.Equal(t, tt.expectedPosition, l.selected) // Selected node is the last children of root node
		})
	}

	t.Run("Send quit on quit key", func(t *testing.T) {
		var cmd tea.Cmd
		l, cmd = l.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
		assert.IsType(t, tea.QuitMsg{}, cmd())
	})
}

func TestListerSelection(t *testing.T) {
	files := []test.File{
		{Name: "./test/", Mode: fs.ModeDir, Body: ""},
		{Name: "./test/nested/", Mode: fs.ModeDir, Body: ""},
		{Name: "./test/nested/.secret", Mode: 0600, Body: "This is not so secret"},
		{Name: "./test/readme.txt", Mode: 0600, Body: "This archive contains some text files."},
		{Name: "./blog/", Mode: fs.ModeDir, Body: ""},
		{Name: "./blog/users.txt", Mode: 0600, Body: "thefumist blaireau_furtif xiunny alejandro coureur_sans_jupon franciscolkdo"},
		{Name: "./gopher.txt", Mode: 0600, Body: "Gopher names:\nGeorge\nGeoffrey\nGonzo"},
		{Name: "./todo.txt", Mode: 0600, Body: "Get animal handling license."},
	}
	buf := test.CreateArchive(t, files)
	root, err := tar.Scan(buf, OnNewNode)
	require.Nil(t, err)
	l := NewLister(root, "")
	l.SetSize(tea.WindowSizeMsg{Height: 10})
	t.Run("Select './test/' directory node", func(t *testing.T) {
		var cmd tea.Cmd
		l, cmd = l.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}}) // Type select button
		assert.Nil(t, cmd)
		nd := l.currentNode.GetChildren()[l.selected]
		assert.Equal(t, Selected, nd.Spec.selectionStatus)
		for _, n := range nd.GetChildren() {
			assert.Equal(t, Selected, n.Spec.selectionStatus)
		}
	})

	t.Run("DeSelect './test/nested/' directory node", func(t *testing.T) {
		var cmd tea.Cmd
		l, cmd = l.Update(tea.KeyMsg{Type: tea.KeyEnter}) // Run open function on KeyEnter
		assert.IsType(t, DirMsg{}, cmd())                 // New command with DirMsg is created
		l, cmd = l.Update(cmd())                          // Update lister with provided command
		assert.Nil(t, cmd)
		l, cmd = l.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}}) // Type select button (deselect on selected node)
		assert.Nil(t, cmd)
		nd := l.currentNode.GetChildren()[l.selected]
		assert.Equal(t, "/test/nested", nd.GetPath())
		assert.Equal(t, NotSelected, nd.Spec.selectionStatus) // Nested directory is deselected
		for _, n := range nd.GetChildren() {
			assert.Equal(t, NotSelected, n.Spec.selectionStatus)
		}
		assert.Equal(t, PartialSelected, l.currentNode.GetParent().Spec.selectionStatus) // Test node is now partial selected
	})

	t.Run("Export selected files", func(t *testing.T) {
		var cmd tea.Cmd
		l, cmd = l.Update(tea.KeyMsg{Type: tea.KeyBackspace}) // Run back function on KeyBackspace
		assert.IsType(t, DirMsg{}, cmd())                     // New command with DirMsg is created
		l, cmd = l.Update(cmd())                              // Update lister with provided command
		assert.Nil(t, cmd)
		tmpDir, err := os.MkdirTemp("", "*")
		require.Nil(t, err)
		defer os.RemoveAll(tmpDir)
		l.exportPath = tmpDir
		// Go in /test

		l, cmd = l.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}}) // Type select button (deselect on selected node)
		assert.Nil(t, cmd)
		testNode := l.currentNode.GetChildren()[l.selected]
		// test directory should be found
		assert.DirExists(t, filepath.Join(tmpDir, testNode.GetPath()))
		// test/readme.txt file should be found
		assert.FileExists(t, filepath.Join(tmpDir, testNode.GetChildren()[1].GetPath()))
	})

	t.Run("View lister", func(t *testing.T) {
		s := l.View()
		for _, nd := range l.currentNode.GetChildren() {
			assert.Contains(t, s, nd.Name())
		}
	})
}
