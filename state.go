package regex

// State this structure represents a finite automata state
type state struct {
	accepting   bool
	transitions map[string][]*state
}

// NewState creates an empty finite automata state with its target transitions
func newState(accepting bool, transitions map[string][]*state) state {
	return state{
		accepting:   accepting,
		transitions: transitions,
	}
}

// EmptyState creates an empty finite automata state meaning without any targets
func emptyState(accepting bool) state {
	return state{
		accepting:   accepting,
		transitions: map[string][]*state{},
	}
}

// AddTransitionForSymbol adds a new symbol's state to the list of transitions for that symbol
func (thisState state) addTransitionForSymbol(symbol string, newState *state) {
	if transitionsForSymbol, ok := thisState.transitions[symbol]; ok {
		thisState.transitions[symbol] = append(transitionsForSymbol, newState)
	} else {
		thisState.transitions[symbol] = append([]*state{}, newState)
	}
}

// getTransitionsForSymbol returns the list of transitions for a given symbol
func (thisState state) getTransitionsForSymbol(symbol string) ([]*state, bool) {
	transitionsForSymbol, ok := thisState.transitions[symbol]
	return transitionsForSymbol, ok
}

func (thisState state) setAccepting(accepting bool) {
	thisState.accepting = accepting
}

func (thisState state) test(symbol string) bool {
	return false
}
