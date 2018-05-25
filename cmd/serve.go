package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/nareshganesan/services/server"
)

var port string

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serve API on port 3333 (default port)",
	Long: `Serve API on port 3333
Can be configured to run on any port using --port flag`,
	Run: func(cmd *cobra.Command, args []string) {
		p := viper.GetString("port")
		author := viper.GetString("owner.name")
		fmt.Println("Starting API server on port", p)
		fmt.Println("Author:", author)

		server.Serve(port)
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	serveCmd.PersistentFlags().StringVarP(&port, "port", "p", "3333", "Port number to serve the API (Default: 3333)")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	viper.BindPFlag("port", serveCmd.PersistentFlags().Lookup("port"))
}
