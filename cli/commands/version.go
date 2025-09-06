package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of OrcaAI CLI",
	Long:  `All software has versions. This is OrcaAI's.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("OrcaAI CLI v1.0.0")
	},
}
