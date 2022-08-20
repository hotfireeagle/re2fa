package route

import (
	"re2fa/model"
	"re2fa/service"

	"github.com/gofiber/fiber/v2"
)

func generateOriginFARoute(ctx *fiber.Ctx) error {
	obj := new(model.GenerateFAPostData)

	if checkIsUnvalidJson(ctx, obj) {
		return nil
	}

	if checkIsValidateFailed(ctx, obj) {
		return nil
	}

	response := service.GenerateOriginFAService(obj)

	return ctx.JSON(response)
}
