package regex

const EPSILON = "ε"

type indexValue struct {
	id        string
	accepting bool
	number    int
}

type augmentedNfa struct {
	nfa             nfa
	transitionTable map[int]map[string][]indexValue
	alphabet        map[string]bool
	accepting       map[int]bool
	starting        int
}

type nfa struct {
	inState  *state
	outState *state
}

func (thisNfa nfa) test(symbol string) bool {
	return thisNfa.inState.test(symbol, map[string]bool{})
}

func mergeMap(someMap map[string]int, key string, value int) {
	if _, ok := someMap[key]; !ok {
		someMap[key] = value
	}
}

func mergeMapAndGet(someMap map[string][]indexValue, key string, value []indexValue) []indexValue {
	if _, ok := someMap[key]; !ok {
		someMap[key] = value
	}
	return someMap[key]
}

func (thisNfa nfa) completeNfa() augmentedNfa {
	augmentedNfa := augmentedNfa{
		nfa: thisNfa,
		starting:  0,
		alphabet:  map[string]bool{},
		accepting: map[int]bool{},
	}
	result := map[int]map[string][]indexValue{}
	visited := map[string]int{}
	stack := []*state{thisNfa.inState}
	for len(stack) > 0 {
		current := stack[0]
		var dests = map[string][]indexValue{}
		if _, ok := result[visited[current.id]]; !ok {
			mergeMap(visited, current.id, len(visited))
			for symbol, value := range current.transitions {
				for _, transition := range value {
					mergeMap(visited, transition.id, len(visited))
					stack = append(stack, transition)
					dests[symbol] = append(mergeMapAndGet(dests, symbol, []indexValue{}), indexValue{
						id:        transition.id,
						accepting: transition.accepting,
						number:    visited[transition.id],
					})
					if symbol != EPSILON {
						augmentedNfa.alphabet[symbol] = true
					}
				}
			}
			dests[EPSILON] = append(mergeMapAndGet(dests, EPSILON, []indexValue{}), indexValue{
				id:        current.id,
				accepting: current.accepting,
				number:    visited[current.id],
			})
			result[visited[current.id]] = dests
			if current.accepting {
				augmentedNfa.accepting[visited[current.id]] = true
			}
		}
		stack = stack[1:]
	}
	augmentedNfa.transitionTable = result
	return augmentedNfa
}


func clearNfaTransitionTable(transitionTable map[int]map[string][]indexValue) map[int]map[string][]int {
	result := map[int]map[string][]int{}
	for key, value := range transitionTable {
		result[key] = map[string][]int{}
		for symbol, transitions := range value {
			ints := []int{}
			for _, transition := range transitions {
				ints = append(ints, transition.number)
			}
			result[key][symbol] = ints
		}
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
