package cmd

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/franciscolkdo/guntar/terminal"
	"github.com/spf13/cobra"
)

// exploreCmd represents the explore command
var exploreCmd = &cobra.Command{
	Use:   "explore <archive file>",
	Args:  cobra.ExactArgs(1),
	Short: "Explore tar archive in memory",
	Long: `Explore your tar archive in memory directly in your cli:

You can browse, look into files and extract selected files/folders.
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		file, err := os.Open(args[0])
		if err != nil {
			return fmt.Errorf("failed to open given file: %s", err)
		}
		terminal, err := terminal.New(file)

		if err != nil {
			return fmt.Errorf("failed to create terminal: %s", err)
		}

		_, err = tea.NewProgram(terminal, tea.WithMouseCellMotion()).Run()
		if err != nil {
			return fmt.Errorf("failed on quit program: %s", err)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(exploreCmd)
}
