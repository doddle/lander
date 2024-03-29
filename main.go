package main

import (
	"flag"
	"os"
	"strings"

	"github.com/doddle/lander/pkg/deployments"
	"github.com/gofiber/fiber/v2"
	"github.com/icza/gox/imagex/colorx"
	"github.com/withmandala/go-log"
)

var (
	// setup some global vars
	flagClusterFQDN = flag.String("clusterFQDN", "k8s.example.com", "The cluster this lander is operating in. Used for display and identicon purposes only.")
	// flagConfig = flag.String("config", "default", "Specify a config file (customised colour scheme)")
	flagColor                = flag.String("color", "light-blue lighten-2", "Main color scheme (See: https://vuetifyjs.com/en/styles/colors/#material-colors)")
	flagHex                  = flag.String("hex", "#26c5e8", "identicon color, hex string, eg #112233, #123, #bAC")
	flagLanderAnnotationBase = flag.String("annotationBase", "lander.doddle.tech", "The base of the annotations used for lander. e.g. lander.doddle.tech for annotations like lander.doddle.tech/show")

	// flag for a list of all clusters
	flagClusters = flag.String("clusters", "cluster1.example.com,cluster2.example.com", "comma seperated list of clusters")

	flagShowTagsFor = flag.String("show-tags-for", "ecr=1122334455.dkr.ecr.eu-west-1.amazonaws.com/acmecorp/,gcr=blabla", "show tags for images matching this string in the Deployment tables")

	flagNodeLabels = flag.String(
		"labels",
		"kubernetes.io/role,node.kubernetes.io/instance-type,kops.k8s.io/instancegroup,topology.kubernetes.io/zone",
		"comma seperated list of node labels you care about",
	)

	availableIcons []string
	flagDebug      = flag.Bool("debug", false, "debug")
	clusterList    []string
	nodeLabelSlice []string
	filteredTags   []deployments.TagFilters

	// TODO: ideally the logger shouldn't be global
	logger     *log.Logger
	kubeConfig = autoClientInit(logger)
)

func init() {
	flag.Parse()
	filteredTags = deployments.ParseShowTagsFor(*flagShowTagsFor)
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

	// app.Use(cors.New(cors.Config{
	// 	AllowHeaders: "Cache-Control: No-Store",
	// }))

	onStartup(logger)

	startRoutes(app)

	logger.Info("starting webserver on http://0.0.0.0:8000")
	logger.Fatal(app.Listen(":8000"))
}
