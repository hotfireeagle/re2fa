package main

import (
	"re2fa/core"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	regexp := "ab"

	n := core.Re2nfaConstructor(regexp)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(n)
	})

	app.Listen(":3000")
}
