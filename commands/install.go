package commands

import (
	"fmt"
	"log"
	"time"

	"github.com/ChrisMcKenzie/caddy/install"
	"github.com/ChrisMcKenzie/caddy/pkg"
	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "install dependencies for the package in the cwd",
	Run:   exec,
}

func exec(cmd *cobra.Command, args []string) {
	p, err := pkg.ReadPackageJSON()
	if err != nil {
		log.Fatal(err)
		return
	}

	for _, dep := range p.Dependencies {
		fmt.Printf("\n      Installing %s@%s...", dep.Name, dep.Spec)
		s := spinner.New(spinner.CharSets[21], 100*time.Millisecond) // Build our new spinner
		s.Start()
		err := install.Download(dep)
		s.Stop()
		if err != nil {
			fmt.Printf(" [\033[0;31mERR\033[0m]\n")
			fmt.Printf("\t  \033[0;31m%s\033[0m\n", err)
			continue
		}
		fmt.Printf("\r[\033[0;32mOK\033[0m]")
	}
	fmt.Printf("\n\n")
}
