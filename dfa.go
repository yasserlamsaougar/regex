package regex


type dfa struct {
	someNfa *nfa
}

func (thisDfa dfa) getAlphabet() map[string]bool {

	return map[string]bool{}
}
