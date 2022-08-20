package core

import "testing"

func TestSuffix(t *testing.T) {
	regexp := "abc"
	nfaObj := Re2nfaConstructor(regexp)

	dfaObj := NewDFAFromNFA(nfaObj)

	nfa := dfaObj.Suffix()

	if !nfa.Match("abc") {
		t.Error("abc should match")
	}

	if !nfa.Match("bc") {
		t.Error("bc should match")
	}

	if !nfa.Match("c") {
		t.Error("c should match")
	}

	if nfa.Match("ac") {
		t.Error("ac should't match")
	}
}
