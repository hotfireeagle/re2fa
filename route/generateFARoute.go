package route

import (
	"re2fa/core"
	"re2fa/model"

	"github.com/gofiber/fiber/v2"
)

func generateFA(ctx *fiber.Ctx) error {
	obj := new(model.GenerateFAPostData)

	if checkIsUnvalidJson(ctx, obj) {
		return nil
	}

	if checkIsValidateFailed(ctx, obj) {
		return nil
	}

	nfaObj := core.Re2nfaConstructor(obj.RegExp)
	okRes := model.Response{
		Code: model.Success,
		Data: nfaObj,
	}

	return ctx.JSON(&okRes)
}
