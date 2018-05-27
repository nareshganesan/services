package globals

import (
	"fmt"
	rotatelogs "github.com/lestrrat/go-file-rotatelogs"
	"github.com/onrik/logrus/filename"
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

	var Log *logrus.Logger
	if Config.CurrentEnv != "testing" {
		CreateFolder(logFolder)
		CreateFile(infoFile)
		CreateFile(errorFile)
		logrus.SetLevel(logrus.InfoLevel)
		Log = logrus.New()
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
		hook := lfshook.NewHook(
			lfshook.WriterMap{
				logrus.InfoLevel:  iwriter,
				logrus.ErrorLevel: ewriter,
			},
			&logrus.JSONFormatter{},
		)
		Log.Hooks.Add(hook)
	} else {
		logrus.SetLevel(logrus.ErrorLevel)
		Log = logrus.New()
		Log.Out = os.Stdout
	}

	// do not use the filename hook in production!
	Log.AddHook(filename.NewHook())
	// Log.AddHook(filename.NewHook())
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
