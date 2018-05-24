package globals

import (
	"fmt"
	rotatelogs "github.com/lestrrat/go-file-rotatelogs"
	// "github.com/nareshganesan/services/shared"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

// ConfigureAPILogger configures logger object to application global enitty
func ConfigureAPILogger() {
	fmt.Println("Configuring Api logger")

	folder := Config.ProjectRoot + string(os.PathSeparator) + Config.Logger.Folder()
	path := Config.ProjectRoot
	infoFilePath := path + string(os.PathSeparator) + Config.Logger.InfoPath()
	errorFilePath := path + string(os.PathSeparator) + Config.Logger.ErrorPath()
	Gbl.Log = getLogger(folder, infoFilePath, errorFilePath)
}

// ConfigureESLogger configures ElasticSearch logger object to application global entity
func ConfigureESLogger() {
	fmt.Println("Configuring ES logger")

	folder := Config.ProjectRoot + string(os.PathSeparator) + Config.ESLogger.Folder()
	path := Config.ProjectRoot
	infoFilePath := path + string(os.PathSeparator) + Config.ESLogger.InfoPath()
	errorFilePath := path + string(os.PathSeparator) + Config.ESLogger.ErrorPath()
	Gbl.ESLog = getLogger(folder, infoFilePath, errorFilePath)
}

func getLogger(logFolder, infoFile, errorFile string) *logrus.Logger {

	CreateFolder(logFolder)
	CreateFile(infoFile)
	CreateFile(errorFile)

	iwriter, err := rotateLogsWriter(infoFile)
	if err != nil {
		fmt.Println("Error configuring info log file.")
		fmt.Println(err)
	}
	ewriter, err := rotateLogsWriter(errorFile)
	if err != nil {
		fmt.Println("Error configuring error log file.")
		fmt.Println(err)
	}

	// Only log the warning severity or above.
	logrus.SetLevel(logrus.InfoLevel)
	Log := logrus.New()
	hook := lfshook.NewHook(
		lfshook.WriterMap{
			logrus.InfoLevel:  iwriter,
			logrus.ErrorLevel: ewriter,
		},
		&logrus.JSONFormatter{},
	)
	Log.Hooks.Add(hook)
	return Log
}

func rotateLogsWriter(path string) (*rotatelogs.RotateLogs, error) {
	return rotatelogs.New(
		path+".%Y%m%d%H%M",
		rotatelogs.WithLinkName(path),
		rotatelogs.WithMaxAge(time.Duration(Config.Logger.Maxage)*time.Hour),
		rotatelogs.WithRotationTime(time.Hour),
	)
}
