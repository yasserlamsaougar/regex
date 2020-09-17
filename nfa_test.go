package regex

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

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

	for index := range alphabet[0 : len(alphabet)-1] {
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

func TestAcceptorSimpleConcat(t *testing.T) {
	regex := concat(
		char("a"),
		char("b"),
	)
	if !regex.test("ab") {
		t.Error("the regex ab should match ab")
	}
	if regex.test("aa") {
		t.Error("the regex ab should not match aa")
	}
}

func TestAcceptorStar(t *testing.T) {
	regex := star(
		concat(
			char("a"),
			char("a"),
		),
	)
	if !regex.test("aaaaaaaaaa") {
		t.Error("the regex should match aaaaaaaaa")
	}
	if regex.test("aaaaaaaab") {
		t.Error("the regex should not match aaaaaaaab")
	}
	if regex.test("aaaaaaaaaaa") {
		t.Error("the regex should not match aaaaaaaaaaa")
	}
}

func TestAcceptorOpt(t *testing.T) {
	regex := opt(
		char("a"),
	)
	if !regex.test("") {
		t.Error("the regex should match empty string")
	}
	if !regex.test("a") {
		t.Error("the regex should match a")
	}
	if regex.test("aa") {
		t.Error("the regex should not match aa")
	}
}

func TestAcceptorOr(t *testing.T) {
	regex := or(
		char("a"),
		char("b"),
		char("c"),
		plus(
			char("a"),
		),
	)
	if !regex.test("a") {
		t.Error("the regex should match a")
	}
	if !regex.test("b") {
		t.Error("the regex should match b")
	}
	if !regex.test("c") {
		t.Error("the regex should match c")
	}
	if !regex.test("aaaaaaaaaa") {
		t.Error("the regex should match aaaaaaaaaa")
	}
	if regex.test("") {
		t.Error("the regex should not match empty string")
	}
	if regex.test("vv") {
		t.Error("the regex should not match vv")
	}
}

func TestGetTransitionTable(t *testing.T) {
	regex := or(
		char("a"),
		char("b"),
	)
	expected := map[int]map[string][]int{
		0: {
			EPSILON: []int{1, 2, 0},
		},
		1: {
			"a":     []int{3},
			EPSILON: []int{1},
		},
		2: {
			"b":     []int{4},
			EPSILON: []int{2},
		},
		3: {
			EPSILON: []int{3, 5},
		},
		4: {
			EPSILON: []int{4, 5},
		},
		5: {
			EPSILON: []int{5},
		},
	}
	transitionTable := regex.getTransitionTable()

	for stateNumber, transitionsBySymbol := range transitionTable {
		expectedStateTransitions := expected[stateNumber]
		for symbol, transitions := range transitionsBySymbol {
			numbers := []int{}
			expectedTransitions := expectedStateTransitions[symbol]
			for _, value := range transitions {
				numbers = append(numbers, value.number)
			}
			assert.ElementsMatchf(t, numbers, expectedTransitions, "the state number %d should point to %+v instead pointing to %+v", stateNumber, expectedTransitions, numbers)
		}
	}

}
