package cmd

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/nathanborror/startapp/gen"
	"github.com/spf13/cobra"
)

var schemaCmd = &cobra.Command{
	Use:   "schema",
	Short: "Schema returns the parsed GraphQL schema",
	Long:  `Schema returns the parsed GraphQL schema`,
	Args:  cobra.MinimumNArgs(1),
	Run:   runSchemaCmd,
}

func init() {
	RootCmd.AddCommand(schemaCmd)
}

func runSchemaCmd(cmd *cobra.Command, args []string) {
	filename := args[0]

	project := gen.NewProject("temp", "", "")
	project.ReadGraphQLSchema(filename)
	checkErr(project.Err())

	fmt.Printf("%+v\n", project.Definition)
}

func readFileContents(filename string) string {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Printf("Failed to open schema file '%s': %v\n", filename, err)
		os.Exit(1)
	}
	return string(file)
}
