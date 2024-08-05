/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	tea "github.com/charmbracelet/bubbletea"
	log "github.com/charmbracelet/log"
	"github.com/franciscolkdo/guntar/terminal"
	"github.com/spf13/cobra"
)

// introCmd represents the intro command
var introCmd = &cobra.Command{
	Use:   "intro",
	Short: "Intro animation",
	RunE: func(cmd *cobra.Command, args []string) error {
		_, err := tea.NewProgram(terminal.NewIntroModel(), tea.WithMouseCellMotion()).Run()
		if err != nil {
			log.Error("failed on quit program", "err", err)
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(introCmd)
}
