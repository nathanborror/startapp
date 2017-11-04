package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var version = "0.1"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Long:  "All software has versions.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(versionString())
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}

func versionString() string {
	return fmt.Sprintf("{{.Name|titlecase}} v%s\n", version)
}
