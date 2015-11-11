package commands

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/ChrisMcKenzie/caddy/pkg"
	"github.com/spf13/cobra"
)

var buildCmd = &cobra.Command{
	Use:     "build",
	Aliases: []string{"b"},
	Short:   "Build will install any missing dependencies and build binary",
	Long: `
If no build script is specified in the caddy.json a default "go build" will
execute. other wise it will run the script specified
	`,
	Run: build,
}

var output string

func init() {
	buildCmd.Flags().StringVarP(&output, "out", "o", "", "write binary to this file path.")
}

func build(cmd *cobra.Command, args []string) {
	p, err := pkg.ReadCaddyJSON()
	if err != nil {
		log.Fatal(err)
		return
	}

	if build, ok := p.Scripts["build"]; ok {
		fmt.Printf("Running script \"build\"\n%s", build)
		fields := strings.Fields(build)

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
	} else {
		if strings.Contains(runtime.Version(), "go1.5") {
			// add GO15VENDOREXPERIMENT=1 env var so it uses local "vendor""
			var goArgs []string

			if output != "" {
				goArgs = []string{"-o", output}
			}

			goArgs = append(goArgs, args...)
			goArgs = append([]string{"build"}, goArgs...)

			cmd := exec.Command("go", goArgs...)
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
		}
	}
}
