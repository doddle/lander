package main

import (
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/patrickmn/go-cache"
	"github.com/withmandala/go-log"
)

// create a short cache.. to prevent us hammering the kube api
var cacheShort = cache.New(1*time.Minute, 2*time.Minute)

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
	// app.Get("/", func(c *fiber.Ctx) error {
	// 	return c.SendString("hello")
	// })
	app.Get("/v1/endpoints", func(c *fiber.Ctx) error {
		logger.Info("v1/endpoints")
		// get ALL endpoints
		allEndpoints := getEndpoints(logger)
		// lets filter them for only ones matching the hostname of the context
		matchedHostnames := onlyHostnamesContaining(allEndpoints, c.Hostname())
		return c.JSON(matchedHostnames)
	})
	app.Static("/", "./public")
	logger.Fatal(app.Listen(":8000"))
}
