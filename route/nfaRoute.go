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
		Data: nfaObj.ConvertToJSON(),
	}

	return ctx.JSON(&okRes)
}

func nfaMatch(ctx *fiber.Ctx) error {
	postObj := new(model.FAMatchPostData)

	if checkIsUnvalidJson(ctx, postObj) {
		return nil
	}

	if checkIsValidateFailed(ctx, postObj) {
		return nil
	}

	nfaObj := core.Re2nfaConstructor(postObj.RegExp)

	okRes := model.Response{
		Code: model.Success,
		Data: nfaObj.Match(postObj.Text),
	}

	return ctx.JSON(&okRes)
}
