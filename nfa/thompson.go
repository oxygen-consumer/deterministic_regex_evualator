package nfa

import (
	"deterministic_regex_evaluator/regex"
	"fmt"
	"strings"
)

type State struct {
	ID                 int
	Transitions        map[rune][]*State
	EpsilonTransitions []*State
	IsAccepting        bool
}

type NFA struct {
	Start  *State
	Accept *State // NOTE: only one accepting state needed for Thompson
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
		if s.IsAccepting {
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

func BuildNFA(postfix []regex.Token) (*NFA, error) {
	var stack []*NFA

	for _, tok := range postfix {
		switch tok.Type {
		case regex.CHAR:
			start := newState()
			accept := newState()
			start.Transitions[tok.Value] = append(start.Transitions[tok.Value], accept)

			nfa := &NFA{Start: start, Accept: accept}
			stack = append(stack, nfa)

		case regex.CONCAT:
			if len(stack) < 2 {
				return nil, fmt.Errorf("invalid '.' usage: not enough operands")
			}

			nfa2 := stack[len(stack)-1]
			nfa1 := stack[len(stack)-2]
			stack = stack[:len(stack)-2]

			nfa1.Accept.EpsilonTransitions = append(nfa1.Accept.EpsilonTransitions, nfa2.Start)
			nfa := &NFA{Start: nfa1.Start, Accept: nfa2.Accept}
			stack = append(stack, nfa)

		case regex.OR:
			if len(stack) < 2 {
				return nil, fmt.Errorf("invalid '|' usage: not enough operands")
			}

			nfa2 := stack[len(stack)-1]
			nfa1 := stack[len(stack)-2]
			stack = stack[:len(stack)-2]

			start := newState()
			accept := newState()

			start.EpsilonTransitions = append(start.EpsilonTransitions, nfa1.Start, nfa2.Start)
			nfa1.Accept.EpsilonTransitions = append(nfa1.Accept.EpsilonTransitions, accept)
			nfa2.Accept.EpsilonTransitions = append(nfa2.Accept.EpsilonTransitions, accept)

			nfa := &NFA{Start: start, Accept: accept}
			stack = append(stack, nfa)

		case regex.KLEENE:
			if len(stack) < 1 {
				return nil, fmt.Errorf("invalid '*' usage: not enough operands")
			}

			nfa1 := stack[len(stack)-1]
			stack = stack[:len(stack)-1]

			start := newState()
			accept := newState()

			start.EpsilonTransitions = append(start.EpsilonTransitions, nfa1.Start, accept)
			nfa1.Accept.EpsilonTransitions = append(nfa1.Accept.EpsilonTransitions, nfa1.Start, accept)

			nfa := &NFA{Start: start, Accept: accept}
			stack = append(stack, nfa)

		case regex.PLUS:
			if len(stack) < 1 {
				return nil, fmt.Errorf("invalid '+' usage: not enough operands")
			}

			nfa1 := stack[len(stack)-1]
			stack = stack[:len(stack)-1]

			start := newState()
			accept := newState()

			start.EpsilonTransitions = append(start.EpsilonTransitions, nfa1.Start)
			nfa1.Accept.EpsilonTransitions = append(nfa1.Accept.EpsilonTransitions, nfa1.Start, accept)

			nfa := &NFA{Start: start, Accept: accept}
			stack = append(stack, nfa)

		case regex.QUESTION:
			if len(stack) < 1 {
				return nil, fmt.Errorf("invalid '?' usage: not enough operands")
			}

			nfa1 := stack[len(stack)-1]
			stack = stack[:len(stack)-1]

			start := newState()
			accept := newState()

			start.EpsilonTransitions = append(start.EpsilonTransitions, nfa1.Start, accept)
			nfa1.Accept.EpsilonTransitions = append(nfa1.Accept.EpsilonTransitions, accept)

			nfa := &NFA{Start: start, Accept: accept}
			stack = append(stack, nfa)

		default:
			return nil, fmt.Errorf("unexpected token: %v", tok)
		}
	}

	if len(stack) != 1 {
		return nil, fmt.Errorf("invalid NFA: leftover elements")
	}

	finalNFA := stack[0]
	finalNFA.Accept.IsAccepting = true
	return finalNFA, nil
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
