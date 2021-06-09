package main

import (
	"github.com/gofiber/fiber/v2"
)

func startRoutes(app *fiber.App) {

	app.Get("/favicon*", getFavicon)
	app.Get("/img/icons/*", getFavicon)

	app.Get("/healthz", getHealthz)
	app.Get("/v1/endpoints", getEndpoints)
	app.Get("/v1/settings", getSettings)

	app.Get("/v1/pie/deployments", getDeployments)
	app.Get("/v1/pie/statefulsets", getStatefulSets)
	app.Get("/v1/pie/nodes", getNodesPie)
	app.Get("/v1/table/nodes", getNodesTable)
}
