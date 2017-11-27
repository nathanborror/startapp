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
	RootCmd.PersistentFlags().String("graphql-schema", "", "GraphQL schema")
	RootCmd.PersistentFlags().Bool("ios-backend-scaffolding", true, "Output iOS backend scaffolding")
	RootCmd.PersistentFlags().Bool("ios-test-scaffolding", true, "Output iOS tests scaffolding")
	RootCmd.PersistentFlags().String("ios-product-name", "", "iOS Product name")
	RootCmd.PersistentFlags().String("ios-team-id", "", "iOS Team ID")

	viper.BindPFlag("domain", RootCmd.PersistentFlags().Lookup("domain"))
	viper.BindPFlag("graphql-schema", RootCmd.PersistentFlags().Lookup("graphql-schema"))
	viper.BindPFlag("ios-backend-scaffolding", RootCmd.PersistentFlags().Lookup("ios-backend-scaffolding"))
	viper.BindPFlag("ios-test-scaffolding", RootCmd.PersistentFlags().Lookup("ios-test-scaffolding"))
	viper.BindPFlag("ios-product-name", RootCmd.PersistentFlags().Lookup("ios-product-name"))
	viper.BindPFlag("ios-team-id", RootCmd.PersistentFlags().Lookup("ios-team-id"))
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
	proj.AddIOSClient(viper.GetString("ios-product-name"), viper.GetString("ios-team-id"), viper.GetBool("ios-backend-scaffolding"), viper.GetBool("ios-test-scaffolding"))
	proj.ReadGraphQLSchema(viper.GetString("graphql-schema"))
	proj.Copy(viper.GetString("graphql-schema"))
	proj.Write()
	checkErr(proj.Err())
}

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
