package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/redirect/v2"
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

	// redirects to capture URL paths ending in "//"
	// .. this sometimes seems to happen when browsers arrive via their "back" buttons/history
	app.Use(redirect.New(redirect.Config{
		Rules: map[string]string{
			"//":  "/",
			"//*": "/$1",
		},
		StatusCode: 301,
	}))

	// start the static asset serving
	// NOTES:
	// - the fiber.Config has "StrictRouting" "true" which ensures unmatched-paths (EG "//")
	//   won't go to static assets
	// - ensure this is the last route declared
	app.Static("/", "./frontend/dist", fiber.Static{
		Compress: true,
		// MaxAge:   300,
	})
}
