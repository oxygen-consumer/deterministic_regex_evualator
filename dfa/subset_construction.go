package dfa

import (
	"deterministic_regex_evaluator/nfa"
	"fmt"
	"sort"
	"strings"
)

func BuildDFAFromNFA(nfaa *nfa.NFA) *DFA {
	nfaAcceptingStateId := nfaa.Accept.ID

	startClosure := epsilonClosure([]*nfa.State{nfaa.Start})
	stateMap := make(map[string]*State)
	queue := [][]*nfa.State{startClosure}
	dfaStates := []*State{}

	startDFAState := newState()
	for _, state := range startClosure {
		if state.ID == nfaAcceptingStateId {
			startDFAState.IsAccepting = true
			break
		}
	}

	startSignature := generateSignatureFromStates(startClosure)
	stateMap[startSignature] = startDFAState
	dfaStates = append(dfaStates, startDFAState)

	for len(queue) > 0 {
		currentSet := queue[0]
		queue = queue[1:]

		currentSignature := generateSignatureFromStates(currentSet)
		currentDFAState := stateMap[currentSignature]

		for _, symbol := range getAlphabet(currentSet) {
			nextStates := move(currentSet, symbol)
			closure := epsilonClosure(nextStates)

			if len(closure) == 0 {
				continue
			}

			closureSignature := generateSignatureFromStates(closure)

			nextDFAState, exists := stateMap[closureSignature]
			if !exists {
				nextDFAState = newState()
				for _, state := range closure {
					if state.ID == nfaAcceptingStateId {
						nextDFAState.IsAccepting = true
						break
					}
				}

				stateMap[closureSignature] = nextDFAState
				dfaStates = append(dfaStates, nextDFAState)
				queue = append(queue, closure)
			}

			currentDFAState.Transitions[symbol] = nextDFAState
		}
	}

	return &DFA{
		Start:  startDFAState,
		States: dfaStates,
	}
}

func move(states []*nfa.State, symbol rune) []*nfa.State {
	result := make(map[int]*nfa.State)

	for _, state := range states {
		for _, next := range state.Transitions[symbol] {
			result[next.ID] = next
		}
	}

	resultList := make([]*nfa.State, 0, len(result))
	for _, state := range result {
		resultList = append(resultList, state)
	}

	return resultList
}

func epsilonClosure(states []*nfa.State) []*nfa.State {
	stack := make([]*nfa.State, len(states))
	copy(stack, states)

	closure := make(map[int]*nfa.State)
	for _, state := range states {
		closure[state.ID] = state
	}

	for len(stack) > 0 {
		top := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		for _, next := range top.EpsilonTransitions {
			if _, exists := closure[next.ID]; !exists {
				closure[next.ID] = next
				stack = append(stack, next)
			}
		}
	}

	result := make([]*nfa.State, 0, len(closure))
	for _, s := range closure {
		result = append(result, s)
	}

	return result
}

func getAlphabet(states []*nfa.State) []rune {
	seen := make(map[rune]bool)
	alphabet := make([]rune, 0, len(states))

	for _, state := range states {
		for symbol := range state.Transitions {
			if !seen[symbol] {
				seen[symbol] = true
				alphabet = append(alphabet, symbol)
			}
		}
	}

	return alphabet
}

func generateSignatureFromStates(nfaStates []*nfa.State) string {
	ids := make([]int, 0, len(nfaStates))
	for _, state := range nfaStates {
		ids = append(ids, state.ID)
	}

	sort.Ints(ids)
	return strings.Trim(strings.Join(strings.Fields(fmt.Sprint(ids)), ","), "[]")
}
