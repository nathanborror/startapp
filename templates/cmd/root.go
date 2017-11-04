package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/nathanborror/{{.Name}}/state"
	"github.com/nathanborror/{{.Name}}/state/postgres"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	configFlag   string
	verboseFlag  bool
	stateBackend state.Stater
)

// RootCmd affords the main command.
var RootCmd = &cobra.Command{
	Use:   "{{.Name|titlecase}}",
	Short: "{{.Name|titlecase}} is an all-purpose tool for working with API.",
	Long: `An all-purpose utility for deploying {{.Name|titlecase}} and interacting
		   with the local or remote server.`,
}

// Execute executes the root command.
func Execute() {
	checkErr(RootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().BoolVarP(&verboseFlag, "verbose", "v", false, "verbose output")
	RootCmd.PersistentFlags().StringVar(&configFlag, "config", "", "config file (default is ~/.{{.Name}})")
	RootCmd.PersistentFlags().String("state", "postgres", "state backend to use")
	RootCmd.PersistentFlags().String("schema", "schema.graphql", "GraphQL schema file to use")
	RootCmd.PersistentFlags().String("token", "", "Authorization token")
	RootCmd.PersistentFlags().String("database", "{{.Name}}", "Database to use")
	RootCmd.PersistentFlags().String("database-user", "postgres", "Database user to use")

	viper.BindPFlag("state", RootCmd.PersistentFlags().Lookup("state"))
	viper.BindPFlag("schema", RootCmd.PersistentFlags().Lookup("schema"))
	viper.BindPFlag("token", RootCmd.PersistentFlags().Lookup("token"))
	viper.BindPFlag("database", RootCmd.PersistentFlags().Lookup("database"))
	viper.BindPFlag("database-user", RootCmd.PersistentFlags().Lookup("database-user"))
}

func initConfig() {
	if configFlag != "" {
		viper.SetConfigFile(configFlag)
	}

	viper.SetConfigName(".{{.Name}}")
	viper.AddConfigPath("$HOME")
	viper.AutomaticEnv()
	viper.SetEnvPrefix("{{.Name}}")

	err := viper.ReadInConfig()
	if err == nil && verboseFlag {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

	stateCfg := make(map[string]string)
	stateCfg["Database"] = viper.GetString("database")
	stateCfg["User"] = viper.GetString("database-user")
	state.Register("postgres", postgres.NewState)
	stateBackend = state.NewState(viper.GetString("state"), stateCfg)
}

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func prettyPrint(in interface{}) string {
	b, err := json.MarshalIndent(in, "", "  ")
	checkErr(err)
	return string(b)
}
