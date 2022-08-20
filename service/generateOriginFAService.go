package service

import (
	"re2fa/core"
	"re2fa/model"
)

func GenerateOriginFAService(obj *model.GenerateFAPostData) *model.Response {
	nfaObj := core.Re2nfaConstructor(obj.RegExp)
	dfaObj := core.NewDFAFromNFA(nfaObj)

	return &model.Response{
		Code: model.Success,
		Data: model.NewDoubleFA(dfaObj.ConvertToJSON(), "dfa", nfaObj.ConvertToJSON(), "nfa"),
	}
}
