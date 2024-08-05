/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	tea "github.com/charmbracelet/bubbletea"
	log "github.com/charmbracelet/log"
	"github.com/franciscolkdo/guntar/terminal"
	"github.com/spf13/cobra"
)

// exploreCmd represents the explore command
var exploreCmd = &cobra.Command{
	Use:   "explore",
	Args:  cobra.ExactArgs(1),
	Short: "Explore tar archive in memory",
	Long: `Explore your tar archive in memory directly in your cli:

You can browse, look into files and extract selected files/folders.
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		file, err := os.Open(args[0])
		if err != nil {
			log.Error("failed to open given file", "err", err)
			return err
		}
		terminal, err := terminal.New(file)

		if err != nil {
			log.Error("failed to create terminal", "err", err)
			return err
		}

		_, err = tea.NewProgram(terminal, tea.WithMouseCellMotion()).Run()
		if err != nil {
			log.Error("failed on quit program", "err", err)
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(exploreCmd)
}
