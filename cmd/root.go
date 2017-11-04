package cmd

import (
	"fmt"
	"os"

	"github.com/nathanborror/startapp/gen"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var configFlag string

var RootCmd = &cobra.Command{
	Use:   "startapp",
	Short: "StartApp is a quick way to get started on a new App",
	Long: `StartApp creates a new app folder with an iOS client and
			a partial Go project.`,
	Args: cobra.MinimumNArgs(2),
	Run:  runGenerateApp,
}

func Execute() {
	checkErr(RootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().StringVar(&configFlag, "config", "", "config file (default is ~/.startapp)")
	RootCmd.PersistentFlags().String("domain", "", "The app's domain")
	RootCmd.PersistentFlags().String("schema", "", "GraphQL schema")
	RootCmd.PersistentFlags().Bool("has-client-backend", true, "Has client backend")
	RootCmd.PersistentFlags().Bool("has-client-tests", true, "Has client tests")
	RootCmd.PersistentFlags().String("has-client-bundleid", "", "Has client Bundle ID")
	RootCmd.PersistentFlags().String("has-client-teamid", "", "Has client Team ID")

	viper.BindPFlag("domain", RootCmd.PersistentFlags().Lookup("domain"))
	viper.BindPFlag("schema", RootCmd.PersistentFlags().Lookup("schema"))
	viper.BindPFlag("has-client-backend", RootCmd.PersistentFlags().Lookup("has-client-backend"))
	viper.BindPFlag("has-client-tests", RootCmd.PersistentFlags().Lookup("has-client-tests"))
	viper.BindPFlag("has-client-bundleid", RootCmd.PersistentFlags().Lookup("has-client-bundleid"))
	viper.BindPFlag("has-client-teamid", RootCmd.PersistentFlags().Lookup("has-client-teamid"))
}

func initConfig() {
	if configFlag != "" {
		viper.SetConfigFile(configFlag)
	} else {
		viper.SetConfigName(".startapp")
		viper.AddConfigPath("$HOME")
	}
	viper.AutomaticEnv()
	viper.SetEnvPrefix("startapp")

	err := viper.ReadInConfig()
	if err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func runGenerateApp(cmd *cobra.Command, args []string) {
	name := args[0]
	dest := args[1]
	proj := gen.NewProject(name, dest, viper.GetString("domain"))
	proj.ReadGraphQLSchema(viper.GetString("schema"))
	proj.AddClient(gen.IOSClientKind, viper.GetString("has-client-bundleid"), viper.GetString("has-client-teamid"), viper.GetBool("has-client-backend"), viper.GetBool("has-client-tests"))
	proj.Write()
	proj.Copy(viper.GetString("schema"))
	checkErr(proj.Err())
}

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
