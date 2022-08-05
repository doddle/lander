package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/doddle/lander/pkg/deployments"
	"github.com/doddle/lander/pkg/endpoints"
	"github.com/doddle/lander/pkg/identicon/identicon"
	"github.com/doddle/lander/pkg/inventory"
	"github.com/doddle/lander/pkg/nodes"
	"github.com/doddle/lander/pkg/statefulsets"
	"github.com/gofiber/fiber/v2"
	"k8s.io/client-go/kubernetes"
)

func getHealthz(c *fiber.Ctx) error {
	return c.SendString("ok")
}

func getSettings(c *fiber.Ctx) error {
	settings := Settings{
		ColorScheme: *flagColor,
		Cluster:     *flagClusterFQDN,
		ClusterList: clusterList,
	}
	return c.JSON(settings)
}

func getDeployments(c *fiber.Ctx) error {
	clientSet, err := kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		logger.Fatal(err)
	}
	resp, err := deployments.AssembleDeploymentPieChart(logger, clientSet)
	if err != nil {
		logger.Fatal(err)
	}
	return c.JSON(resp)
}

func getStatefulSets(c *fiber.Ctx) error {
	clientSet, err := kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		logger.Fatal(err)
	}
	resp, err := statefulsets.AssembleStatefulSetPieChart(logger, clientSet)
	if err != nil {
		logger.Fatal(err)
	}
	return c.JSON(resp)
}

func getClusterLinks(c *fiber.Ctx) error {
	clientSet, err := kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		logger.Fatal(err)
	}
	resp := endpoints.ReallyAssemble(
		logger,
		clientSet,
		*flagLanderAnnotationBase,
	)

	// if there's no Icon set, try and find one by name match
	for i, endpoint := range resp {
		if endpoint.Icon == "" {
			resp[i].Icon = "link.png"
			for _, iconFile := range availableIcons {
				if strings.Contains(iconFile, strings.ToLower(endpoint.Name)) {
					resp[i].Icon = iconFile
				}
			}
		}
	}
	if len(resp) < 1 {
		return c.JSON([]string{})
	}
	return c.JSON(resp)
}

func getDeploymentsTable(c *fiber.Ctx) error {
	clientSet, err := kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		logger.Fatal(err)
	}
	resp := deployments.AssembleDeploymentsTable(
		logger,
		clientSet,
		filteredTags,
	)
	if len(resp) < 1 {
		return c.JSON([]string{})
	}
	return c.JSON(resp)
}

func getRoutes(c *fiber.Ctx) error {
	clientSet, err := kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		logger.Fatal(err)
	}
	resp := endpoints.AssembleRouteMetaData(
		logger,
		clientSet,
	)
	if len(resp) < 1 {
		return c.JSON([]string{})
	}
	return c.JSON(resp)
}

func getStatefulSetsTable(c *fiber.Ctx) error {
	clientSet, err := kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		logger.Fatal(err)
	}

	resp := statefulsets.AssembleDeploymentsTable(
		logger,
		clientSet,
	)
	if len(resp) < 1 {
		return c.JSON([]string{})
	}
	return c.JSON(resp)
}

// TODO: detect desired sizes from URI and generate smaller/bigger ones also
func getFavicon(c *fiber.Ctx) error {
	hex := *flagHex
	name := *flagClusterFQDN
	uri := c.Context().Request.URI()
	uriPath := uri.LastPathSegment()

	// buildPixelMap() has hard-coded sizes. cannot use guessSize() until
	// that's figured out I guess
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
		err = f.Close()
		if err != nil {
			return err
		}
		logger.Infof("rendered a new icon for: %s, hex: %s", name, hex)
	}
	logger.Debugf("%s served: %s", uriPath, fileName)
	return c.SendFile(fileName)
}

func getNodesPie(c *fiber.Ctx) error {
	clientSet, err := kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		logger.Fatal(err)
	}
	resp, err := nodes.AssembleNodesPieChart(logger, clientSet)
	if err != nil {
		logger.Fatal(err)
	}
	return c.JSON(resp)
}

func getNodesTable(c *fiber.Ctx) error {
	clientSet, err := kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		logger.Fatal(err)
	}
	resp, err := nodes.AssembleTable(logger, clientSet, nodeLabelSlice)
	if err != nil {
		logger.Fatal(err)
	}
	return c.JSON(resp)
}

func getFluxIgnored(c *fiber.Ctx) error {
	data := inventory.AssembleFluxIgnored(logger, kubeConfig)
	// if the data received is empty, at least return an empty nodeLabelSlice
	// c.JSON will convert an empty slice into a valid `[]`
	// ... if we don't do this the endpoint returns `null` which isn't valid json
	if len(*data) < 1 {
		return c.JSON([]string{})
	}
	return c.JSON(data)
}
