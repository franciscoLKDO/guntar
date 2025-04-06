package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var output string = "./extracted"

func parseExtractPath() error {
	if strings.HasPrefix(output, "~/") {
		dirname, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("failed to get user home dir: %s", err)
		}
		output = filepath.Join(dirname, output[2:])
	}
	return nil
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use: "guntar",
	CompletionOptions: cobra.CompletionOptions{
		DisableDefaultCmd: true,
	},
	Short: "Guntar your archive like a pro",
	Long: `Guntar is a cli experience for tar archives:

It can read tar archive and allow you to browse, read and extract files directly in memory.
⚠️ Limitations:
  - All parent directories must be explicitly included in the archive.
    For example, 'a/b/c/file.txt' requires 'a/', 'a/b/', and 'a/b/c/'
    to be included as separate entries for correct extraction.
  - GunTar does not handle symbolic links (symlinks). Only regular files and directories are processed.
`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
