package main

import (
	"flag"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/redirect/v2"
	"github.com/icza/gox/imagex/colorx"
	"github.com/withmandala/go-log"
	"os"
	"strings"
)

var (
	// setup some global vars
	flagHost = flag.String("host", "k8s.prod.example.com", "filter ingresses matching this hostname [required]")
	// flagConfig = flag.String("config", "default", "Specify a config file (customised colour scheme)")
	flagColor = flag.String("color", "light-blue lighten-2", "Main color scheme (See: https://vuetifyjs.com/en/styles/colors/#material-colors)")
	flagHex   = flag.String("hex", "#26c5e8", "identicon color, hex string, eg #112233, #123, #bAC")

	// flag for a list of all clusters
	flagClusters = flag.String("clusters", "cluster1.example.com,cluster2.example.com", "comma seperated list of clusters")

	flagNodeLabels = flag.String(
		"node-labels",
		"kubernetes.io/role,node.kubernetes.io/instance-type,node.kubernetes.io/instancegroup,topology.kubernetes.io/zone",
		"comma seperated list of node labels you care about",
	)

	flagDebug      = flag.Bool("debug", false, "debug")
	clusterList    []string
	nodeLabelSlice []string

	// TODO: ideally the logger shouldn't be global
	logger     *log.Logger
	kubeConfig = autoClientInit(logger)
)

func init() {
	flag.Parse()
	clusterList = strings.Split(*flagClusters, ",")
	nodeLabelSlice = strings.Split(*flagNodeLabels, ",")
	logger = newLogger(*flagDebug)
}

// Settings to be returned to the browser/client
type Settings struct {
	ColorScheme string   `json:"colorscheme"`
	Cluster     string   `json:"cluster"`
	ClusterList []string `json:"clusters"`
}

type IconSize struct {
	Width  int
	Height int
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
		StrictRouting:         true, // tell fiber to not try to redirect to "//" all the time with static content
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

	app.Get("/v1/pie/deployments", getDeployments)
	app.Get("/v1/pie/statefulsets", getStatefulSets)
	app.Get("/v1/pie/nodes", getNodesPie)
	app.Get("/v1/table/nodes", getNodesTable)

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
