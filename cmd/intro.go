package cmd

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/franciscolkdo/guntar/terminal"
	"github.com/spf13/cobra"
)

// introCmd represents the intro command
var introCmd = &cobra.Command{
	Use:   "intro",
	Short: "Intro animation",
	RunE: func(cmd *cobra.Command, args []string) error {
		_, err := tea.NewProgram(terminal.NewIntroModel()).Run()
		if err != nil {
			return fmt.Errorf("failed on quit program: %s", err)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(introCmd)
}
