package terminal

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/dustin/go-humanize"
	"github.com/franciscolkdo/guntar/tar"
)

type SelectedState int

const (
	NotSelected SelectedState = iota
	PartialSelected
	Selected
)

type listerData struct {
	selectionStatus SelectedState
	style           lipgloss.Style
}

// listerNode is an alias to Node[listerData]
type listerNode = tar.Node[listerData]

// Callback function for tar.Scan on node creation.
// This add the style and selection state on creation based on FileInfo
func OnNewNode(n *listerNode) error {
	n.Spec.selectionStatus = NotSelected
	n.Spec.style = defaultStyle.File
	if n.IsDir() {
		n.Spec.style = defaultStyle.Directory
	}
	return nil
}

func formatNode(n listerNode) string {
	// Add file mode
	line := " " + defaultStyle.Permission.Render(n.Mode().String())
	// Add file size
	line += fmt.Sprintf("%"+strconv.Itoa(defaultStyle.FileSize.GetWidth())+"s", strings.Replace(humanize.Bytes(uint64(n.Size())), " ", "", 1))
	// Add file name
	line += " " + n.Spec.style.Render(n.Name())
	return line
}

// Get Selection status
func getSelectionStatus(n listerNode) SelectedState {
	if !n.IsDir() {
		return n.Spec.selectionStatus
	}
	sc := 0
	for _, f := range n.GetChildren() {
		if f.Spec.selectionStatus == Selected {
			sc++
		}
	}
	if sc == n.LenChildren() {
		return Selected
	}
	if sc > 0 && sc < n.LenChildren() {
		return PartialSelected
	}
	return NotSelected
}

func getSelectionStyle(n listerNode) lipgloss.Style {
	switch getSelectionStatus(n) {
	case NotSelected:
		return defaultStyle.Cursor
	case PartialSelected:
		return defaultStyle.PartialSelectedStatus
	case Selected:
		return defaultStyle.SelectedStatus
	}
	return defaultStyle.Cursor
}
