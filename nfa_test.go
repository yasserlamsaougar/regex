package regex

import "testing"

func TestChar(t *testing.T) {
	myNfa := char("a")
	if _, ok := myNfa.inState.getTransitionsForSymbol("a"); !ok {
		t.Errorf("could not find the transitions for %s in the instate instead found this %v", "a", myNfa.inState.transitions)
	}
	if !myNfa.outState.accepting {
		t.Errorf("outState accepting value should be accepting")
	}
}

func TestEpsilon(t *testing.T) {
	myNfa := epsilon()
	if _, ok := myNfa.inState.getTransitionsForSymbol(EPSILON); !ok {
		t.Errorf("could not find the transitions for %s in the inState instead found this %v", EPSILON, myNfa.inState.transitions)
	}
	if !myNfa.outState.accepting {
		t.Errorf("outState accepting value should be accepting")
	}
}

func TestConcatPair(t *testing.T) {
	first := char("a")
	second := char("b")
	concatPair(first, second)
	if first.outState.accepting {
		t.Error("first nfa's outstate accepting value should be false")
	}
	if transitions, ok := first.outState.getTransitionsForSymbol(EPSILON); !ok && len(transitions) == 1 {
		t.Error("first's outstate should be epsilon")
	}
	if !second.outState.accepting {
		t.Error("second nfa's outstate accepting value should be true")
	}

}

func TestConcat(t *testing.T) {
	alphabet := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z"}
	rest := []nfa{}
	for _, letter := range alphabet {
		rest = append(rest, char(letter))
	}
	result := concat(rest[0], rest[1:]...)

	for index, _ := range alphabet[0 : len(alphabet)-1] {
		if rest[index].outState.accepting {
			t.Error("first nfa's outstate accepting value should be false", index)
		}
		if transitions, ok := rest[index].outState.getTransitionsForSymbol(EPSILON); !ok && len(transitions) == 1 {
			t.Error("first's outstate should be epsilon")
		}
	}
	if !result.outState.accepting {
		t.Error("result's nfa's outstate accepting value should be true")
	}
}

func TestOr(t *testing.T) {
	first := char("a")
	second := char("b")
	third := char("c")
	result := or(first, second, third)
	if result.inState.accepting {
		t.Error("first nfa's outstate accepting value should be false")
	}
	if transitions, ok := result.inState.getTransitionsForSymbol(EPSILON); !ok && len(transitions) == 2 {
		t.Error("result's inState should be epsilon and have a length of 2", transitions)
	}
	if !result.outState.accepting {
		t.Error("result's nfa's inState accepting value should be true")
	}
	if len(result.outState.transitions) != 0 {
		t.Error("The number of transitions from or outstate should be 0")
	}
}

func TestStar(t *testing.T) {
	result := star(char("a"))
	if result.inState.accepting {
		t.Error("first nfa's outstate accepting value should be false")
	}
	if transitions, ok := result.inState.getTransitionsForSymbol(EPSILON); !ok && len(transitions) == 2 {
		t.Error("result's inState should be epsilon and have a length of 2", transitions)
	}
	if !result.outState.accepting {
		t.Error("result's nfa's inState accepting value should be true")
	}
	if len(result.outState.transitions) != 1 {
		t.Error("The number of transitions from or outstate should be 1", result.outState.transitions)
	}
}

func TestSmall(t *testing.T) {
	complex := concat(
		char("x"),
		star(char("z")),
		char("b"),
		or(char("a"), char("b"), star(char("c")), char("d")),
	)
	t.Log(complex)
}
