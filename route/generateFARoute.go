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

	faItems := service.GenerateOriginFAService(obj.RegExp)

	return ctx.JSON(getOkResponseFromFaItems(faItems))
}

func generateNFAAndSuffixNFARoute(ctx *fiber.Ctx) error {
	obj := new(model.GenerateFAPostData)

	if checkIsUnvalidJson(ctx, obj) {
		return nil
	}

	if checkIsValidateFailed(ctx, obj) {
		return nil
	}

	faItems := service.GenerateNfaAndSuffixNfa(obj.RegExp)

	return ctx.JSON(getOkResponseFromFaItems(faItems))
}
