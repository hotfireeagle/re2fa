package service

import (
	"re2fa/core"
	"re2fa/model"
)

func GenerateOriginFAService(regexp string) []*model.FAItem {
	nfaObj := core.Re2nfaConstructor(regexp)
	dfaObj := core.NewDFAFromNFA(nfaObj)

	return []*model.FAItem{
		{FA: nfaObj, Title: "nfa"},
		{FA: dfaObj, Title: "dfa"},
	}
}
