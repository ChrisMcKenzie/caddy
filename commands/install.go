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
	Use:   "install [registry dependency]",
	Short: "install dependencies for the package in the cwd",
	Long: `
Install gets any dependencies specified in the caddy.json
if a registry dependency is given it will be installed to
the "vendor" dir in the cwd.

if --save,-s is given the dependency given will be save to the 
caddy.json.

if --global,-g is given the dependency given will not go to the 
local "vendor" dir and instead will be installed to your $GOPATH
	`,
	Run: runInstall,
}

var save bool
var isGlobal bool

func init() {
	installCmd.Flags().BoolVarP(&save, "save", "s", false, "save dependency to caddy.json")
	installCmd.Flags().BoolVarP(&isGlobal, "global", "g", false, "save dependency to global $GOPATH")
}

func runInstall(cmd *cobra.Command, args []string) {
	p, err := pkg.ReadCaddyJSON()
	if err != nil {
		log.Fatal(err)
		return
	}

	if len(args) >= 1 {
		name := args[0]
		fmt.Printf("\n      Installing %s...", name)
		s := spinner.New(spinner.CharSets[21], 100*time.Millisecond) // Build our new spinner
		s.Start()
		defer s.Stop()
		dep, err := pkg.Parse(name, "*")
		if err != nil {
			fmt.Printf(" [\033[0;31mERR\033[0m]\n")
			fmt.Printf("\t  \033[0;31m%s\033[0m\n", err)
			return
		}
		err = install.Download(isGlobal, &dep)
		if err != nil {
			fmt.Printf(" [\033[0;31mERR\033[0m]\n")
			fmt.Printf("\t  \033[0;31m%s\033[0m\n", err)
			return
		}
		fmt.Printf("\r[\033[0;32mOK\033[0m]")
		fmt.Printf("\n\n")
		if save {
			p.RawDependencies[dep.Name] = dep.Spec
			pkg.WriteCaddyJSON(p)
		}
	} else {
		for _, dep := range p.Dependencies {
			fmt.Printf("\n      Installing %s@%s...", dep.Name, dep.Spec)
			s := spinner.New(spinner.CharSets[21], 100*time.Millisecond) // Build our new spinner
			s.Start()
			defer s.Stop()
			err := install.Download(isGlobal, &dep)
			if err != nil {
				fmt.Printf(" [\033[0;31mERR\033[0m]\n")
				fmt.Printf("\t  \033[0;31m%s\033[0m\n", err)
				continue
			}
			fmt.Printf("\r[\033[0;32mOK\033[0m]")
		}
		fmt.Printf("\n\n")
	}

}
