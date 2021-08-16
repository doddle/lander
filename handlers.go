package main

import (
	"fmt"
	"github.com/digtux/lander/pkg/deployments"
	"github.com/digtux/lander/pkg/endpoints"
	"github.com/digtux/lander/pkg/identicon/identicon"
	"github.com/digtux/lander/pkg/nodes"
	"github.com/digtux/lander/pkg/statefulsets"
	"github.com/gofiber/fiber/v2"
	"k8s.io/client-go/kubernetes"
	"os"
)

func getHealthz(c *fiber.Ctx) error {
	return c.SendString("ok")
}

func getSettings(c *fiber.Ctx) error {
	settings := Settings{
		ColorScheme: *flagColor,
		Cluster:     *flagHost,
		ClusterList: clusterList,
	}
	return c.JSON(settings)
}

func getDeployments(c *fiber.Ctx) error {
	clientSet, err := kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		logger.Fatal(err)
	}
	stats, err := deployments.AssembleDeploymentPieChart(logger, clientSet)
	if err != nil {
		logger.Fatal(err)
	}
	return c.JSON(stats)
}

func getStatefulSets(c *fiber.Ctx) error {
	clientSet, err := kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		logger.Fatal(err)
	}
	stats, err := statefulsets.AssembleStatefulSetPieChart(logger, clientSet)
	if err != nil {
		logger.Fatal(err)
	}
	return c.JSON(stats)
}

func getEndpoints(c *fiber.Ctx) error {
	clientSet, err := kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		logger.Fatal(err)
	}
	result := endpoints.ReallyAssemble(
		logger,
		clientSet,
		*flagIngressAnnotation,
	)
	return c.JSON(result)
}

// TODO: detect desired sizes from URI and generate smaller/bigger ones also
func getFavicon(c *fiber.Ctx) error {
	hex := *flagHex
	name := *flagHost
	uri := c.Context().Request.URI()
	uriPath := uri.LastPathSegment()

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
		err = f.Close()
		if err != nil {
			return err
		}
		logger.Infof("rendered a new icon for: %s, hex: %s", name, hex)
	}
	logger.Infof("%s served: %s", uriPath, fileName)
	return c.SendFile(fileName)
}

func getNodesPie(c *fiber.Ctx) error {
	clientSet, err := kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		logger.Fatal(err)
	}
	stats, err := nodes.AssembleNodesPieChart(logger, clientSet)
	if err != nil {
		logger.Fatal(err)
	}
	return c.JSON(stats)
}

func getNodesTable(c *fiber.Ctx) error {
	clientSet, err := kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		logger.Fatal(err)
	}
	stats, err := nodes.AssembleTable(logger, clientSet, nodeLabelSlice)
	if err != nil {
		logger.Fatal(err)
	}
	return c.JSON(stats)
}
