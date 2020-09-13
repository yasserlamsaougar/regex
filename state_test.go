package regex

import (
	"testing"
)

func TestEmptyState(t *testing.T) {
	newState := emptyState(true)
	if !newState.accepting {
		t.Error("the new state should have been created with an accepting field with value true")
	}
	if len(newState.transitions) != 0 {
		t.Errorf("The new state should have an empty transitions field not a size %d", len(newState.transitions))
	}
}

func TestNewState(t *testing.T) {
	newState := newState(false, map[string][]state{
		"a": []state{},
	})
	if newState.accepting {
		t.Error("the new state should have been created with an accepting field with value false")
	}
	if len(newState.transitions) != 1 {
		t.Errorf("The new state should have a transitions field with length 1 not a length %d", len(newState.transitions))
	}
	if _, ok := newState.transitions["a"]; !ok {
		t.Errorf("could not find the key a in the list of transitions %v", newState.transitions)
	}
}

func TestAddTransitionForSymbol(t *testing.T) {
	newState := emptyState(true)
	newState.addTransitionForSymbol("a", emptyState(false))
	if _, ok := newState.transitions["a"]; !ok {
		t.Errorf("could not find the key a in the list of transitions %v", newState.transitions)
	}
}

func TestGetTransitionsForSymbol(t *testing.T) {
	newState := newState(false, map[string][]state{
		"a": []state{},
	})
	if _, ok := newState.getTransitionsForSymbol("a"); !ok {
		t.Errorf("could not find the key a in the list of transitions %v", newState.transitions)
	}
	if _, ok := newState.getTransitionsForSymbol("b"); ok {
		t.Errorf("should not find the key b in the list of transitions %v", newState.transitions)
	}
}
