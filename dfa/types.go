package dfa

import (
	"fmt"
	"strings"
)

type State struct {
	ID          int
	Transitions map[rune]*State
	IsAccepting bool
}

type DFA struct {
	Start  *State
	States []*State
}

var nextStateID = 0

func newState() *State {
	state := &State{
		ID:          nextStateID,
		Transitions: make(map[rune]*State),
	}

	nextStateID++
	return state
}

func (dfa *DFA) ToDot() string {
	// NOTE: written with chat gipiti
	var b strings.Builder

	b.WriteString("digraph DFA {\n")
	b.WriteString("\trankdir=LR;\n")
	b.WriteString("\tsize=\"10,7\";\n")
	b.WriteString("\tdpi=300;\n") // High quality

	// Step 1: Accepting states
	b.WriteString("\tnode [shape=doublecircle];\n")
	for _, state := range dfa.States {
		if state.IsAccepting {
			b.WriteString(fmt.Sprintf("\tq%d;\n", state.ID))
		}
	}

	// Step 2: Start node and normal nodes
	b.WriteString("\tnode [shape=none]; start;\n")
	b.WriteString("\tnode [shape=circle];\n")
	b.WriteString(fmt.Sprintf("\tstart -> q%d;\n", dfa.Start.ID))

	// Step 3: Transitions
	for _, state := range dfa.States {
		for symbol, target := range state.Transitions {
			b.WriteString(fmt.Sprintf("\tq%d -> q%d [label=\"%c\"];\n", state.ID, target.ID, symbol))
		}
	}

	b.WriteString("}\n")
	return b.String()
}
