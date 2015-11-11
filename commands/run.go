package commands

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/ChrisMcKenzie/caddy/pkg"
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:     "run [script]",
	Aliases: []string{"r"},
	Short:   "Run will run the script specified in you caddy.json",
	Run:     run,
}

func run(cmd *cobra.Command, args []string) {
	p, err := pkg.ReadCaddyJSON()
	if err != nil {
		log.Fatal(err)
		return
	}

	if script, ok := p.Scripts[args[0]]; ok {
		fmt.Printf("Running script \"%s\"\n%s\n", args[0], script)
		fields := strings.Fields(script)

		cmd := exec.Command(fields[0], fields[1:]...)
		cmd.Env = append(os.Environ(), "GO15VENDOREXPERIMENT=1")
		var buf bytes.Buffer
		cmd.Stderr = &buf
		cmd.Stdout = &buf
		err := cmd.Run()
		out := buf.String()
		fmt.Println(out)
		if err != nil {
			fmt.Println(err)
		}
		return
	}
}
