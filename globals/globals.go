package globals

import (
	"database/sql"
	// Mysql driver
	_ "github.com/go-sql-driver/mysql"
	"github.com/olivere/elastic"
	"github.com/sirupsen/logrus"
)

// Globals application globals to be shared across api layer
type Globals struct {
	DB    *sql.DB
	Log   *logrus.Logger
	ES    *elastic.Client
	ESLog *logrus.Logger
}

// Gbl application global entity
var Gbl Globals

// GetGlobals return global application entity object
func GetGlobals() *Globals {
	return &Gbl
}
