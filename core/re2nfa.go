package core

import (
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
	// vals []*State
	vals []int
}

func stateStackConstructor() *stateStack {
	return &stateStack{
		vals: make([]int, 0),
		// vals: make([]*State, 0),
	}
}

func (s *stateStack) in(v int) {
	s.vals = append(s.vals, v)
}

func (s *stateStack) out() int {
	st := s.vals[len(s.vals)-1]
	s.vals = s.vals[:len(s.vals)-1]
	return st
}

type NFA struct {
	States        []int                    // 状态slice
	InputSymbols  []string                 // 输入inputSymbol
	TransitionMap map[int]map[string][]int // 状态转移函数
	StartStates   []int                    // 开始状态
	AcceptStates  []int                    // 接受状态
	BeginEndPairs map[int]int              // 一对状态的开始和结束
}

func NewNFA() *NFA {
	return &NFA{
		States:        make([]int, 0),
		InputSymbols:  make([]string, 0),
		TransitionMap: make(map[int]map[string][]int),
		StartStates:   make([]int, 0),
		AcceptStates:  make([]int, 0),
		BeginEndPairs: make(map[int]int),
	}
}

func (n *NFA) AddState() int {
	stateId := stateFactory()
	n.States = append(n.States, stateId)
	return stateId
}

func (n *NFA) AddInputSymbol(inputSymbol string) {
	n.InputSymbols = append(n.InputSymbols, inputSymbol)
}

func (n *NFA) AddTransition(inputSymbol string, fromStateId int, toStateId int) {
	if _, ok := n.TransitionMap[fromStateId]; !ok {
		n.TransitionMap[fromStateId] = make(map[string][]int)
	}

	if n.TransitionMap[fromStateId][inputSymbol] == nil {
		n.TransitionMap[fromStateId][inputSymbol] = make([]int, 0)
	}

	n.TransitionMap[fromStateId][inputSymbol] = append(n.TransitionMap[fromStateId][inputSymbol], toStateId)
}

func (n *NFA) SetBeginEndPairs(beginState int, endState int) {
	n.BeginEndPairs[beginState] = endState
}

func (n *NFA) GetEndState(beginState int) int {
	return n.BeginEndPairs[beginState]
}

func (n *NFA) AddStartState(stateId int) {
	n.StartStates = append(n.StartStates, stateId)
}

func (n *NFA) AddAcceptState(stateId int) {
	n.AcceptStates = append(n.AcceptStates, stateId)
}

var idCount int
var epsilonSymbol = "ε"

type State struct {
	Id          int                 `json:"id"`
	Transitions map[string][]*State `json:"transitions"`
	End         *State              `json:"end"`
}

func stateFactory() int {
	result := idCount
	idCount += 1
	return result
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

func Re2nfaConstructor(regexp string) *NFA {
	idCount = 0

	n := NewNFA()
	// n := &Nfa{}

	postfix := re2postfix(regexp)

	n.Postfix2NFA(postfix)

	return n
}

func (n *NFA) Postfix2NFA(postfix string) {
	stateStack := stateStackConstructor()

	for _, character := range postfix {
		if character == '.' {
			rightState := stateStack.out()
			leftState := stateStack.out()

			n.AddTransition(epsilonSymbol, n.GetEndState(leftState), rightState)
			n.SetBeginEndPairs(leftState, n.GetEndState(rightState))
			// epsilonSymbol := "-1"
			// leftState.End.addTransition(epsilonSymbol, rightState)
			// leftState.End = rightState.End

			stateStack.in(leftState)
		} else if character == '|' {
			rightState := stateStack.out()
			leftState := stateStack.out()

			newBegin := n.AddState()
			newEnd := n.AddState()

			n.SetBeginEndPairs(newBegin, newEnd)

			n.AddTransition(epsilonSymbol, newBegin, leftState)
			n.AddTransition(epsilonSymbol, newBegin, rightState)

			rightStateEnd := n.GetEndState(rightState)
			leftStateEnd := n.GetEndState(leftState)

			n.AddTransition(epsilonSymbol, rightStateEnd, newEnd)
			n.AddTransition(epsilonSymbol, leftStateEnd, newEnd)

			stateStack.in(newBegin)

			// newBegin := stateConstructor()
			// newEnd := stateConstructor()
			// newBegin.setEnd(newEnd)

			// epsilonSymbol := "-1"

			// newBegin.addTransition(epsilonSymbol, leftState)
			// newBegin.addTransition(epsilonSymbol, rightState)

			// rightStateEnd := rightState.End
			// leftStateEnd := leftState.End

			// rightStateEnd.addTransition(epsilonSymbol, newEnd)
			// leftStateEnd.addTransition(epsilonSymbol, newEnd)

			// stateStack.in(newBegin)
		} else if character == '*' {
			state := stateStack.out()

			newBegin := n.AddState()
			newEnd := n.AddState()
			n.SetBeginEndPairs(newBegin, newEnd)

			stateEnd := n.GetEndState(state)

			n.AddTransition(epsilonSymbol, newBegin, state)
			n.AddTransition(epsilonSymbol, stateEnd, state)
			n.AddTransition(epsilonSymbol, stateEnd, newEnd)
			n.AddTransition(epsilonSymbol, newBegin, newEnd)

			stateStack.in(newBegin)

			// newBegin := stateConstructor()
			// newEnd := stateConstructor()
			// newBegin.setEnd(newEnd)

			// epsilonSymbol := "-1"

			// stateEnd := state.End

			// newBegin.addTransition(epsilonSymbol, state)
			// stateEnd.addTransition(epsilonSymbol, state)
			// stateEnd.addTransition(epsilonSymbol, newEnd)
			// newBegin.addTransition(epsilonSymbol, newEnd)

			// stateStack.in(newBegin)
		} else {
			beginStateId := n.AddState()
			endStateId := n.AddState()

			n.SetBeginEndPairs(beginStateId, endStateId)
			stateStack.in(beginStateId)

			n.AddTransition(string(character), beginStateId, endStateId)
			// begin := stateConstructor()
			// end := stateConstructor()
			// begin.setEnd(end)

			// characterSymbol := string(character)

			// begin.addTransition(characterSymbol, end)

			// stateStack.in(begin)
		}
	}

	startState := stateStack.out()
	n.AddStartState(startState)
	n.AddAcceptState(n.GetEndState(startState))

	// state := stateStack.out()

	// n.StartState = state
	// n.EndState = state.End
}

// what's frontend need?
func (n *NFA) ConvertToJSON() *model.DrawNFAResponse {
	edges := make([]*model.Edge, 0)

	for startStateId, pathMap := range n.TransitionMap {
		for inputSymbol, endStates := range pathMap {
			for _, endStateId := range endStates {
				edge := &model.Edge{
					From:  startStateId,
					To:    endStateId,
					Label: inputSymbol,
				}
				edges = append(edges, edge)
			}
		}
	}

	return &model.DrawNFAResponse{
		Edges:        edges,
		Nodes:        n.States,
		StartStates:  n.StartStates,
		AcceptStates: n.AcceptStates,
	}
	// startState := n.StartState
	// endState := n.EndState

	// startId := startState.Id
	// endId := endState.Id

	// edges := make([]model.Edge, 0)
	// visited := make(map[string]bool) // key is fromStateId-inputSymbol-idx

	// var dfs func(s *State)

	// dfs = func(state *State) {
	// 	if state == nil {
	// 		return
	// 	}

	// 	from := state.Id

	// 	for inputSymbol, nextStates := range state.Transitions {
	// 		for idx, nextState := range nextStates {
	// 			to := nextState.Id
	// 			edgeIdx := fmt.Sprintf("%d-%s-%d", from, inputSymbol, idx)
	// 			if visited[edgeIdx] {
	// 				continue
	// 			}
	// 			visited[edgeIdx] = true

	// 			var label string
	// 			if inputSymbol == "-1" {
	// 				label = "ε"
	// 			} else {
	// 				label = inputSymbol
	// 			}
	// 			edges = append(edges, model.Edge{
	// 				From:  from,
	// 				To:    to,
	// 				Label: label,
	// 			})
	// 			dfs(nextState)
	// 		}
	// 	}
	// }

	// dfs(startState)

	// return &model.DrawNFAResponse{
	// 	Edges:       edges,
	// 	StartState:  startId,
	// 	AcceptState: endId,
	// }
}
