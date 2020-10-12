package main

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/starkers/lander/database"
	"github.com/withmandala/go-log"
)

// returns true if PLUGIN_DEBUG!=""
func newLogger(debug bool) *log.Logger {
	// check if debug enabled
	if debug {
		logger := log.New(os.Stdout).WithDebug().WithColor()
		logger.Debug("debug enabled")
		return logger
	} else {
		return log.New(os.Stdout).WithColor()
	}
}

func envVarExists(key string) bool {
	_, exists := os.LookupEnv(key)
	if exists {
		return true
	}
	return false
}

func initDatabase(logger *log.Logger) {
	dbObj := "cache.db"
	var err error
	database.DBConn, err = gorm.Open("sqlite3", dbObj)
	if err != nil {
		logger.Error(err)
		logger.Fatalf("failed to connect database: %s", dbObj)
	}
	logger.Infof("opened db: %s", dbObj)
}

func main() {
	logger := newLogger(true)
	//listIngresses(logger)

	initDatabase(logger)
	defer database.DBConn.Close()

	fiberCfg := fiber.Config{
		DisableStartupMessage: true,
	}
	app := fiber.New(fiberCfg)
	app.Get("/", func(c *fiber.Ctx) error {
		logger.Info(c.Hostname())
		listIngresses(logger)
		return c.SendString("Hello")
	})
	logger.Fatal(app.Listen(":8000"))
}

