package main

type transition struct {
	id    int
	from  int  // from state
	to    int  // to state
	input byte // input symbol
}

func new_transition(id int, fromStateId int, toStateId int, inputC byte) *transition {
	return &transition{
		id:    id,
		from:  fromStateId,
		to:    toStateId,
		input: inputC,
	}
}

type state struct {
	id          int
	transitions []*transition
}

func new_state(id int, transitions []*transition) *state {
	return &state{
		id:          id,
		transitions: transitions,
	}
}

// move
type move struct {
	start *state
	end   *state
}

func new_move(s, e *state) *move {
	return &move{
		start: s,
		end:   e,
	}
}

type nfa struct {
	start          int
	end            int
	state_map      map[int]*state
	transition_map map[int]*transition
}

func new_nfa() *nfa {
	return &nfa{
		start:          -1,
		end:            -1,
		state_map:      make(map[int]*state),
		transition_map: make(map[int]*transition),
	}
}

func (n *nfa) createFromRegexp(str string) {
	if str == "" {
		panic("empty regexp")
	}

	postfix := re2postfix(str)
	n.createFromPostfix(postfix)
}

func (n *nfa) createFromPostfix(postfix string) {
	if postfix == "" {
		panic("empty postfix")
	}

	uuid_state := 0
	uuid_transition := 0

	var newState = func() *state {
		empty_transitions := make([]*transition, 0)
		uuid_state += 1
		state := new_state(uuid_state, empty_transitions)
		n.state_map[uuid_state] = state
		return state
	}

	var newTransition = func(fromState *state, toState *state, inputC byte) *transition {
		uuid_transition += 1
		transition := new_transition(uuid_transition, fromState.id, toState.id, inputC)
		fromState.transitions = append(fromState.transitions, transition)
		n.transition_map[uuid_transition] = transition
		return transition
	}

	moves := new_movestack()

	for i := 0; i < len(postfix); i++ {
		ch := postfix[i]

		if ch == '|' {

		} else if ch == '.' {
			leftVal := moves.out()
			rightVal := moves.out()

			for _, transition := range rightVal.start.transitions {
				transition_id := transition.id
				transition := n.transition_map[transition_id]
				transition.from = leftVal.end.id
				leftVal.end.transitions = append(leftVal.end.transitions, transition)
			}
		} else if ch == '*' {

		} else if ch == '?' {

		} else if ch == '+' {

		} else {
			// normal character
			start := newState()
			end := newState()
			newTransition(start, end, ch)
			move := new_move(start, end)
			moves.in(move)
		}
	}
}
