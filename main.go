package main

import (
	"os"

	"github.com/gofiber/fiber"
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

func main() {
	app := fiber.New()
	setupRoutes(app)
	app.Listen(8000)
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
