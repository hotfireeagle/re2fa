package route

import "github.com/gofiber/fiber/v2"

func RegisterRouter(app *fiber.App) {
	apiModule := app.Group("/api")
	apiModule.Post("/generateFA", generateFA)
	apiModule.Post("/generateFANoEpsilon", generateFANoEpsilonRoute)
}
