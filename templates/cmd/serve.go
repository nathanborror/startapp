package cmd

import (
	"fmt"
	"net/http"

	"github.com/nathanborror/{{.Name}}/api"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serve the {{.Name|titlecase}} API.",
	Long: `{{.Name|titlecase}} has a web-based API for clients to access, this 
	is the command to start the web-server.`,
	Run: runServe,
}

func init() {
	RootCmd.AddCommand(serveCmd)
	serveCmd.Flags().StringP("port", "p", ":8080", "port to serve on")
	viper.BindPFlag("port", serveCmd.Flags().Lookup("port"))
}

func runServe(cmd *cobra.Command, args []string) {
	
	api.Configure(viper.GetString("schema"), api.Backends{
		State: stateBackend,
	})
	http.HandleFunc("/", indexHandler)

	fmt.Printf("Starting to serve on %s\n", viper.GetString("port"))
	checkErr(http.ListenAndServe(viper.GetString("port"), nil))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, versionString())
}