package commands

import "github.com/spf13/cobra"

var publishCmd = &cobra.Command{
	Use:   "init",
	Short: "initialize a new go package",
	Run:   publish,
}

func publish(cmd *cobra.Command, args []string) {

}
