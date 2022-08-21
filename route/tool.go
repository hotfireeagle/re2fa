package route

import (
	"re2fa/model"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func checkIsUnvalidJson(ctx *fiber.Ctx, data interface{}) bool {
	if err := ctx.BodyParser(data); err != nil {
		res := model.Response{Code: model.Error, Msg: "无效JSON", ErrorLog: err.Error()}
		ctx.JSON(&res)
		return true
	}
	return false
}

func checkIsValidateFailed(ctx *fiber.Ctx, data interface{}) bool {
	validate := validator.New()
	if err := validate.Struct(data); err != nil {
		res := model.Response{Code: model.Error, Msg: "无效参数", ErrorLog: err.Error()}
		ctx.JSON(&res)
		return true
	}
	return false
}

func getOkResponseFromFaItems(faItems []*model.FAItem) *model.Response {
	var result = make([]*model.DrawFAResponse, 0)
	for _, item := range faItems {
		uiItem := item.FA.ConvertToJSON()
		uiItem.Title = item.Title
		result = append(result, uiItem)
	}
	return &model.Response{
		Code: model.Success,
		Data: result,
	}
}
