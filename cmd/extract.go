package cmd

import (
	"fmt"
	"os"

	"github.com/franciscolkdo/guntar/tar"
	"github.com/spf13/cobra"
)

// extractCmd represents the extract command
var extractCmd = &cobra.Command{
	Use:   "extract <archive>",
	Short: "Extract archive",
	RunE: func(cmd *cobra.Command, args []string) error {
		file, err := os.Open(args[0])
		if err != nil {
			return fmt.Errorf("failed to open given file: %s", err)
		}
		node, err := tar.Scan(file, func(n *tar.SimpleNode) error { return nil })
		if err != nil {
			return fmt.Errorf("failed to list archive: %s", err)
		}

		if err := parseExtractPath(); err != nil {
			return err
		}

		return tar.Extract(node, output, func(n *tar.SimpleNode) bool { return false })
	},
}

func init() {
	rootCmd.AddCommand(extractCmd)
	extractCmd.Flags().StringVarP(&output, "output", "o", "", "Output directory to extract archive")
}
