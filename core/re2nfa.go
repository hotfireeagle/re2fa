package core

import (
	"fmt"
	"re2fa/model"
	"strings"
)

// 操作符优先级
var operatorPriority = map[rune]int{
	'*': 4,
	'?': 4,
	'+': 4,
	'.': 3,
	'|': 2,
	'(': 1, // TODO: 左右括号算操作符吗？没用的话去掉
}

var operatorStack *runeStack
var postfixResult strings.Builder

func re2postfix(re string) string {
	operatorStack = runeStackConstructor()
	postfixResult.Reset()

	shouldAddConcat := false

	for _, ch := range re {
		if ch == '*' || ch == '?' || ch == '+' {
			shouldAddConcat = true
			pushOperator(ch)
		} else if ch == '|' {
			shouldAddConcat = false
			pushOperator(ch)
		} else if ch == '(' {
			if shouldAddConcat {
				pushOperator('.')
			}
			operatorStack.in(ch)
			shouldAddConcat = false
		} else if ch == ')' {
			var operator rune

			for !operatorStack.isEmpty() {
				operator = operatorStack.out()
				if operator == '(' {
					break
				}
				postfixResult.WriteRune(operator)
			}

			if operator != '(' {
				panic("unmatched ')'")
			}

			shouldAddConcat = true
		} else {
			if shouldAddConcat {
				pushOperator('.')
			}

			postfixResult.WriteRune(ch)
			shouldAddConcat = true
		}
	}

	for !operatorStack.isEmpty() {
		operator := operatorStack.out()
		if operator == '(' {
			panic("unmatched '('")
		}
		postfixResult.WriteRune(operator)
	}

	return postfixResult.String()
}

func pushOperator(operator rune) {
	currentPriority := operatorPriority[operator]

	for !operatorStack.isEmpty() {
		top := operatorStack.out()
		topPriority := operatorPriority[top]

		if topPriority >= currentPriority {
			postfixResult.WriteRune(top)
		} else {
			operatorStack.in(top)
			break
		}
	}

	operatorStack.in(operator)
}

type runeStack struct {
	vals []rune
}

func runeStackConstructor() *runeStack {
	return &runeStack{
		vals: make([]rune, 0),
	}
}

func (s *runeStack) in(val rune) {
	s.vals = append(s.vals, val)
}

func (s *runeStack) out() rune {
	val := s.vals[len(s.vals)-1]
	s.vals = s.vals[:len(s.vals)-1]
	return val
}

func (s *runeStack) isEmpty() bool {
	return len(s.vals) == 0
}

type stateStack struct {
	vals []*State
}

func stateStackConstructor() *stateStack {
	return &stateStack{
		vals: make([]*State, 0),
	}
}

func (s *stateStack) in(st *State) {
	s.vals = append(s.vals, st)
}

func (s *stateStack) out() *State {
	st := s.vals[len(s.vals)-1]
	s.vals = s.vals[:len(s.vals)-1]
	return st
}

var idCount int

type State struct {
	Id          int                 `json:"id"`
	Transitions map[string][]*State `json:"transitions"`
	End         *State              `json:"end"`
}

func stateConstructor() *State {
	state := &State{
		Id:          idCount,
		Transitions: make(map[string][]*State),
	}
	idCount += 1
	return state
}

func (s *State) setEnd(end *State) {
	s.End = end
}

func (s *State) addTransition(input string, next *State) {
	if s.Transitions[input] == nil {
		s.Transitions[input] = make([]*State, 0)
	}
	s.Transitions[input] = append(s.Transitions[input], next)
}

type Nfa struct {
	StartState *State `json:"startState"`
	EndState   *State `json:"endState"`
}

func Re2nfaConstructor(regexp string) *Nfa {
	idCount = 0
	n := &Nfa{}

	postfix := re2postfix(regexp)

	n.post2nfa(postfix)

	return n
}

func (n *Nfa) post2nfa(postfix string) {
	stateStack := stateStackConstructor()

	for _, character := range postfix {
		if character == '.' {
			rightState := stateStack.out()
			leftState := stateStack.out()

			epsilonSymbol := "-1"
			leftState.End.addTransition(epsilonSymbol, rightState)
			leftState.End = rightState.End

			stateStack.in(leftState)
		} else if character == '|' {
			rightState := stateStack.out()
			leftState := stateStack.out()

			newBegin := stateConstructor()
			newEnd := stateConstructor()
			newBegin.setEnd(newEnd)

			epsilonSymbol := "-1"

			newBegin.addTransition(epsilonSymbol, leftState)
			newBegin.addTransition(epsilonSymbol, rightState)

			rightStateEnd := rightState.End
			leftStateEnd := leftState.End

			rightStateEnd.addTransition(epsilonSymbol, newEnd)
			leftStateEnd.addTransition(epsilonSymbol, newEnd)

			stateStack.in(newBegin)
		} else if character == '*' {
			state := stateStack.out()

			newBegin := stateConstructor()
			newEnd := stateConstructor()
			newBegin.setEnd(newEnd)

			epsilonSymbol := "-1"

			stateEnd := state.End

			newBegin.addTransition(epsilonSymbol, state)
			stateEnd.addTransition(epsilonSymbol, state)
			stateEnd.addTransition(epsilonSymbol, newEnd)
			newBegin.addTransition(epsilonSymbol, newEnd)

			stateStack.in(newBegin)
		} else {
			begin := stateConstructor()
			end := stateConstructor()
			begin.setEnd(end)

			characterSymbol := string(character)

			begin.addTransition(characterSymbol, end)

			stateStack.in(begin)
		}
	}

	state := stateStack.out()

	n.StartState = state
	n.EndState = state.End
}

func (n *Nfa) ConvertToJSON() *model.DrawNFAResponse {
	startState := n.StartState
	endState := n.EndState

	startId := startState.Id
	endId := endState.Id

	edges := make([]model.Edge, 0)
	visited := make(map[string]bool) // key is fromStateId-inputSymbol-idx

	var dfs func(s *State)

	dfs = func(state *State) {
		if state == nil {
			return
		}

		from := state.Id

		for inputSymbol, nextStates := range state.Transitions {
			for idx, nextState := range nextStates {
				to := nextState.Id
				edgeIdx := fmt.Sprintf("%d-%s-%d", from, inputSymbol, idx)
				if visited[edgeIdx] {
					continue
				}
				visited[edgeIdx] = true

				var label string
				if inputSymbol == "-1" {
					label = "ε"
				} else {
					label = inputSymbol
				}
				edges = append(edges, model.Edge{
					From:  from,
					To:    to,
					Label: label,
				})
				dfs(nextState)
			}
		}
	}

	dfs(startState)

	return &model.DrawNFAResponse{
		Edges:       edges,
		StartState:  startId,
		AcceptState: endId,
	}
}