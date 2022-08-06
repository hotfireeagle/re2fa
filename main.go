package main

import (
	"fmt"
	"strings"
)

// 操作符优先级
var operatorPriority = map[byte]int{
	'*': 4,
	'?': 4,
	'+': 4,
	'.': 3,
	'|': 2,
	'(': 1, // TODO: 左右括号算操作符吗？没用的话去掉
}

var operatorStack *byteStack
var postfixResult strings.Builder

func re2postfix(re string) string {
	operatorStack = byteStackConstructor()
	postfixResult.Reset()

	shouldAddConcat := false

	// TODO: del continue change to if-else
	for i := 0; i < len(re); i++ {
		ch := re[i]

		if ch == '*' || ch == '?' || ch == '+' {
			shouldAddConcat = true
			pushOperator(ch)
			continue
		}

		if ch == '|' {
			shouldAddConcat = false
			pushOperator(ch)
			continue
		}

		if ch == '(' {
			if shouldAddConcat {
				pushOperator('.')
			}
			operatorStack.in(ch)
			shouldAddConcat = false
			continue
		}

		if ch == ')' {
			var operator byte

			for !operatorStack.isEmpty() {
				operator = operatorStack.out()
				if operator == '(' {
					break
				}
				postfixResult.WriteByte(operator)
			}

			if operator != '(' {
				panic("unmatched ')'")
			}

			shouldAddConcat = true
			continue
		}

		if shouldAddConcat {
			pushOperator('.')
		}

		postfixResult.WriteByte(ch)
		shouldAddConcat = true
	}

	for !operatorStack.isEmpty() {
		operator := operatorStack.out()
		if operator == '(' {
			panic("unmatched '('")
		}
		postfixResult.WriteByte(operator)
	}

	return postfixResult.String()
}

func pushOperator(operator byte) {
	currentPriority := operatorPriority[operator]

	for !operatorStack.isEmpty() {
		top := operatorStack.out()
		topPriority := operatorPriority[top]

		if topPriority >= currentPriority {
			postfixResult.WriteByte(top)
		} else {
			operatorStack.in(top)
			break
		}
	}

	operatorStack.in(operator)
}

type byteStack struct {
	vals []byte
}

func byteStackConstructor() *byteStack {
	return &byteStack{
		vals: make([]byte, 0),
	}
}

func (s *byteStack) in(val byte) {
	s.vals = append(s.vals, val)
}

func (s *byteStack) out() byte {
	val := s.vals[len(s.vals)-1]
	s.vals = s.vals[:len(s.vals)-1]
	return val
}

func (s *byteStack) peek() byte {
	return s.vals[len(s.vals)-1]
}

func (s *byteStack) isEmpty() bool {
	return len(s.vals) == 0
}

type stateStack struct {
	vals []*state
}

func stateStackConstructor() *stateStack {
	return &stateStack{
		vals: make([]*state, 0),
	}
}

func (s *stateStack) in(st *state) {
	s.vals = append(s.vals, st)
}

func (s *stateStack) out() *state {
	st := s.vals[len(s.vals)-1]
	s.vals = s.vals[:len(s.vals)-1]
	return st
}

var idCount int

type inputSymbol struct {
	character byte
	isEpsilon bool
}

func inputSymbolConstructor(character byte, isEpsilon bool) *inputSymbol {
	return &inputSymbol{
		character: character,
		isEpsilon: isEpsilon,
	}
}

type state struct {
	id          int
	transitions map[*inputSymbol][]*state
	end         *state
}

func stateConstructor() *state {
	idCount++
	return &state{
		id:          idCount,
		transitions: make(map[*inputSymbol][]*state),
	}
}

func (s *state) setEnd(end *state) {
	s.end = end
}

func (s *state) addTransition(input *inputSymbol, next *state) {
	if s.transitions[input] == nil {
		s.transitions[input] = make([]*state, 0)
	}
	s.transitions[input] = append(s.transitions[input], next)
}

type nfa struct {
	startState *state
	endState   *state
}

func re2nfaConstructor(regexp string) *nfa {
	n := &nfa{}

	postfix := re2postfix(regexp)

	n.post2nfa(postfix)

	return n
}

func (n *nfa) post2nfa(postfix string) {
	stateStack := stateStackConstructor()

	for i := 0; i < len(postfix); i++ {
		character := postfix[i]

		if character == '.' {
			rightState := stateStack.out()
			leftState := stateStack.out()

			epsilonSymbol := inputSymbolConstructor(0, true)
			leftState.addTransition(epsilonSymbol, rightState)

			stateStack.in(leftState)
		} else if character == '|' {
			rightState := stateStack.out()
			leftState := stateStack.out()

			newBegin := stateConstructor()
			newEnd := stateConstructor()
			newBegin.setEnd(newEnd)

			epsilonSymbol := inputSymbolConstructor(0, true)

			newBegin.addTransition(epsilonSymbol, leftState)
			newBegin.addTransition(epsilonSymbol, rightState)

			rightStateEnd := rightState.end
			leftStateEnd := leftState.end

			rightStateEnd.addTransition(epsilonSymbol, newEnd)
			leftStateEnd.addTransition(epsilonSymbol, newEnd)

			stateStack.in(newBegin)
		} else if character == '*' {
			state := stateStack.out()

			newBegin := stateConstructor()
			newEnd := stateConstructor()
			newBegin.setEnd(newEnd)

			epsilonSymbol := inputSymbolConstructor(0, true)

			stateEnd := state.end

			newBegin.addTransition(epsilonSymbol, state)
			stateEnd.addTransition(epsilonSymbol, state)
			stateEnd.addTransition(epsilonSymbol, newEnd)
			newBegin.addTransition(epsilonSymbol, newEnd)

			stateStack.in(newBegin)
		} else {
			begin := stateConstructor()
			end := stateConstructor()
			begin.setEnd(end)

			characterSymbol := inputSymbolConstructor(character, false)

			begin.addTransition(characterSymbol, end)

			stateStack.in(begin)
		}
	}

	state := stateStack.out()

	n.startState = state
	n.endState = state.end
}

func main() {
	regexp := "ab"
	n := re2nfaConstructor(regexp)
	fmt.Println(n.startState.transitions)
	fmt.Println(n.endState.id)
}
