package route

import (
	"re2fa/model"
	"re2fa/service"
	"time"

	"github.com/gofiber/fiber/v2"
)

func faMatch(ctx *fiber.Ctx) error {
	postObj := new(model.FAMatchPostData)

	if checkIsUnvalidJson(ctx, postObj) {
		return nil
	}

	if checkIsValidateFailed(ctx, postObj) {
		return nil
	}

	var faItems []*model.FAItem

	api := postObj.Api

	if api == "/api/generateOriginFA" {
		faItems = service.GenerateOriginFAService(postObj.RegExp)
	} else if api == "/api/generateNFAAndSuffixNFA" {
		faItems = service.GenerateNfaAndSuffixNfa(postObj.RegExp)
	}

	answers := make([]model.MatchAnswer, 0)

	for _, faItem := range faItems {
		startTime := time.Now().UnixNano()
		matchResult := faItem.FA.Match(postObj.Text)
		endTime := time.Now().UnixNano()

		item := model.MatchAnswer{
			Result: matchResult,
			Time:   float64((endTime - startTime) / 1e3),
		}

		answers = append(answers, item)
	}

	okRes := model.Response{
		Code: model.Success,
		Data: answers,
	}

	return ctx.JSON(&okRes)
}
