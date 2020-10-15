package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/patrickmn/go-cache"
	"github.com/starkers/lander/identicon"
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
	app.Get("/favicon.ico", func(c *fiber.Ctx) error {
		name := func(c *fiber.Ctx) string {
			if c.Hostname() == "" {
				return "unknown"
			}
			// strip any trailing ports off and just return the hostname
			return strings.Split(c.Hostname(), ":")[0]
		}(c)

		fileName := fmt.Sprintf("icon-%s.png", name)

		if _, err := os.Stat(fileName); os.IsNotExist(err) {
			icon := identicon.Generate([]byte(name))

			f, err := os.Create(fileName)
			if err != nil {
				logger.Error(err)
			}
			err = icon.WriteImage(f)
			if err != nil {
				logger.Error(err)
			}
			f.Close()
			logger.Infof("rendered a new icon for: %s", name)
		}
		return c.SendFile(fileName)
	})
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
