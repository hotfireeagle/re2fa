package route

import "github.com/gofiber/fiber/v2"

func RegisterRouter(app *fiber.App) {
	apiModule := app.Group("/api")

	apiModule.Get("/apiList", getApiList)
	apiModule.Post("/generateOriginFA", generateOriginFARoute)
	apiModule.Post("/match", faMatch)
}
