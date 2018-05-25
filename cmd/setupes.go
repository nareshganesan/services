package cmd

import (
	"fmt"
	g "github.com/nareshganesan/services/globals"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"strings"
	"time"
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
		getFiles(mappingsFolder, forceCreate)
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

func getFiles(folderPath string, forceCreate bool) {
	err := filepath.Walk(folderPath, func(path string, f os.FileInfo, err error) error {
		if err != nil {
			fmt.Println("Walkpath error")
			fmt.Println(err)
		} else {

			if !f.IsDir() && strings.Contains(f.Name(), ".json") {
				index := strings.Split(f.Name(), ".")[0]
				fmt.Printf("creating index: %s mappings: %s\n", index, path)
				now := time.Now().UTC()
				suffix := fmt.Sprintf("%d%02d%02d", now.Year(), now.Month(), now.Day())
				suffix = "-" + suffix
				alias := index
				newindex := index + suffix
				createIndexFromJSON(newindex, path, forceCreate)
				createAlias(newindex, alias, forceCreate)
			}
		}
		return nil
	})
	if err != nil {
		fmt.Println("Get Files error")
		fmt.Printf(err.Error())
	}
}

func createIndexFromJSON(index, mappingsFile string, forceCreate bool) {
	es := g.GetGlobals()
	var MappingJSON map[string]interface{}
	g.LoadJSON(mappingsFile, &MappingJSON)
	status := es.CreateIndex(index, MappingJSON, forceCreate)
	if status {
		fmt.Printf("Index: %s Created\n", index)
	} else {
		fmt.Printf("Error creating %s index\n", index)
	}
}

func createAlias(index, alias string, forceCreate bool) {
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
