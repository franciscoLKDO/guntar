package cmd

import (
	"fmt"
	"os"

	"github.com/franciscolkdo/guntar/tar"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list <archive file>",
	Short: "List all files in current archive",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		file, err := os.Open(args[0])
		if err != nil {
			return fmt.Errorf("failed to open given file: %s", err)
		}
		ls, err := tar.List(file)
		if err != nil {
			return fmt.Errorf("failed to list archive: %s", err)
		}
		for _, f := range ls {
			fmt.Println(f)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
