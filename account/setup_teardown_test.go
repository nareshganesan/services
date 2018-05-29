package account

import (
	"fmt"
	"github.com/gin-gonic/gin"
	g "github.com/nareshganesan/services/globals"
	mw "github.com/nareshganesan/services/middleware"
	"github.com/nareshganesan/services/shared"
	"github.com/olivere/elastic"
	"os"
	// "github.com/spf13/viper"
	"strings"
	"testing"
)

var cfgFile string

func TestMain(m *testing.M) {
	// call flag.Parse() here if TestMain uses flags
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

func setup() {
	fmt.Println("Setup for account package of services app")
	projectRoot := string(os.Getenv("GOPATH")) + "/src/github.com/nareshganesan/services"
	cfgFile = projectRoot + string(os.PathSeparator) + ".config.yaml"
	g.InitConfig(cfgFile)
	g.LoadConfig()

	g.Config.ProjectRoot = projectRoot
	mappingsFolder := g.Config.ProjectRoot +
		string(os.PathSeparator) +
		"mappings" +
		string(os.PathSeparator)
	g.Config.CurrentEnv = "testing"
	// Configure Logrus application logger
	g.ConfigureAPILogger()
	g.ConfigureESLogger()
	g.Gbl.ES = MockESService(g.Config.ES.Urls[0])
	CreateTestIndex(mappingsFolder)
}

func teardown() {
	fmt.Println("teardown for account package of services app")
	projectRoot := string(os.Getenv("GOPATH")) + "/src/github.com/nareshganesan/services"
	g.Config.ProjectRoot = projectRoot
	g.Config.ProjectRoot = projectRoot
	mappingsFolder := g.Config.ProjectRoot +
		string(os.PathSeparator) +
		"mappings" +
		string(os.PathSeparator)
	DeleteTestIndex(mappingsFolder)

}

// MockESService creates a elastic search client with localhost for testing
func MockESService(url string) *elastic.Client {
	client, err := elastic.NewSimpleClient(elastic.SetURL(url))
	if err != nil {
		return nil
	}
	return client
}

// GetRouter helper to create a router during testing
func GetRouter(withTemplates bool) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/login", Login)
	router.POST("/signup", Signup)
	router.POST("/update", mw.AuthDecorator(), UpdateAccount)
	router.POST("/delete", mw.AuthDecorator(), DeleteAccount)
	router.POST("/list", mw.AuthDecorator(), ListAccount)
	return router
}

// CreateTestIndex creates the test index & alias for account testing
func CreateTestIndex(mappingsFolder string) {
	indexJsons := shared.GetFiles(mappingsFolder, ".json")
	for _, f := range indexJsons {
		path := mappingsFolder + f
		index := strings.Split(f, ".")[0]
		alias := index
		datepattern := "%d%02d%02d"
		newindex := "testing-" + index + shared.DateString(datepattern)
		fmt.Printf("creating index: %s mappings: %s\n", index, path)
		fmt.Printf("creating newindex: %s alias: %s\n", newindex, alias)
		forceCreate := true
		// Delete all existing indexes for alias
		existingIndexes, _ := g.Gbl.GetIndexesByAlias(index)
		for _, idx := range existingIndexes {
			_ = g.Gbl.DeleteIndex(idx)
		}
		g.Gbl.CreateIndexFromJSON(newindex, path, forceCreate)
		CreateAlias(newindex, alias, forceCreate)
	}
	CreateTestAccount()
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

// DeleteTestIndex deletes the test index & alias after testing
func DeleteTestIndex(mappingsFolder string) {
	indexJsons := shared.GetFiles(mappingsFolder, ".json")
	for _, f := range indexJsons {
		index := strings.Split(f, ".")[0]
		alias := index
		datepattern := "%d%02d%02d"
		newindex := "testing-" + index + shared.DateString(datepattern)
		fmt.Printf("deleting test index: %s alias: %s\n", index, alias)
		status := g.Gbl.DeleteAlias(newindex, alias)
		fmt.Printf("deleting test index: %s\n", newindex)
		_ = g.Gbl.DeleteIndex(newindex)
		fmt.Printf("%s alias for index %s deleted - %v", alias, newindex, status)
	}
}

func CreateTestAccount() {
	var account Entity
	account.Email = "test1@email.com"
	account.Password = "Testpassword#123"
	id, status := account.Create()
	if status {
		fmt.Println("account created")
		fmt.Println(id)
	} else {
		fmt.Println("Could not create account")
	}

}
