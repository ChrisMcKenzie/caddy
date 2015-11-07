package commands

import (
	"os"

	"github.com/spf13/cobra"
)

var caddyCmd = &cobra.Command{
	Use:   "caddy",
	Short: "A package management tool for go",
}

func Execute() {
	AddCommands()
	if err := caddyCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}

func AddCommands() {
	caddyCmd.AddCommand(installCmd)
}
