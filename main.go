package main

import (
	"flag"
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
	}
	return log.New(os.Stdout).WithColor()
}

func envVarExists(key string) bool {
	_, exists := os.LookupEnv(key)
	if exists {
		return true
	}
	return false
}

func main() {

	flagHost := flag.String("host", "clustername.example.com", "filter ingresses matching this hostname")
	flagExcludeEndpoints := flag.String("excludeEndpoints", "https://example.com/foo,https://example.com/", "exclude (comma sperated) specific endpoints")
	flagDebug := flag.Bool("debug", false, "enable or disable debug logging")

	flag.Parse()
	fmt.Println(*flagDebug)

	logger := newLogger((*flagDebug))

	fiberCfg := fiber.Config{
		DisableStartupMessage: true,
	}
	app := fiber.New(fiberCfg)
	app.Get("/img/icons/favicon*", func(c *fiber.Ctx) error {
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
	app.Get("/favicon-k8s.*", func(c *fiber.Ctx) error {
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
		// get ALL endpoints
		allEndpoints := getEndpoints(logger)
		// lets filter them for only ones matching the hostname of the context
		matchedHostnames := onlyHostnamesContaining(allEndpoints, *flagHost, *flagExcludeEndpoints)
		// matchedHostnames := onlyHostnamesContaining(allEndpoints, c.Hostname())
		logger.Infof("/v1/endpoints filtered %v known endpoints and returned %v results", len(allEndpoints), len(matchedHostnames))
		return c.JSON(matchedHostnames)
	})

	app.Static("/", "./frontend/dist")

	onStartup(logger)

	logger.Info("starting webserver on :8000")
	logger.Fatal(app.Listen(":8000"))
}

func onStartup(logger *log.Logger) {
	logger.Info("getting some initial data bootstrapped")
	_ = getEndpoints(logger)
}
