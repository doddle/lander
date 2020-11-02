package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"strconv"
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

type IconSize struct {
	Width  int
	Height int
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
		StrictRouting:         true, // tell fabric to not try to redirect to "//" all the time with static content
	}

	app := fiber.New(fiberCfg)

	app.Static("/", "./frontend/dist", fiber.Static{
		Compress: true,
		// MaxAge:   300,
	})

	app.Get("/favicon*", getFavicon)
	app.Get("/img/icons/*", getFavicon)

	app.Get("/healthz", getHealthz)
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

	// app.Use(cors.New(cors.Config{
	// 	AllowHeaders: "Cache-Control: No-Store",
	// }))

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

// guessSize will look at an inbound URI/path.. for example:
// - img/icons/apple-touch-icon-152x152.png
// - img/icons/apple-touch-icon-120x120.png
// it will attempt to match the 152x152 (or similar) patterns
// and return a string such as 152x152
// if none are matched a default of 250 is returned
func guessSize(input string) IconSize {

	// defaults
	var defaultSize = IconSize{}
	defaultSize.Width = 250
	defaultSize.Height = 250

	patternAll := `[0-9]{1,}x[0-9]{1,}`
	matched, err := regexp.MatchString(patternAll, input)

	// error compiling.. return the default
	if err != nil {
		return defaultSize
	}
	// didn't match a regex.. return the default
	if !matched {
		return defaultSize
	}

	// if it did match.. lets get the value from it
	re := regexp.MustCompile(patternAll)
	matchedString := re.FindString(input)
	x := strings.Split(matchedString, "x")

	desiredWidth := x[0]
	desiredHeight := x[1]
	if desiredHeight != desiredWidth {
		// didn't request identical Width + Height
		// this is intended for Identicons.. lets just return the default
		return defaultSize
	} else {
		w, _ := strconv.Atoi(desiredWidth)
		if w > 250 {
			return defaultSize
		}
		defaultSize.Width = w
		defaultSize.Height = w
		return defaultSize
	}
}

// TODO: detect desired sizes from URI and generate smaller/bigger ones also
func getFavicon(c *fiber.Ctx) error {
	hex := *flagHex
	uri := c.Context().Request.URI()
	name := uri.LastPathSegment()

	// buildPixelMap() has hard-coded sizes.. cannot use guessSize() until
	// thats figured out I guess
	// size := guessSize(string(uri.LastPathSegment()))
	size := IconSize{
		Width:  250,
		Height: 250,
	}

	fileName := fmt.Sprintf("/tmp/.lander-%dx%d-%s",
		size.Width,
		size.Height,
		name,
	)

	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		icon := identicon.Generate([]byte(name), hex, size.Width, size.Height)

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
	logger.Infof("%s served: %s", uri.Path(), fileName)
	return c.SendFile(fileName)
}
func getEndpoints(c *fiber.Ctx) error {
	// uri := c.Context().Request.URI()

	// get ALL endpoints
	allEndpoints := getIngressEndpoints(logger)

	// lets filter them for only ones matching the hostname of the context
	containsHostname := onlyHostnamesContaining(allEndpoints, *flagHost)

	// create an empty slice of endpoints
	matchedHostnames := []Endpoint{}
	for _, i := range containsHostname {
		// exclude hostnames that are not identical to flagHost
		// input here will be like: https://example.com, http://example.com/foo
		split := strings.Split(i.Address, "/")

		// len 3 includes only addresses with something after the "/" .. eg:
		// - example.com         <- not match
		// - example.com/foo     <- would match
		// this should exclude the lander itself (designed to be on https://cluster.example.com/)
		if len(split) > 3 {
			got := split[2]
			want := *flagHost
			if strings.Compare(got, want) == 0 {
				matchedHostnames = append(matchedHostnames, i)
			}
		}
	}
	// logger.Infof("%s returned %d of %d",
	// 	uri.Path(),
	// 	len(matchedHostnames),
	// 	len(allEndpoints),
	// )
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
