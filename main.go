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

var (
	// create a short cache.. to prevent us hammering the kube api
	cacheShort = cache.New(1*time.Minute, 2*time.Minute)

	// lets guess a colorscheme based on strings from a hostname

)

// Settings to be returned to the browser/client
type Settings struct {
	ColorScheme string `json:"colorscheme"`
	Cluster     string `json:"cluster"`
}

func onStartup(logger *log.Logger) {
	logger.Info("getting some initial data bootstrapped")
	_ = getEndpoints(logger)
}

func main() {

	flagHost := flag.String("host", "clustername.example.com", "filter ingresses matching this hostname")
	flagExcludeEndpoints := flag.String("excludeEndpoints", "https://example.com/foo,https://example.com/", "exclude (comma separated) specific endpoints")
	flagConfig := flag.String("config", "default", "Specify a config file (customised colour scheme)")
	flagDebug := flag.Bool("debug", false, "enable or disable debug logging")

	flag.Parse()
	//fmt.Println(*flagDebug)

	logger := newLogger(*flagDebug)

	fiberCfg := fiber.Config{
		DisableStartupMessage: true,
	}

	//fmt.Println(*flagConfig)
	// sets up an initial config (for colourschemes)
	cfg := initialConfig(logger, *flagConfig)

	logger.Info(cfg)

	app := fiber.New(fiberCfg)

	app.Get("*favicon*", func(c *fiber.Ctx) error {
		var name string
		configuredHostname := *flagHost
		if configuredHostname == "clustername.example.com" {
			if c.Hostname() == "" {
				name = "unknown"
			}
			name = strings.Split(c.Hostname(), ":")[0]
		}
		// the -host param was used.. infer it
		name = configuredHostname
		fileName := fmt.Sprintf("icon-%s.png", name)
		hex := guessHex(logger, cfg.Colorschemes, name)
		logger.Debugf("guessed the hex: %s for %s", hex, name)
		if _, err := os.Stat(fileName); os.IsNotExist(err) {
			icon := identicon.Generate([]byte(name), hex)

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
		// logger.Debugf("served %s", fileName)
		return c.SendFile(fileName)
	})

	app.Get("/v1/endpoints", func(c *fiber.Ctx) error {
		// get ALL endpoints
		allEndpoints := getEndpoints(logger)
		// lets filter them for only ones matching the hostname of the context
		matchedHostnames := onlyHostnamesContaining(allEndpoints, *flagHost, *flagExcludeEndpoints)
		// matchedHostnames := onlyHostnamesContaining(allEndpoints, c.Hostname())
		logger.Debugf("/v1/endpoints filtered %v known endpoints and returned %v results", len(allEndpoints), len(matchedHostnames))
		return c.JSON(matchedHostnames)
	})

	app.Get("/v1/settings", func(c *fiber.Ctx) error {
		hostname := *flagHost
		colour := guessColour(logger, cfg.Colorschemes, hostname)
		logger.Debugf("guessed the colour: %s", colour)

		settings := Settings{
			ColorScheme: colour,
			Cluster:     hostname,
		}
		return c.JSON(settings)
	})

	app.Get("/healthz", func(c *fiber.Ctx) error {
		return c.SendString("ok")
	})

	app.Static("/", "./frontend/dist")

	onStartup(logger)

	logger.Info("starting webserver on :8000")
	logger.Fatal(app.Listen(":8000"))
}

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
