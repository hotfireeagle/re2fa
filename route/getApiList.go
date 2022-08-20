package route

import (
	"re2fa/model"

	"github.com/gofiber/fiber/v2"
)

func getApiList(ctx *fiber.Ctx) error {
	result := []*model.ApiListItem{
		{Name: "GenerateOriginFA", Api: "/api/generateOriginFA"},
	}

	res := &model.Response{
		Code: model.Success,
		Data: result,
	}

	return ctx.JSON(res)
}
