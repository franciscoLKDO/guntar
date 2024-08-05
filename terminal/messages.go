package terminal

import (
	tea "github.com/charmbracelet/bubbletea"
)

type setViewTypeMsg int

const (
	directoryLister setViewTypeMsg = iota
	fileReader      setViewTypeMsg = iota
)

func setView(vt setViewTypeMsg) tea.Cmd {
	return func() tea.Msg {
		return vt
	}
}

type errMsg error

type readDataMsg []byte

func ReadData(p []byte) tea.Cmd {
	return func() tea.Msg {
		return readDataMsg(p)
	}
}
