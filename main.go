package main

import (
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/patrickmn/go-cache"
	"github.com/withmandala/go-log"
)

var cacheShort = cache.New(5*time.Minute, 5*time.Minute)

func newLogger(debug bool) *log.Logger {
	// check if debug enabled
	if debug {
		logger := log.New(os.Stdout).WithDebug().WithColor()
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
	logger := newLogger(true)
	//listIngresses(logger)

	//initDatabase(logger)
	//defer database.DBConn.Close()

	fiberCfg := fiber.Config{
		DisableStartupMessage: true,
	}
	app := fiber.New(fiberCfg)
	app.Get("/", func(c *fiber.Ctx) error {
		logger.Info(c.Hostname())
		foo := getEndpoints(logger)
		return c.JSON(foo)
	})
	logger.Fatal(app.Listen(":8000"))
}
