package regex

const EPSILON = "ε"

type nfa struct {
	inState  *state
	outState *state
}

func (thisNfa nfa) test(symbol string) bool {
	return thisNfa.inState.test(symbol)
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

func star(someNfa nfa) nfa {
	start := emptyState(false)
	end := emptyState(true)
	start.addTransitionForSymbol(EPSILON, someNfa.inState)
	start.addTransitionForSymbol(EPSILON, &end)

	end.addTransitionForSymbol(EPSILON, someNfa.inState)

	someNfa.outState.addTransitionForSymbol(EPSILON, &end)
	someNfa.outState.accepting = false
	return nfa{
		inState:  &start,
		outState: &end,
	}
}

func plus(someNfa nfa) nfa {
	return concat(
		someNfa,
		star(someNfa),
	)
}

func opt(someNfa nfa) nfa {
	return or(
		someNfa,
		epsilon(),
	)
}

func epsilon() nfa {
	return char(EPSILON)
}
