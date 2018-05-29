package server

import (
	"flag"
	"fmt"
	"github.com/DeanThompson/ginpprof"
	"github.com/aviddiviner/gin-limit"
	"os"
	// "os"
	"github.com/gin-gonic/gin"
	"github.com/itsjamie/gin-cors"
	// application modules as apps
	"github.com/nareshganesan/services/account"
	g "github.com/nareshganesan/services/globals"
	mw "github.com/nareshganesan/services/middleware"
	"github.com/nareshganesan/services/status"
	"github.com/sirupsen/logrus"
	"runtime"
	"strconv"
	"time"
)

var (
	router *gin.Engine
)

// init sets runtime settings.
func setup() error {
	// load app config details
	if err := g.LoadConfig(); err != nil {
		return err
	}
	// Configure Logrus application logger
	g.ConfigureAPILogger()
	g.ConfigureESLogger()
	l := g.Gbl.Log
	g.ConfigureDB()
	g.ConfigureElasticDB()

	// Use all CPU cores
	cores := runtime.NumCPU() - 1
	coresVal := strconv.Itoa(cores)
	runtime.GOMAXPROCS(cores)
	l.WithFields(logrus.Fields{
		"cores": coresVal,
	}).Info("No of cores")
	return nil
}

func setEnv(env string) string {
	l := g.Gbl.Log
	var currentEnv string
	if env == "" && os.Getenv("SERVICES_ENV") == "" {
		currentEnv = g.Config.Envs.Types[0]
	} else if env == "" {
		currentEnv = os.Getenv("SERVICES_ENV")
	} else if os.Getenv("SERVICES_ENV") == "" {
		currentEnv = env
	} else {
		currentEnv = g.Config.Envs.Types[0]
	}
	l.WithFields(logrus.Fields{
		"ENV": currentEnv,
	}).Info("Current Env")
	return currentEnv
}

func configCORS() cors.Config {
	cors := cors.Config{
		Origins:         "*",
		Methods:         "GET, PUT, POST, DELETE",
		RequestHeaders:  "Origin, Authorization, Content-Type",
		ExposedHeaders:  "",
		MaxAge:          50 * time.Second,
		Credentials:     true,
		ValidateHeaders: false,
	}
	return cors
}

// Serve method serves the services app api
func Serve(port, env string) {
	flag.Parse()
	if err := setup(); err != nil {
		fmt.Println(err.Error())
	} else {
		g.Config.CurrentEnv = setEnv(env)
		// Set Gin to release mode
		gin.SetMode(gin.ReleaseMode)
		router := gin.Default()
		router.Use(cors.Middleware(configCORS()))
		router.Use(limit.MaxAllowed(500))
		// router.LoadHTMLFiles("static/*")
		// router.LoadHTMLGlob("templates/*")
		router.HandleMethodNotAllowed = true
		// configure middlewares
		configureMiddlewares(router)
		// map api handlers
		v1 := router.Group("/v1")
		mapAPIHandlers(v1)
		ginpprof.Wrap(router)
		router.Run(":" + port)
	}
}

func configureMiddlewares(router *gin.Engine) {
	// router.Use(mw.DBMiddleware())
	router.Use(mw.RequestID())
	router.Use(mw.AuthMiddleware())
	router.Use(mw.LogrusMiddleware())
}

func mapAPIHandlers(router *gin.RouterGroup) {
	// configure all the url mappings
	account.RegisterAccount(router.Group("/accounts"))
	status.RegisterStatus(router.Group("/status"))
}
