package cmd

import (
	g "github.com/nareshganesan/services/globals"
)

func setup() {
	// load app config details
	g.LoadConfig()
	// Configure Logrus application logger
	g.ConfigureAPILogger()
	g.ConfigureESLogger()
	g.ConfigureDB()
	g.ConfigureElasticDB()
}
