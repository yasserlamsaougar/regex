package regex

import (
	"github.com/google/uuid"
)

// State this structure represents a finite automata state
type state struct {
	accepting   bool
	transitions map[string][]*state
	id          string
}

// NewState creates an empty finite automata state with its target transitions
func newState(accepting bool, transitions map[string][]*state) state {
	id, _ := uuid.NewRandom()
	return state{
		accepting:   accepting,
		transitions: transitions,
		id:          id.String(),
	}
}

// EmptyState creates an empty finite automata state meaning without any targets
func emptyState(accepting bool) state {
	id, _ := uuid.NewRandom()
	return state{
		accepting:   accepting,
		transitions: map[string][]*state{},
		id:          id.String(),
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

func (thisState state) test(symbol string, visited map[string]bool) bool {
	formatted := symbol + "_" + thisState.id
	if _, ok := visited[formatted]; ok {
		return false
	}
	if len(symbol) == 0 {
		if thisState.accepting {
			return true
		} else {
			if epsilonTransitions, ok := thisState.getTransitionsForSymbol(EPSILON); ok {
				for _, transition := range epsilonTransitions {
					visited[formatted] = true
					if transition.test(symbol, visited) {
						return true
					}
				}
			}
			return false
		}
	}
	char := symbol[0:1]
	rest := symbol[1:]
	if transitions, ok := thisState.getTransitionsForSymbol(char); ok {
		for _, transition := range transitions {
			if transition.test(rest, visited) {
				return true
			}
		}
	}
	if transitions, ok := thisState.getTransitionsForSymbol(EPSILON); ok {
		for _, transition := range transitions {
			visited[formatted] = true
			if transition.test(symbol, visited) {
				return true
			}
		}
	}
	return false
}
