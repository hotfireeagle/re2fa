package core

import "testing"

func TestNFAConcatenate(t *testing.T) {
	regexp := "abc"
	nfaObj := Re2nfaConstructor(regexp)

	if !nfaObj.Match("abc") {
		t.Error("abc should match")
	}

	if nfaObj.Match("abcd") {
		t.Error("abcd should't match")
	}
}

func TestNFAOr(t *testing.T) {
	regexp := "a|b"
	nfaObj := Re2nfaConstructor(regexp)

	if !nfaObj.Match("a") {
		t.Error("a should match")
	}

	if !nfaObj.Match("b") {
		t.Error("b should match")
	}

	if nfaObj.Match("ab") {
		t.Error("ab should't match")
	}

	if nfaObj.Match("c") {
		t.Error("c should't match")
	}
}

func TestNFAStar(t *testing.T) {
	regexp := "a*"

	nfaObj := Re2nfaConstructor(regexp)

	if nfaObj.Match("b") {
		t.Error("b should't match")
	}

	if !nfaObj.Match("a") {
		t.Error("a should match")
	}

	if !nfaObj.Match("aa") {
		t.Error("aa should match")
	}

	if nfaObj.Match("ab") {
		t.Error("ab should't match")
	}

	if nfaObj.Match("ba") {
		t.Error("ba should't match")
	}
}

func TestQuickOr(t *testing.T) {
	regexp := "[1-9]"
	nfaObj := Re2nfaConstructor(regexp)

	if !nfaObj.Match("1") {
		t.Error("nfa match error")
	}

	if !nfaObj.Match("5") {
		t.Error("5 should match")
	}

	if !nfaObj.Match("9") {
		t.Error("9 should match")
	}

	if nfaObj.Match("12") {
		t.Error("12 should't match")
	}
}

func TestNfaComplex(t *testing.T) {
	regexp := "12[1-9]sdjk"
	nfaObj := Re2nfaConstructor(regexp)

	if !nfaObj.Match("121sdjk") {
		t.Error("nfa match error")
	}

	if nfaObj.Match("121") {
		t.Error("121 should't match")
	}
}
