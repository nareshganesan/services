package globals

import (
	"database/sql"
	// mysql driver
	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
	"os"
)

func newDB(driver, dataSourceName string) (*sql.DB, error) {
	l := Gbl.Log
	db, err := sql.Open(driver, dataSourceName)
	if err != nil {
		l.WithFields(logrus.Fields{
			"error": err,
		}).Error("Error creating new MySQL server connection")
		return nil, err
	}
	if err = db.Ping(); err != nil {
		l.WithFields(logrus.Fields{
			"error": err,
		}).Error("Error pinging MySQL server")
		return nil, err
	}
	return db, nil
}

// ConfigureDB configures global application object with MySQL DB connection
func ConfigureDB() {
	l := Gbl.Log
	driver := Config.Database.Driver
	dsn := Config.Database.DSN()
	serverName, err := os.Hostname()
	if err != nil {
		l.WithFields(logrus.Fields{
			"error": err,
		}).Error("Could not get server name.")
		os.Exit(1)
	}
	Config.Servername = serverName
	database, err := newDB(driver, dsn)
	if err != nil {
		l.WithFields(logrus.Fields{
			"error": err,
		}).Error("Could not create MySQL server connection.")
		os.Exit(1)
	}

	if err = database.Ping(); err != nil {
		l.WithFields(logrus.Fields{
			"error": err,
		}).Error("Could not ping MySQL server")
		os.Exit(1)
	}

	database.SetMaxOpenConns(1000)
	database.SetMaxIdleConns(1000)
	Gbl.DB = database
}
