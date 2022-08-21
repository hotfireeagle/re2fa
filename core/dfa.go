package core

import (
	"re2fa/model"
	"sort"
	"strconv"
	"strings"
)

type DFA struct {
	States          []int
	InputSymbols    []string
	TransitionMap   map[int]map[string]int
	StartState      int
	AcceptStates    []int
	StateIdToSetMap map[string]int
	DeadStateId     int
	dfaStatesCount  int
}

func NewDFA() *DFA {
	return &DFA{
		States:          make([]int, 0),
		AcceptStates:    make([]int, 0),
		InputSymbols:    make([]string, 0),
		TransitionMap:   make(map[int]map[string]int),
		StateIdToSetMap: make(map[string]int),
		DeadStateId:     -1,
	}
}

func (d *DFA) AddTransition(fromStateId int, toStateId int, inputSymbol string) {
	if _, ok := d.TransitionMap[fromStateId]; !ok {
		d.TransitionMap[fromStateId] = make(map[string]int)
	}
	d.TransitionMap[fromStateId][inputSymbol] = toStateId
}

func (d *DFA) AddState(nfaStates []int) (int, bool) {
	if len(nfaStates) == 0 {
		return d.AddDeadState()
	}

	idStr := slice2str(nfaStates)
	if _, ok := d.StateIdToSetMap[idStr]; ok {
		return d.StateIdToSetMap[idStr], true
	}

	stateId := d.dfaStatesCount
	d.States = append(d.States, stateId)
	d.AddStateIdToSetMap(stateId, idStr)

	d.dfaStatesCount += 1
	return stateId, false
}

func (d *DFA) AddDeadState() (int, bool) {
	sId := d.dfaStatesCount

	if d.DeadStateId == -1 {
		d.DeadStateId = sId
		d.States = append(d.States, sId)
		d.AddTransitionForDeadState()
		d.dfaStatesCount += 1
		return sId, false
	} else {
		return d.DeadStateId, true
	}
}

func (d *DFA) AddTransitionForDeadState() {
	for _, inputSymbol := range d.InputSymbols {
		d.AddTransition(d.DeadStateId, d.DeadStateId, inputSymbol)
	}
}

func (d *DFA) AddStateIdToSetMap(stateId int, set string) {
	d.StateIdToSetMap[set] = stateId
}

func (d *DFA) GetStateIdByStr(stateStr string) int {
	if stateStr == "" {
		return d.DeadStateId
	}
	return d.StateIdToSetMap[stateStr]
}

func (d *DFA) SetStartState(stateId int) {
	d.StartState = stateId
}

func (d *DFA) SetInputSymbols(inputSymbols []string) {
	ism := make([]string, 0)
	for _, v := range inputSymbols {
		if v == epsilonSymbol {
			continue
		}
		ism = append(ism, v)
	}
	d.InputSymbols = ism
}

func NewDFAFromNFA(n *NFA) *DFA {
	dfa := NewDFA()

	dfa.SetInputSymbols(n.InputSymbols)

	var findCurrentStateCanGoAnyStateByEpsilon = func(state int) []int {
		canGoStates := make([]int, 0)
		visited := make(map[int]bool)

		var dfs func(s int)

		dfs = func(currentState int) {
			if visited[currentState] {
				return
			}
			visited[currentState] = true

			transitions := n.TransitionMap[currentState]
			for inputSymbol, toStates := range transitions {
				if inputSymbol == epsilonSymbol {
					for _, stateId := range toStates {
						canGoStates = append(canGoStates, stateId)
						dfs(stateId)
					}
				}
			}
		}

		dfs(state)

		return canGoStates
	}

	// 开始节点state
	startStates := findCurrentStateCanGoAnyStateByEpsilon(n.StartState)
	startStates = append(startStates, n.StartState)

	startStateId, _ := dfa.AddState(startStates)
	dfa.SetStartState(startStateId)

	needBeSettle := make([][]int, 0)
	needBeSettle = append(needBeSettle, startStates)

	for len(needBeSettle) > 0 {
		nextNeedBeSettle := make([][]int, 0)

		for _, states := range needBeSettle {

			for _, inputSymbol := range dfa.InputSymbols {

				nextCanGoStates := make(map[int]bool)

				for _, fromStateId := range states {

					canGoStateList := n.TransitionMap[fromStateId][inputSymbol]
					for _, canGoStateId := range canGoStateList {
						nextCanGoStates[canGoStateId] = true
						// we also need find the epsilon move
						thisStateCanGoByEpsilonMove := findCurrentStateCanGoAnyStateByEpsilon(canGoStateId)
						for _, id := range thisStateCanGoByEpsilonMove {
							nextCanGoStates[id] = true
						}
					}

				}

				nextCanGoStateIds := getKeys(nextCanGoStates)
				dfaStateId, hasExist := dfa.AddState(nextCanGoStateIds)

				fromStateStr := slice2str(states)
				dfaFromStateId := dfa.GetStateIdByStr(fromStateStr)
				dfa.AddTransition(dfaFromStateId, dfaStateId, inputSymbol)

				if !hasExist {
					nextNeedBeSettle = append(nextNeedBeSettle, nextCanGoStateIds)
				}
			}
		}

		needBeSettle = nextNeedBeSettle
	}

	endStates := make([]int, 0)
	for stateStr, stateId := range dfa.StateIdToSetMap {
		nfaFinalState := n.AcceptStates[0]
		if strings.Contains(stateStr, strconv.Itoa(nfaFinalState)) {
			endStates = append(endStates, stateId)
		}
	}

	dfa.AcceptStates = endStates

	return dfa
}

func (d *DFA) Match(str string) bool {
	currentState := d.StartState

	for _, character := range str {
		currentStateTransitions := d.TransitionMap[currentState]
		currentState = -1
		for inputSymbol, nextStateId := range currentStateTransitions {
			if inputSymbol == string(character) {
				currentState = nextStateId
				break
			}
		}
	}

	answer := false

	for _, acceptState := range d.AcceptStates {
		if acceptState == currentState {
			answer = true
			break
		}
	}

	return answer
}

func (d *DFA) ConvertToJSON() *model.DrawFAResponse {
	edges := make([]*model.Edge, 0)

	for fromStateId, transitions := range d.TransitionMap {
		for inputSymbol, toStateId := range transitions {
			edge := &model.Edge{
				From:  fromStateId,
				To:    toStateId,
				Label: inputSymbol,
			}
			edges = append(edges, edge)
		}
	}

	return &model.DrawFAResponse{
		Edges:        edges,
		Nodes:        d.States,
		StartState:   d.StartState,
		AcceptStates: d.AcceptStates,
		DeadState:    d.DeadStateId,
	}
}

func slice2str(arr []int) string {
	sort.Slice(arr, func(i, j int) bool {
		return arr[i] < arr[j]
	})

	var result string

	for idx, v := range arr {
		if idx != 0 {
			result += ","
		}
		s := strconv.Itoa(v)
		result += s
	}

	return result
}

func getKeys(m map[int]bool) []int {
	keys := make([]int, 0)
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func (d *DFA) Suffix() *NFA {
	nfa := NewNFA()
	nfa.InputSymbols = d.InputSymbols
	nfa.InputSymbols = append(nfa.InputSymbols, epsilonSymbol)
	nfa.States = d.States
	nfa.InputSymbols = d.InputSymbols
	nfa.AcceptStates = d.AcceptStates
	nfa.IDCount = d.dfaStatesCount

	newBeginStateId := nfa.AddState()
	nfa.SetStartState(newBeginStateId)

	transitionMap := make(map[int]map[string][]int)
	for fromStateId, transitions := range d.TransitionMap {
		for inputSymbol, toState := range transitions {
			if _, ok := transitionMap[fromStateId]; !ok {
				transitionMap[fromStateId] = make(map[string][]int)
			}

			if _, ok := transitionMap[fromStateId][inputSymbol]; !ok {
				transitionMap[fromStateId][inputSymbol] = make([]int, 0)
			}

			transitionMap[fromStateId][inputSymbol] = append(transitionMap[fromStateId][inputSymbol], toState)
		}
	}

	transitionMap[newBeginStateId] = make(map[string][]int)
	transitionMap[newBeginStateId][epsilonSymbol] = make([]int, 0)

	transitionMap[newBeginStateId][epsilonSymbol] = append(transitionMap[newBeginStateId][epsilonSymbol], nfa.States...)
	nfa.TransitionMap = transitionMap

	return nfa
}
