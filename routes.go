package main

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/redirect/v2"
)

func startRoutes(app *fiber.App) {

	app.Get("/favicon*", getFavicon)
	app.Get("/img/icons/*", getFavicon)

	app.Get("/healthz", getHealthz)
	app.Get("/v1/endpoints", getEndpoints)
	app.Get("/v1/routes", getRoutes)
	app.Get("/v1/table/deployments", getDeploymentsTable)
	app.Get("/v1/settings", getSettings)

	app.Get("/v1/pie/deployments", getDeployments)
	app.Get("/v1/pie/statefulsets", getStatefulSets)
	app.Get("/v1/pie/nodes", getNodesPie)
	app.Get("/v1/table/nodes", getNodesTable)

	// redirects to capture URL paths ending in "//"
	// .. this sometimes seems to happen when browsers arrive via their "back" buttons/history

	app.Get("//", func(c *fiber.Ctx) error {
		return c.Redirect("/")
	})

	// match incoming queries prefixed with "//" and suggest a redirect
	// EG a GET like "//foo" will get a 301 -> "/foo"
	app.Use(redirect.New(redirect.Config{
		Rules: map[string]string{
			"//*": "/$1",
		},
		StatusCode: 301,
	}))

	// start the static asset serving
	// NOTES:
	// - the fiber.Config has "StrictRouting" "true" which ensures unmatched-paths (EG "//")
	//   won't go to static assets
	// - ensure this is the last route declared
	//
	//  read about fiber static settings: https://docs.gofiber.io/api/app#static

	// Cache-Control tweaks
	// read more https://developers.google.com/web/fundamentals/performance/optimizing-content-efficiency/http-caching#defining_optimal_cache-control_policy
	app.Static("/assets", "./frontend/dist/assets", fiber.Static{
		Compress:      true,
		MaxAge:        3600,
		CacheDuration: 60 * time.Minute,
	})
	app.Static("/css", "./frontend/dist/css", fiber.Static{
		Compress:      true,
		MaxAge:        3600,
		CacheDuration: 60 * time.Minute,
	})

	// potentially the js should be cached for less as its changed when we deploy new versions of lander
	app.Static("/js", "./frontend/dist/js", fiber.Static{
		Compress:      true,
		MaxAge:        300,
		CacheDuration: 5 * time.Minute,
	})

	// NEVER cache /index.html
	// currently fiber doesn't allow setting custom headers (EG: "no-cache", "no-store" or "expires 0")
	// we can set "expires 1" and "max-age=1" however
	app.Static("/", "./frontend/dist", fiber.Static{
		Compress:      true,
		CacheDuration: 1 * time.Second,
		MaxAge:        1,
	})
}
