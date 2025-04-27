package nfa

import (
	"fmt"
	"strings"
)

type State struct {
	ID                 int
	Transitions        map[rune][]*State
	EpsilonTransitions []*State
}

type NFA struct {
	Start  *State
	Accept *State // NOTE: only one accepting state needed for Thompson
}

var nextStateID int = 0

func newState() *State {
	state := &State{
		ID:          nextStateID,
		Transitions: make(map[rune][]*State),
	}

	nextStateID++
	return state
}

func (nfa *NFA) ToDot() string {
	// NOTE: written with chat gipiti
	var b strings.Builder

	b.WriteString("digraph NFA {\n")
	b.WriteString("\trankdir=LR;\n")
	b.WriteString("\tsize=\"10,7\";\n")
	b.WriteString("\tdpi=300;\n")

	// Step 1: Collect all states and accepting states
	allStates := make(map[int]*State)
	acceptingStates := make(map[int]bool)

	var collect func(s *State)
	visited := make(map[int]bool)

	collect = func(s *State) {
		if visited[s.ID] {
			return
		}
		visited[s.ID] = true
		allStates[s.ID] = s
		if s.ID == nfa.Accept.ID {
			acceptingStates[s.ID] = true
		}
		for _, next := range s.EpsilonTransitions {
			collect(next)
		}
		for _, nexts := range s.Transitions {
			for _, next := range nexts {
				collect(next)
			}
		}
	}

	collect(nfa.Start)

	// Step 2: Print accepting states
	b.WriteString("\tnode [shape=doublecircle];\n")
	for id := range acceptingStates {
		b.WriteString(fmt.Sprintf("\tq%d;\n", id))
	}

	// Step 3: Back to normal states
	b.WriteString("\tnode [shape=none]; start;\n")
	b.WriteString("\tnode [shape=circle];\n")

	// Step 4: Start arrow
	b.WriteString(fmt.Sprintf("\tstart -> q%d;\n", nfa.Start.ID))

	// Step 5: Transitions
	for _, state := range allStates {
		for symbol, targets := range state.Transitions {
			for _, target := range targets {
				b.WriteString(fmt.Sprintf("\tq%d -> q%d [label=\"%c\"];\n", state.ID, target.ID, symbol))
			}
		}
		for _, target := range state.EpsilonTransitions {
			b.WriteString(fmt.Sprintf("\tq%d -> q%d [label=\"ε\"];\n", state.ID, target.ID))
		}
	}

	b.WriteString("}\n")
	return b.String()
}
