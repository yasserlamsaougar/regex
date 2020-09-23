package regex

import (
	"fmt"
	"strconv"
	"strings"
)

const SEP = "|"

type dfa struct {
	nfa augmentedNfa
}

func (thisDfa dfa) getAlphabet() map[string]bool {
	return thisDfa.nfa.alphabet
}

func (thisDfa dfa) getAccepting() map[int]bool {
	return thisDfa.nfa.accepting
}

func (thisDfa dfa) getNfaTransitionTable() map[int]map[string][]indexValue {
	return thisDfa.nfa.transitionTable
}

func (thisDfa dfa) getStartingState() int {
	return thisDfa.nfa.starting
}

func joinInts(sep string, ints ...int) string {
	stringBuilder := strings.Builder{}
	intsLength := len(ints)
	for i := 0; i < intsLength-1; i++ {
		stringBuilder.WriteString(strconv.Itoa(ints[i]))
		stringBuilder.WriteString(sep)
	}
	if intsLength > 0 {
		stringBuilder.WriteString(strconv.Itoa(ints[intsLength-1]))
	}
	return stringBuilder.String()
}

func getAndFilterStateNumbers(indexValues []indexValue, predicate func(value indexValue) bool) []int {
	results := []int{}
	for _, iv := range indexValues {
		if predicate(iv) {
			results = append(results, iv.number)
		}
	}
	return results
}

func getStateNumbers(indexValues []indexValue) []int {
	return getAndFilterStateNumbers(indexValues, func(_ indexValue) bool {
		return true
	})
}

func intsToWorkStates(symbol string, hasSymbol bool, ints ...int) []workState {
	results := []workState{}
	for _, iv := range ints {
		results = append(results, workState{
			number:    iv,
			symbol:    symbol,
			hasSymbol: hasSymbol,
		})
	}
	return results
}

type workState struct {
	number    int
	hasSymbol bool
	symbol    string
}

type joinedState struct {
	joinedStateInts  []int
	joinedState      string
	joinedStateIndex map[string]map[int]bool
}

func (thisDfa dfa) getTransitionTable() map[string]map[string]string {
	finalResult := map[string]map[string]string{}
	transitionTable := thisDfa.getNfaTransitionTable()
	startingState := thisDfa.getStartingState()
	init := transitionTable[startingState][EPSILON]
	stateNumbers := getStateNumbers(init)
	workStack := []workState{}
	stateStack := []joinedState{{
		joinedStateInts:  stateNumbers,
		joinedState:      joinInts(SEP, stateNumbers...),
		joinedStateIndex: map[string]map[int]bool{},
	}}
	visited := map[string]bool{}
	for len(stateStack) > 0 {
		currentJoinedstate := stateStack[0]
		newJoinedState := map[string][]int{}
		for _, currentState := range currentJoinedstate.joinedStateInts[:len(currentJoinedstate.joinedStateInts)] {
			for symbol, transitions := range transitionTable[currentState] {
				if symbol == EPSILON && len(transitions) == 1 {
					continue
				}
				workStack = append(workStack, workState{
					number:    currentState,
					symbol:    symbol,
					hasSymbol: symbol != EPSILON,
				})
			}
		}
		fmt.Println(currentJoinedstate.joinedState)
		for len(workStack) > 0 {
			currentWorkState := workStack[0]
			transitionsBySymbol := transitionTable[currentWorkState.number]
			// if we only have an epsilon state
			if len(transitionsBySymbol) == 1 {
				epsilonTransitions := transitionsBySymbol[EPSILON]
				if currentWorkState.hasSymbol {
					for _, value := range epsilonTransitions {
						if _, ok := currentJoinedstate.joinedStateIndex[currentWorkState.symbol]; !ok {
							currentJoinedstate.joinedStateIndex[currentWorkState.symbol] = map[int]bool{}
						}
						if _, ok := currentJoinedstate.joinedStateIndex[currentWorkState.symbol][value.number]; !ok {
							newJoinedState[currentWorkState.symbol] = append(newJoinedState[currentWorkState.symbol], value.number)
							currentJoinedstate.joinedStateIndex[currentWorkState.symbol][value.number] = true
						}
					}
				} else {
					for _, value := range epsilonTransitions {
						visitedId := fmt.Sprintf("%s:%d", currentJoinedstate.joinedState, value.number)
						if _, ok := visited[visitedId]; !ok && value.number != currentWorkState.number {
							workStack = append(workStack, workState{
								number:    value.number,
								symbol:    "",
								hasSymbol: false,
							})
							visited[visitedId] = true
						}
					}
				}
			} else {
				for symbol, transitions := range transitionsBySymbol {
					newWorkStateNumbers := getStateNumbers(transitions)
					for _, value := range newWorkStateNumbers {
						if value != currentWorkState.number && value != startingState {
							workStack = append(workStack, workState{
								number:    value,
								hasSymbol: true,
								symbol:    symbol,
							})
						}
					}
				}
			}
			workStack = workStack[1:]
		}
		stateStack = stateStack[1:]
		tempNewJoinedState := finalResult[currentJoinedstate.joinedState]
		if tempNewJoinedState == nil {
			finalResult[currentJoinedstate.joinedState] = map[string]string{}
		}
		for symbol, value := range newJoinedState {
			joinedStateToStack := joinInts(SEP, value...)
			if joinedStateFromFinalResult, ok := tempNewJoinedState[symbol]; !ok || joinedStateFromFinalResult != joinedStateToStack {
				finalResult[currentJoinedstate.joinedState][symbol] = joinedStateToStack
				stateStack = append(stateStack, joinedState{
					joinedStateInts:  value,
					joinedState:      joinedStateToStack,
					joinedStateIndex: map[string]map[int]bool{},
				})
			}
		}
	}
	return finalResult
}

func clearDfaTransitionTable(dfaTable map[string]map[string]string) map[string]map[string]string {
	results := map[string]map[string]string{}
	currentId := 1
	transformMap := map[string]string{}
	for joinedState, statesBySymbol := range dfaTable {
		if _, ok := transformMap[joinedState]; !ok {
			transformMap[joinedState] = strconv.Itoa(currentId)
			currentId++
		}
		results[transformMap[joinedState]] = map[string]string{}
		for symbol, transitions := range statesBySymbol {
			currentTransformation := results[transformMap[joinedState]]
			if _, ok := transformMap[transitions]; !ok {
				transformMap[transitions] = strconv.Itoa(currentId)
				currentId++
			}
			currentTransformation[symbol] = transformMap[transitions]
		}
	}
	return results
}
