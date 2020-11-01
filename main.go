package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/digtux/lander/identicon"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/redirect/v2"
	"github.com/icza/gox/imagex/colorx"
	"github.com/patrickmn/go-cache"
	"github.com/withmandala/go-log"
)

var (
	// create a short cache.. to prevent us hammering the kube api
	cacheShort = cache.New(1*time.Minute, 2*time.Minute)

	// setup some global vars
	flagHost = flag.String("host", "k8s.prod.example.com", "filter ingresses matching this hostname [required]")
	// flagConfig = flag.String("config", "default", "Specify a config file (customised colour scheme)")
	flagColor = flag.String("color", "light-blue lighten-2", "Main color scheme (See: https://vuetifyjs.com/en/styles/colors/#material-colors)")
	flagHex   = flag.String("hex", "#26c5e8", "identicon color, hex string, eg #112233, #123, #bAC")

	// TODO: ideally the logger shouldn't be global
	logger = newLogger(false)
)

func init() {
	flag.Parse()
}

// Endpoint is for the metadata returned to the browser/frontend
type Endpoint struct {
	// these are not ac
	Address     string `json:"address"`     // combined hostname + /paths
	Class       string `json:"class"`       // name of the ingress path
	Https       bool   `json:"https"`       // https
	Oauth2proxy bool   `json:"oauth2proxy"` // if its secured by an oauth2proxy
	Icon        string `json:"icon"`        // we will attempt to guess the ICON for endpoints
	Description string `json:"description"` // if we can match this to an app, we can propogate this
	Name        string `json:"name"`        // Application name
}

// EndpointList is just a list/slice of Endpoint objects to make it easier to work with them
type EndpointList struct {
	Endpoints []Endpoint `json:"endpoints"`
}

// Settings to be returned to the browser/client
type Settings struct {
	ColorScheme string `json:"colorscheme"`
	Cluster     string `json:"cluster"`
}

func onStartup(logger *log.Logger) {
	logger.Info("getting some initial data bootstrapped")
	_ = getIngressEndpoints(logger)
}

func main() {

	checkRequredFlag()

	// handy hex parser: https://github.com/icza/gox/blob/7dc3510ae515f0a6e8479d9a382bc8bb04f3a37d/imagex/colorx/colorx_test.go#L10-L14
	_, err := colorx.ParseHexColor(*flagHex)
	if err != nil {
		logger.Errorf("unable to parse hex value: %s", *flagHex)
		os.Exit(1)
	}

	logger.Infof("hex: %s", *flagHex)
	logger.Infof("colorscheme: %s", *flagColor)

	fiberCfg := fiber.Config{
		DisableStartupMessage: true,
	}

	app := fiber.New(fiberCfg)

	app.Static("/", "./frontend/dist")

	app.Get("/favicon*", getFavicon)
	app.Get("/healthz", getHealthz)
	app.Get("/img/icons/*", getFavicon)
	app.Get("/v1/endpoints", getEndpoints)
	app.Get("/v1/settings", getSettings)

	// sometimes in firefox (pressing "back") you can end up with the url example.com//
	// redirect that back

	app.Use(redirect.New(redirect.Config{
		Rules: map[string]string{
			"//":  "/",
			"//*": "/",
			// "//*": "/new/$1",
		},
		StatusCode: 301,
	}))

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

func getSettings(c *fiber.Ctx) error {
	settings := Settings{
		ColorScheme: *flagColor,
		Cluster:     *flagHost,
	}
	return c.JSON(settings)
}

func getHealthz(c *fiber.Ctx) error {
	return c.SendString("ok")
}

// TODO: detect desired sizes from URI and generate smaller/bigger ones also
func getFavicon(c *fiber.Ctx) error {
	hex := *flagHex
	name := *flagHost

	fileName := fmt.Sprintf("/tmp/%s.png", name)

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
		logger.Infof("rendered a new icon for: %s, hex: %s", name, hex)
	}
	uri := c.Context().Request.URI().LastPathSegment()
	logger.Infof("/%s served: %s", uri, fileName)
	return c.SendFile(fileName)
}
func getEndpoints(c *fiber.Ctx) error {
	// get ALL endpoints
	allEndpoints := getIngressEndpoints(logger)
	// lets filter them for only ones matching the hostname of the context
	containsHostname := onlyHostnamesContaining(allEndpoints, *flagHost)
	matchedHostnames := []Endpoint{}
	for _, i := range containsHostname {
		// exclude hostnames that are not identical to flagHost
		// input here will be like: https://example.com, http://example.com/foo
		split := strings.Split(i.Address, "/")

		// len 3 includes only addresses with something after the "/" .. eg:
		// - example.com/foo
		// this should exclude the lander itself (designed to be on https://cluster.example.com)
		if len(split) > 3 {
			got := split[2]
			want := *flagHost
			if strings.Compare(got, want) == 0 {
				matchedHostnames = append(matchedHostnames, i)

			}
		}
	}
	logger.Infof("/v1/endpoints filtered %v known endpoints and returned %v results", len(allEndpoints), len(matchedHostnames))
	return c.JSON(matchedHostnames)
}

// using flag.Visit, check if a flag was provided
// if not.. tell the user, print `-help` and bail
func checkRequredFlag() {
	required := []string{"host"}
	seen := make(map[string]bool)
	flag.Visit(func(f *flag.Flag) { seen[f.Name] = true })
	for _, req := range required {
		if !seen[req] {
			fmt.Printf("ERROR: missing required flag: '-%s'\n-------------\n", req)
			flag.PrintDefaults()
			os.Exit(2)
		}
	}
}
