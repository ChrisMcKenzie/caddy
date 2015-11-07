package commands

import (
	"fmt"
	"log"

	"github.com/ChrisMcKenzie/caddy/install"
	"github.com/ChrisMcKenzie/caddy/module"
	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "install dependencies for the package in the cwd",
	Run:   exec,
}

func exec(cmd *cobra.Command, args []string) {
	pkg, err := module.ReadPackageJSON()
	if err != nil {
		log.Fatal(err)
	}

	for name, version := range pkg.Dependencies {
		install.Download("vendor/"+name, name, version)
	}
	fmt.Println()
}
