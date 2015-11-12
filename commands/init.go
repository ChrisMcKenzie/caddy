package commands

import (
	"fmt"

	"github.com/ChrisMcKenzie/caddy/pkg"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "initialize a new go package",
	Run:   initC,
}

var withBin bool

func init() {
	initCmd.Flags().BoolVarP(&withBin, "bin", "b", false, "will also create main.go for a binary package")
}

func initC(cmd *cobra.Command, args []string) {
	p := pkg.Package{
		Version: "1.0.0",
		Scripts: map[string]string{
			"setup": "echo hello world",
		},
		RawDependencies: make(map[string]string),
	}

	fmt.Printf("What is the name of your package?: ")
	_, err := fmt.Scanln(&p.Name)
	if err != nil {
		fmt.Printf(" [\033[0;31mERR\033[0m]\n")
		fmt.Printf("\t  \033[0;31m%s\033[0m\n", err)
	}

	err = pkg.WriteCaddyJSON(&p)
	if err != nil {
		fmt.Printf(" [\033[0;31mERR\033[0m]\n")
		fmt.Printf("\t  \033[0;31m%s\033[0m\n", err)
	}
}

//
// func createMainFile() {
// 	template := `
// package main
//
// func main() {
//
// }
// `
// 	err := writeTemplateToFile(ProjectPath(), "main.go", template, data)
// 	_ = err
// 	// if err != nil {
// 	// 	er(err)
// 	// }
// }
