package core

import (
	"fmt"
	"re2fa/model"
)

type DFAState struct {
	Id int `json:"id"`
	// 确定性的，同一个state，同一个input symbol，通过delta function只能转移到一个确定的state
	Transitions map[string]*DFAState `json:"transitions"`
}

type Transition struct {
	from        int
	to          int
	inputSymbol string
}

type DFA struct {
	Start        *DFAState   `json:"start"`
	AcceptStates []*DFAState `json:"acceptStates"`
}

type NoEpsilonFA struct {
	StartStates map[int]*State `json:"startStates"`
	EndStates   map[int]*State `json:"endStates"`
	States      map[int]*State `json:"states"`
}

func NoEpsilonFAConstructor(nfa *Nfa) *NoEpsilonFA {
	nefa := &NoEpsilonFA{
		StartStates: make(map[int]*State),
		EndStates:   make(map[int]*State),
		States:      make(map[int]*State),
	}

	nefa.StartStates[nfa.StartState.Id] = nfa.StartState
	nefa.EndStates[nfa.EndState.Id] = nfa.EndState

	var dfs func(state *State)

	dfs = func(state *State) {
		if state == nil {
			return
		}

		stateId := state.Id

		if _, ok := nefa.States[stateId]; !ok {
			nefa.States[stateId] = state
		}

		for _, nextStates := range state.Transitions {
			for _, nextState := range nextStates {
				dfs(nextState)
			}
		}
	}

	dfs(nfa.StartState)

	nefa.removeEpsilonTransition()

	return nefa
}

func (n *NoEpsilonFA) findAllEpsilonStartState() map[int]bool {
	result := make(map[int]bool)

	var dfsFindAllEpsilonStartState func(state *State)

	dfsFindAllEpsilonStartState = func(s *State) {
		if s == nil {
			return
		}

		for inputSymbol, nextStates := range s.Transitions {
			if inputSymbol == "-1" {
				result[s.Id] = true
			}

			for _, nextState := range nextStates {
				dfsFindAllEpsilonStartState(nextState)
			}
		}
	}

	for _, state := range n.States {
		dfsFindAllEpsilonStartState(state)
	}

	return result
}

func (n *NoEpsilonFA) removeEpsilonTransition() {
	epsilonStartStates := n.findAllEpsilonStartState()

	for len(epsilonStartStates) > 0 {
		for stateId := range epsilonStartStates {
			epsilonStartState := n.States[stateId]

			isStartState := false

			if _, ok := n.StartStates[epsilonStartState.Id]; ok {
				isStartState = true
			}

			for _, epsilonNextState := range epsilonStartState.Transitions["-1"] {
				if isStartState {
					n.StartStates[epsilonNextState.Id] = epsilonNextState
				}

				if _, ok := n.EndStates[epsilonNextState.Id]; ok {
					n.EndStates[epsilonStartState.Id] = epsilonStartState
				}

				for inputSymbol, nextStates := range epsilonNextState.Transitions {
					for _, nextState := range nextStates {
						if epsilonStartState.Transitions[inputSymbol] == nil {
							epsilonStartState.Transitions[inputSymbol] = make([]*State, 0)
						}

						// TODO: use hashmap
						alreadyExist := false

						for _, st := range epsilonStartState.Transitions[inputSymbol] {
							if st.Id == nextState.Id {
								alreadyExist = true
								break
							}
						}

						if !alreadyExist {
							epsilonStartState.Transitions[inputSymbol] = append(epsilonStartState.Transitions[inputSymbol], nextState)
						}
					}
				}
			}

			delete(epsilonStartState.Transitions, "-1")
		}

		epsilonStartStates = n.findAllEpsilonStartState()
	}
}

func (n *NoEpsilonFA) ConvertToDrawRes() *model.DrawDFAResponse {
	edges := make([]*model.Edge, 0)
	nodes := make([]int, 0)
	startStates := make([]int, 0)
	accetpStates := make([]int, 0)

	for id := range n.States {
		nodes = append(nodes, id)
	}

	for id := range n.StartStates {
		startStates = append(startStates, id)
	}

	for id := range n.EndStates {
		accetpStates = append(accetpStates, id)
	}

	var dfsFindEdge func(state *State)
	visitedEdge := make(map[string]bool)

	dfsFindEdge = func(currentState *State) {
		if currentState == nil {
			return
		}

		from := currentState.Id

		for inputSymbol, nextStates := range currentState.Transitions {
			for _, nextState := range nextStates {
				to := nextState.Id
				fti := fmt.Sprintf("%d-%d-%s", from, to, inputSymbol)

				if _, ok := visitedEdge[fti]; ok {
					continue
				}
				visitedEdge[fti] = true

				edge := model.Edge{
					From:  from,
					To:    to,
					Label: inputSymbol,
				}

				edges = append(edges, &edge)

				dfsFindEdge(nextState)
			}
		}
	}

	for _, sta := range n.States {
		dfsFindEdge(sta)
	}

	jsonResult := model.DrawDFAResponse{
		Nodes:        nodes,
		Edges:        edges,
		StartStates:  startStates,
		AcceptStates: accetpStates,
	}

	return &jsonResult
}

// func NewDFAFromNFA(n *Nfa) interface{} {
// 	noEpsilonFA := NoEpsilonFAConstructor(n)
// 	noEpsilonFA.removeEpsilonTransition()
// 	return noEpsilonFA
// }
