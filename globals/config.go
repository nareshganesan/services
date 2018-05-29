package globals

import (
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"os"
	"reflect"
)

// ownerInfo config entity
type ownerInfo struct {
	Name         string
	Organization string
	Domain       string
	Email        string
}

type env struct {
	Types []string
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
	CurrentEnv  string
	Envs        env
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
func LoadConfig() error {
	fmt.Println("Loading configuration")
	viper.Unmarshal(&Config)
	Config.ProjectRoot = SetProjectHome()
	// fmt.Println(Config)
	return Config.ValidateESIndex()
}

func (cfg *yamlConfig) ValidateESIndex() error {
	esIndex := cfg.ES.Index
	fields := reflect.ValueOf(&esIndex).Elem()
	// fieldType := fields.Type()
	for i := 0; i < fields.NumField(); i++ {
		// fieldName := fieldType.Field(i).Name
		fieldValue := fields.Field(i).Interface()
		switch index := fieldValue.(type) {
		case accounts:
			if index.Name == "" {
				return errors.New("Config validation error!. Accounts index name cannot be empty! ")
			}
			if index.DocType == "" {
				return errors.New("Config validation error!. Accounts index doctype cannot be empty! ")
			}
			if index.Name != "" && index.DocType != "" {
				fmt.Println(fmt.Sprintf("Accounts: index: %s, doctype: %s", index.Name, index.DocType))
			}
		default:
			fmt.Println("index type is unknown")
			return errors.New("Configuration validation! Unkown index! ")
		}
	}
	return nil
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

// InitConfig reads in config file and ENV variables if set.
func InitConfig(cfgFile string) {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		projectHome, err := os.Getwd()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(projectHome)

		// Search config in home directory with name ".config" (without extension).
		viper.AddConfigPath(projectHome)
		viper.SetConfigName(".config")
		viper.SetConfigType("yaml")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	} else {
		fmt.Println(err)
	}
}
