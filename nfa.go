package regex

const EPSILON = "ε"

type indexValue struct {
	id        string
	accepting bool
}

type nfa struct {
	inState  *state
	outState *state
}

func (thisNfa nfa) test(symbol string) bool {
	return thisNfa.inState.test(symbol, map[string]bool{})
}

func (thisNfa nfa) getTransitionTable() map[string]map[string][]indexValue {
	result := map[string]map[string][]indexValue{}
	stack := []*state{thisNfa.inState}
	for len(stack) > 0 {
		current := stack[0]
		var dests = map[string][]indexValue{}
		if _, ok := result[current.id]; !ok {
			for symbol, value := range current.transitions {
				for _, transition := range value {
					stack = append(stack, transition)
					if _, ok := dests[symbol]; !ok {
						dests[symbol] = []indexValue{}
					}
					dests[symbol] = append(dests[symbol], indexValue{
						id:        transition.id,
						accepting: transition.accepting,
					})
				}
				if _, ok := dests[EPSILON]; !ok {
					dests[EPSILON] = []indexValue{}
				}
				dests[EPSILON] = append(dests[EPSILON], indexValue{
					id:        current.id,
					accepting: current.accepting,
				})
			}
			result[current.id] = dests
		}
		stack = stack[1:]
	}
	return result

}

func char(symbol string) nfa {
	inState := emptyState(false)
	outState := emptyState(true)
	inState.addTransitionForSymbol(symbol, &outState)
	return nfa{
		inState:  &inState,
		outState: &outState,
	}
}

func concatPair(first nfa, second nfa) nfa {
	first.outState.accepting = false
	second.outState.accepting = true
	first.outState.addTransitionForSymbol(EPSILON, second.inState)
	return nfa{
		inState:  first.inState,
		outState: second.outState,
	}
}

func concat(first nfa, rest ...nfa) nfa {
	for _, element := range rest {
		first = concatPair(first, element)
	}
	return first
}

func orPair(first nfa, second nfa) nfa {
	start := emptyState(false)
	end := emptyState(true)
	start.addTransitionForSymbol(EPSILON, first.inState)
	start.addTransitionForSymbol(EPSILON, second.inState)
	first.outState.addTransitionForSymbol(EPSILON, &end)
	second.outState.addTransitionForSymbol(EPSILON, &end)
	first.outState.accepting = false
	second.outState.accepting = false
	return nfa{
		inState:  &start,
		outState: &end,
	}
}

func or(first nfa, rest ...nfa) nfa {
	for _, element := range rest {
		first = orPair(first, element)
	}
	return first
}

func class(symbol string, symbols ...string) nfa {
	someNfa := char(symbol)
	for _, element := range symbols {
		someNfa.inState.addTransitionForSymbol(element, someNfa.outState)
	}
	return someNfa
}

func star(someNfa nfa) nfa {
	someNfa.inState.addTransitionForSymbol(EPSILON, someNfa.outState)
	someNfa.outState.addTransitionForSymbol(EPSILON, someNfa.inState)
	return someNfa
}

func plus(someNfa nfa) nfa {
	someNfa.outState.addTransitionForSymbol(EPSILON, someNfa.inState)
	return someNfa
}

func opt(someNfa nfa) nfa {
	someNfa.inState.addTransitionForSymbol(EPSILON, someNfa.outState)
	return someNfa
}

func epsilon() nfa {
	return char(EPSILON)
}
