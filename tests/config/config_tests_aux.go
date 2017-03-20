package config

import (
	"gost/config"
	"gost/httphandle"
	"gost/servers"
	"log"
	"os"
)

const (
	envApplicationName    = "GOST_TESTAPP_NAME"
	envAPIInstance        = "GOST_TESTAPP_INSTANCE"
	envHTTPServerAddress  = "GOST_TESTAPP_HTTP"
	envDatabaseName       = "GOST_TESTAPP_DB_NAME"
	envDatabaseConnection = "GOST_TESTAPP_DB_CONN"
)

// InitTestsApp initializes the application used for testing
func InitTestsApp() {
	config.ApplicationName = os.Getenv(envApplicationName)
	config.APIInstance = os.Getenv(envAPIInstance)
	config.HTTPServerAddress = os.Getenv(envHTTPServerAddress)
}

// InitTestsDatabase initializes the connection parameters to the database used for testing
func InitTestsDatabase() {
	dbName := os.Getenv(envDatabaseName)
	dbConn := os.Getenv(envDatabaseConnection)

	if len(dbName) == 0 || len(dbConn) == 0 {
		log.Fatal("Environment variables for the test database are not set!")
	}

	config.DbName = dbName
	config.DbConnectionString = dbConn
}

// InitTestsRoutes initializez the routes used for testing the endpoints
func InitTestsRoutes() {
	httphandle.CreateFrameworkRoutes()
	httphandle.CreateDevelopmentRoutes()
	httphandle.CreateAPIRoutes()
	httphandle.InitRoutes(servers.Multiplexer)
}
