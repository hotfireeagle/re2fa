package main

import (
	"re2fa/route"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	route.RegisterRouter(app)

	err := app.Listen(":3000")

	if err != nil {
		panic(err)
	}
}
