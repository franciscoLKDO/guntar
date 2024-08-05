package terminal

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/franciscolkdo/guntar/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewTerminal(t *testing.T) {
	files := []test.File{
		{Name: "./test/", Mode: 493, Body: ""},
		{Name: "./test/readme.txt", Mode: 0600, Body: "This archive contains some text files."},
		{Name: "./test/hello.txt", Mode: 0600, Body: "world"},
		{Name: "./gopher.txt", Mode: 0600, Body: "Gopher names:\nGeorge\nGeoffrey\nGonzo"},
		{Name: "./todo.txt", Mode: 0600, Body: "Get animal handling license."},
	}
	buf := test.CreateArchive(t, files)
	var term TerminalModel

	update := func(msg tea.Msg) tea.Cmd {
		m, cmd := term.Update(msg)
		term = m.(TerminalModel)
		return cmd
	}
	t.Run("Create new model", func(t *testing.T) {
		var err error
		term, err = New(buf)
		require.Nil(t, err)
		assert.Equal(t, directoryLister, term.CurrentView)
	})

	t.Run("Init model", func(t *testing.T) {
		cmd := term.Init()
		assert.Equal(t, AppName, fmt.Sprintf("%s", cmd()))
	})

	t.Run("Update Models Size on WindowSizeMsg", func(t *testing.T) {
		width, height := 15, 20
		cmd := update(tea.WindowSizeMsg{Width: width, Height: height})
		assert.Nil(t, cmd)
		dl := term.directoryLister
		tb := term.textBox
		assert.Equal(t, height-marginBottom, dl.Height)
		assert.Equal(t, height-marginBottom, tb.viewport.Height)
		assert.Equal(t, width, tb.viewport.Width)
	})

	t.Run("Update directoryLister current node", func(t *testing.T) {
		dl := &term.directoryLister
		cn := dl.currentNode
		cmd := update(readDirNode(cn.GetChildren()[dl.selected])())
		assert.Nil(t, cmd)
		assert.NotEqual(t, dl.currentNode, cn)
	})

	t.Run("Update view to textbox", func(t *testing.T) {
		cmd := update(setViewTypeMsg(fileReader))
		assert.NotNil(t, cmd)
		assert.Equal(t, fileReader, term.CurrentView)
	})

	t.Run("Update textbox body", func(t *testing.T) {
		data := "hello"
		cmd := update(ReadData([]byte(data))())
		assert.Nil(t, cmd)
		assert.Contains(t, term.textBox.View(), data)
	})

	t.Run("Get terminal view with textbox", func(t *testing.T) {
		assert.Contains(t, term.View(), term.textBox.View())
	})

	t.Run("Get terminal view with lister", func(t *testing.T) {
		update(setViewTypeMsg(directoryLister))
		assert.Contains(t, term.View(), term.directoryLister.View())
	})

	t.Run("Get quit message and view", func(t *testing.T) {
		cmd := update(tea.QuitMsg{})
		assert.Equal(t, tea.QuitMsg{}, cmd())
		assert.True(t, term.quitting)
		assert.Equal(t, "Good Bye!", term.View())
	})

	term.quitting = false // reset term quitting status

	t.Run("Get error message and view", func(t *testing.T) {
		err := errors.New("this is a test error")
		cmd := update(errMsg(err))
		assert.Equal(t, tea.QuitMsg{}, cmd())
		assert.True(t, term.quitting)
		assert.Equal(t, err.Error(), term.View())
	})

}

func TestErrorTerminal(t *testing.T) {
	t.Run("Return error on bad tar scanning", func(t *testing.T) {
		term, err := New(strings.NewReader("hello, this is not an archive reader"))
		assert.ErrorContains(t, err, "error on scanning tar file")
		assert.Equal(t, TerminalModel{}, term)
	})
}
