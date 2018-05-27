package cmd

import (
	"fmt"
	g "github.com/nareshganesan/services/globals"
	"github.com/nareshganesan/services/shared"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"strings"
)

// force create index flag
var force bool

// setupesCmd represents the setupes command
var setupesCmd = &cobra.Command{
	Use:   "setupes",
	Short: "Setup ElasticSearch Index for Services (App)",
	Long:  `Creates the ElasticSearch Indices required for the Services App.`,
	Run: func(cmd *cobra.Command, args []string) {
		forceCreate := viper.GetBool("force")
		setup()
		mappingsFolder := g.Config.ProjectRoot +
			string(os.PathSeparator) +
			"mappings" +
			string(os.PathSeparator)
		// GetFiles(mappingsFolder, forceCreate)
		indexJsons := shared.GetFiles(mappingsFolder, ".json")
		for _, f := range indexJsons {
			path := mappingsFolder + f
			index := strings.Split(f, ".")[0]
			alias := index
			datepattern := "%d%02d%02d"
			newindex := index + shared.DateString(datepattern)
			fmt.Printf("creating index: %s mappings: %s\n", index, path)
			g.Gbl.CreateIndexFromJSON(newindex, path, forceCreate)
			CreateAlias(newindex, alias, forceCreate)
		}
	},
}

func init() {
	rootCmd.AddCommand(setupesCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	setupesCmd.PersistentFlags().BoolVarP(&force, "force", "f", false, "Force create new index if index exists. (default false)")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// setupesCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	viper.BindPFlag("force", setupesCmd.PersistentFlags().Lookup("force"))
}

// CreateAlias creates alias for index given index name
func CreateAlias(index, alias string, forceCreate bool) {
	if index == alias {
		fmt.Println("Index name and alias cannot be equal!")
		return
	}
	es := g.GetGlobals()
	status := es.CreateAlias(index, alias, forceCreate)
	if status {
		fmt.Printf("alias: %s created for index %s\n", alias, index)
	} else {
		fmt.Printf("Error creating alias: %s for index: %s \n", alias, index)
	}
}
