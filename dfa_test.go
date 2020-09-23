package regex

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetAlphabet(t *testing.T) {
	regex := or(
		char("a"),
		char("b"),
		char("c"),
		char("d"),
	)
	expected := map[string]bool{
		"a": true,
		"b": true,
		"c": true,
		"d": true,
	}
	someDfa := dfa{nfa: regex.completeNfa()}
	assert.Equal(t, expected, someDfa.getAlphabet(), "the alphabet calculated is incorrect")
}

func TestJoinManyInts(t *testing.T) {
	joinded := joinInts("|", 1, 2, 3)
	assert.Equal(t, "1|2|3", joinded)
}
func TestJoinOneInt(t *testing.T) {
	joinded := joinInts("|", 1)
	assert.Equal(t, "1", joinded)
}

func TestJoinZeroInt(t *testing.T) {
	assert.Equal(t, "", joinInts("|"))
}

func TestGetTransitionTable(t *testing.T) {
		regex := plus(or(char("a"), char("b"), concat(char("a"), char("b"))))
	dfaValue := dfa{nfa: regex.completeNfa()}
	fmt.Println(clearNfaTransitionTable(dfaValue.getNfaTransitionTable()))
	fmt.Println(dfaValue.getTransitionTable())
	fmt.Println(dfaValue.getAccepting())

}
