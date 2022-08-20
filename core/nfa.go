package core

import (
	"re2fa/model"
	"strings"
)

var epsilonSymbol = "ε"

// 操作符优先级
var operatorPriority = map[rune]int{
	'*': 4,
	'?': 4,
	'+': 4,
	'.': 3,
	'|': 2,
	'(': 1,
}

func isNum(r byte) bool {
	return r >= '0' && r <= '9'
}

func isBigCharacter(r byte) bool {
	return r >= 'a' && r <= 'z'
}

func isSmallestCharacter(r byte) bool {
	return r >= 'A' && r <= 'Z'
}

// 应该可优化
func preConvert(str string) string {
	var result strings.Builder
	var needAddUnion bool = false

	var doNextWork = func(startPos int) int {
		var returnIdx int
		idx := startPos
		for idx < len(str) {
			ch := str[idx]
			if ch == '-' {
				if idx+1 < len(str) {
					beforeCharacter := str[idx-1]
					nextCharacter := str[idx+1]

					if isNum(beforeCharacter) && isNum(nextCharacter) {
						n1 := int(beforeCharacter - '0')
						n2 := int(nextCharacter - '0')

						for i := n1 + 1; i <= n2; i++ {
							result.WriteRune('|')
							result.WriteRune(rune(i + '0'))
						}

						idx += 2
					} else if isBigCharacter(beforeCharacter) && isBigCharacter(nextCharacter) {
						n1 := int(beforeCharacter - 'A')
						n2 := int(nextCharacter - 'A')

						for i := n1 + 1; i <= n2; i++ {
							result.WriteRune('|')
							result.WriteRune(rune(i + 'A'))
						}

						idx += 2
					} else if isSmallestCharacter(beforeCharacter) && isSmallestCharacter(nextCharacter) {
						n1 := int(beforeCharacter - 'a')
						n2 := int(nextCharacter - 'a')

						for i := n1 + 1; i <= n2; i++ {
							result.WriteRune('|')
							result.WriteRune(rune(i + 'a'))
						}

						idx += 2
					} else {
						if beforeCharacter == '[' && nextCharacter == ']' {
							result.WriteRune('-')
							idx += 1
						} else {
							panic("invalid character")
						}
					}
				} else {
					panic("invalid range")
				}
			} else if ch == ']' {
				result.WriteRune(')')
				returnIdx = idx + 1
				break
			} else {
				if needAddUnion {
					result.WriteRune('|')
				}
				result.WriteByte(ch)
				idx += 1
				needAddUnion = true
			}
		}
		return returnIdx
	}

	idx := 0
	for idx < len(str) {
		ch := str[idx]
		if ch == '[' {
			needAddUnion = false
			result.WriteRune('(')
			idx = doNextWork(idx + 1)
		} else {
			result.WriteByte(ch)
			idx += 1
		}
	}

	return result.String()
}

func re2postfix(re string) string {
	re2 := preConvert(re)

	operatorStack := runeStackConstructor()
	var postfixResult strings.Builder

	shouldAddConcat := false

	pushOperator := func(operator rune) {
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

	for _, ch := range re2 {
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
	vals []int
}

func stateStackConstructor() *stateStack {
	return &stateStack{
		vals: make([]int, 0),
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
	States              []int
	InputSymbols        []string                 // 输入inputSymbol
	TransitionMap       map[int]map[string][]int // 状态转移函数
	StartState          int                      // 开始状态
	AcceptStates        []int                    // 接受状态
	BeginEndPairs       map[int]int              // 一对状态的开始和结束
	IDCount             int                      // 状态计数器
	inputSymbolAddedMap map[string]bool
}

func NewNFA() *NFA {
	return &NFA{
		States:              make([]int, 0),
		InputSymbols:        make([]string, 0),
		TransitionMap:       make(map[int]map[string][]int),
		AcceptStates:        make([]int, 0),
		BeginEndPairs:       make(map[int]int),
		inputSymbolAddedMap: make(map[string]bool),
	}
}

func (n *NFA) AddState() int {
	stateId := n.IDCount
	n.States = append(n.States, stateId)
	n.IDCount += 1
	return stateId
}

func (n *NFA) AddInputSymbol(inputSymbol string) {
	if !n.inputSymbolAddedMap[inputSymbol] {
		n.InputSymbols = append(n.InputSymbols, inputSymbol)
		n.inputSymbolAddedMap[inputSymbol] = true
	}
}

func (n *NFA) AddTransition(inputSymbol string, fromStateId int, toStateId int) {
	n.AddInputSymbol(inputSymbol)
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

func (n *NFA) SetStartState(stateId int) {
	n.StartState = stateId
}

func (n *NFA) AddAcceptState(stateId int) {
	n.AcceptStates = append(n.AcceptStates, stateId)
}

func Re2nfaConstructor(regexp string) *NFA {
	n := NewNFA()

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
		} else if character == '?' {
			// zero or one time
			state := stateStack.out()

			newBegin := n.AddState()
			newEnd := n.AddState()
			n.SetBeginEndPairs(newBegin, newEnd)

			stateEnd := n.GetEndState(state)

			n.AddTransition(epsilonSymbol, newBegin, state)
			n.AddTransition(epsilonSymbol, stateEnd, newEnd)
			n.AddTransition(epsilonSymbol, newBegin, newEnd)

			stateStack.in(newBegin)
		} else if character == '+' {
			// one or more time
			state := stateStack.out()

			newBegin := n.AddState()
			newEnd := n.AddState()
			n.SetBeginEndPairs(newBegin, newEnd)

			stateEnd := n.GetEndState(state)

			n.AddTransition(epsilonSymbol, newBegin, state)
			n.AddTransition(epsilonSymbol, stateEnd, state)
			n.AddTransition(epsilonSymbol, stateEnd, newEnd)

			stateStack.in(newBegin)
		} else {
			beginStateId := n.AddState()
			endStateId := n.AddState()

			n.SetBeginEndPairs(beginStateId, endStateId)
			stateStack.in(beginStateId)

			n.AddTransition(string(character), beginStateId, endStateId)
		}
	}

	startState := stateStack.out()
	n.SetStartState(startState)
	n.AddAcceptState(n.GetEndState(startState))
}

// 判断一个state是否存在epsilon move
func (n *NFA) CheckStateHasEpsilonTransition(stateId int) bool {
	if _, ok := n.TransitionMap[stateId]; ok {
		if _, ok := n.TransitionMap[stateId][epsilonSymbol]; ok {
			return true
		} else {
			return false
		}
	} else {
		return false
	}
}

// 判断一个状态是否存在有效的转移
func (n *NFA) CheckStateHasValidTransition(stateId int) bool {
	if _, ok := n.TransitionMap[stateId]; ok {
		transitonMap := n.TransitionMap[stateId]
		result := false

		for inputSymbol, _ := range transitonMap {
			if inputSymbol != epsilonSymbol {
				result = true
				break
			}
		}

		return result
	} else {
		return false
	}
}

// 判断一个状态只存在epsilon transition
func (n *NFA) CheckStateJustHasEpsilonTransition(stateId int) bool {
	if _, ok := n.TransitionMap[stateId]; ok {
		inputSymbolLen := len(n.TransitionMap[stateId])
		if inputSymbolLen == 1 {
			if _, ok := n.TransitionMap[stateId][epsilonSymbol]; ok {
				return true
			} else {
				return false
			}
		} else {
			return false
		}
	} else {
		return false
	}
}

// 是自己所到达的state存在epsilon transition的话才可继续转移
func (n *NFA) Match(str string) bool {
	// TODO: 完成NFA的match
	dfaObj := NewDFAFromNFA(n)
	return dfaObj.Match(str)
}

func (n *NFA) ConvertToJSON() *model.DrawFAResponse {
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

	return &model.DrawFAResponse{
		Edges:        edges,
		Nodes:        n.States,
		StartState:   n.StartState,
		AcceptStates: n.AcceptStates,
		DeadState:    -1,
	}
}

// 返回一个NFA的逆NFA
// 如果regular language "abc"对应了NFA a
// 那么NFA a的逆NFA就可以用来表示cba
func (n *NFA) Reverse() *NFA {
	reverseNFA := &NFA{
		States:        n.States,
		InputSymbols:  n.InputSymbols,
		BeginEndPairs: n.BeginEndPairs,
		IDCount:       n.IDCount,
	}

	newAcceptStates := []int{
		reverseNFA.StartState,
	}

	reverseNFA.AcceptStates = newAcceptStates
	newBeginStateId := reverseNFA.AddState()
	reverseNFA.StartState = newBeginStateId

	reverseTransitionMap := make(map[int]map[string][]int)
	for originFromStateId, stateTransitions := range n.TransitionMap {
		for inputSymbol, originToStateIds := range stateTransitions {
			for _, originToStateId := range originToStateIds {
				if _, ok := reverseTransitionMap[originToStateId]; !ok {
					reverseTransitionMap[originToStateId] = make(map[string][]int)
				}
				if _, ok := reverseTransitionMap[originToStateId][inputSymbol]; !ok {
					reverseTransitionMap[originToStateId][inputSymbol] = make([]int, 0)
				}
				reverseTransitionMap[originToStateId][inputSymbol] = append(reverseTransitionMap[originToStateId][inputSymbol], originFromStateId)
			}
		}
	}

	// startState通过epsilon move转移到原来的endState
	reverseTransitionMap[newBeginStateId] = make(map[string][]int)
	reverseTransitionMap[newBeginStateId][epsilonSymbol] = n.AcceptStates

	reverseNFA.TransitionMap = reverseTransitionMap

	return reverseNFA
}
