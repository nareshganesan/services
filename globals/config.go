package globals

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
)

// ownerInfo config entity
type ownerInfo struct {
	Name         string
	Organization string
	Domain       string
}

// databaseInfo config entity
type databaseInfo struct {
	Username  string
	Password  string
	Hostname  string
	Port      int
	DBName    string
	Parameter string
	ConnMax   int
	Enabled   bool
	Driver    string
}

// elastic config entity
type es struct {
	Username string
	Password string
	Sniff    bool
	Health   int
	Gzip     bool
	Retries  int
	Urls     []string
	Index    index
}

// index config entity
type index struct {
	Accounts accounts
}

// accounts index config entity
type accounts struct {
	Name    string
	DocType string
}

// loggerInfo config entity
type loggerInfo struct {
	Path       string
	InfoFile   string
	ErrorFile  string
	Maxsize    int
	Maxbackups int
	Maxage     int
}

// esloggerInfo config entity
type esloggerInfo struct {
	Path       string
	InfoFile   string
	ErrorFile  string
	Maxsize    int
	Maxbackups int
	Maxage     int
}

// tokens config entity
type tokens struct {
	Auth  auth
	Email email
}

// auth token config eitity
type auth struct {
	Secret    string
	Algorithm string
	Maxage    int // no of days
}

// email token config entity
type email struct {
	Secret    string
	Algorithm string // Algorithm used for signing HSA512, HSA256
	Maxage    int    // no of days
}

// smtpInfo config entity
type smtpInfo struct {
	Host    string
	Port    string
	ToEmail []string
}

// yamlConfig config entity
type yamlConfig struct {
	Title       string
	AppName     string
	Owner       ownerInfo
	Database    databaseInfo
	ES          es
	Logger      loggerInfo
	ESLogger    esloggerInfo
	Tokens      tokens
	Servername  string
	SMTP        smtpInfo
	ProjectRoot string
}

// Config object used globally across the services app
var Config yamlConfig

// LoadConfig helper to load config entity
func LoadConfig() {
	fmt.Println("Loading configuration")
	viper.Unmarshal(&Config)
	Config.ProjectRoot = SetProjectHome()
	// fmt.Println(Config)
}

// DSN returns the data source name for MySQL DB
func (dbconfig *databaseInfo) DSN() string {
	// Example: root:@tcp(localhost:3306)/test
	return dbconfig.Username +
		":" +
		dbconfig.Password +
		"@tcp(" +
		dbconfig.Hostname +
		":" +
		fmt.Sprintf("%d", dbconfig.Port) +
		")/" +
		dbconfig.DBName + dbconfig.Parameter
}

// Folder returns Api logger's log folder path
func (l *loggerInfo) Folder() string {
	return l.Path
}

// InfoPath returns Info log file path for Api logger
func (l *loggerInfo) InfoPath() string {
	return l.Path + string(os.PathSeparator) + l.InfoFile
}

// ErrorPath returns Error log file path for Api logger
func (l *loggerInfo) ErrorPath() string {
	return l.Path + string(os.PathSeparator) + l.ErrorFile
}

// Folder returns ElasticSearch logger's log folder path
func (l *esloggerInfo) Folder() string {
	return l.Path
}

// InfoPath returns Info log file path for ElasticSearch logger
func (l *esloggerInfo) InfoPath() string {
	return l.Path + string(os.PathSeparator) + l.InfoFile
}

// ErrorPath returns Error log file path for ElasticSearch logger
func (l *esloggerInfo) ErrorPath() string {
	return l.Path + string(os.PathSeparator) + l.ErrorFile
}
