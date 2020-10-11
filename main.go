package main

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber"
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

func initDatabase() {
	var err error
	database.DBConn, err = gorm.Open("sqlite3", "cache.db")
	if err != nil {
		panic("failed to connect database")
	}
	fmt.Println("Connection Opened to Database")
}

func main() {
	initDatabase()
	app := fiber.New()
	setupRoutes(app)
	app.Listen(8000)
	defer database.DBConn.Close()
}

func helloWorld(c *fiber.Ctx) {
	c.Send("Hello, World!")
}

func setupRoutes(app *fiber.App) {
	app.Get("/", helloWorld)
	app.Get("/v1/ingress", getIngress)
}

func getIngress(c *fiber.Ctx) {
	c.Send("foo")
}
