package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/dreikanter/npub"
	"github.com/dreikanter/npub/internal/config"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init [path]",
	Short: "Create a sample npub.yml in the given directory (default: current directory)",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		dir := "."
		if len(args) > 0 {
			dir = args[0]
		}
		dir = expandHome(os.ExpandEnv(dir))

		if err := os.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf("create directory: %w", err)
		}

		target := filepath.Join(dir, config.DefaultConfigFile)
		if _, err := os.Stat(target); err == nil {
			return fmt.Errorf("%s already exists", target)
		} else if !os.IsNotExist(err) {
			return fmt.Errorf("check %s: %w", target, err)
		}

		if err := os.WriteFile(target, npub.SampleConfig, 0o644); err != nil {
			return fmt.Errorf("write config: %w", err)
		}

		fmt.Fprintf(cmd.OutOrStdout(), "Created %s\n", target)
		return nil
	},
}
