package terminal

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTextBox(t *testing.T) {
	tb, err := NewTextBox()

	require.Nil(t, err)
	t.Run("New textbox get default", func(t *testing.T) {
		assert.Equal(t, textboxDefaultWidth, tb.viewport.Width)
		assert.Equal(t, textboxDefaultWHeight, tb.viewport.Height)
		assert.Equal(t, directoryLister, tb.exitView)
	})

	t.Run("Update textbox Size on windowSize cmd", func(t *testing.T) {
		newWidth, newHeigh := 50, 60
		var cmd tea.Cmd
		tb, cmd = tb.Update(tea.WindowSizeMsg{Width: newWidth, Height: newHeigh})
		assert.Nil(t, cmd)
		assert.Equal(t, newWidth, tb.viewport.Width)
		assert.Equal(t, newHeigh-helpHeight, tb.viewport.Height)
	})

	t.Run("Update data on read Data message", func(t *testing.T) {
		assert.Equal(t, 0, tb.viewport.TotalLineCount())
		var cmd tea.Cmd
		tb, cmd = tb.Update(readDataMsg("hello there"))
		require.Nil(t, cmd)
		assert.Greater(t, tb.viewport.TotalLineCount(), 0)
	})

	t.Run("Get view", func(t *testing.T) {
		assert.Contains(t, tb.View(), "hello") // from test [Update data on read Data message]
		assert.Contains(t, tb.View(), tb.helpView())
	})

	t.Run("Ask exitView", func(t *testing.T) {
		var cmd tea.Cmd
		tb, cmd = tb.Update(tea.KeyMsg{Type: tea.KeyBackspace})
		assert.Equal(t, tb.exitView, cmd().(setViewTypeMsg))
	})

	t.Run("Update viewport on unhandled keys", func(t *testing.T) {
		var cmd tea.Cmd
		tb, cmd = tb.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}}) // Type go to bottom button
		assert.Nil(t, cmd)
	})

	t.Run("Send quit on quit key", func(t *testing.T) {
		var cmd tea.Cmd
		tb, cmd = tb.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
		assert.IsType(t, tea.QuitMsg{}, cmd())
	})
}
