package nfa

import (
	"deterministic_regex_evaluator/regex"
	"fmt"
)

func BuildNFA(postfix []regex.Token) (*NFA, error) {
	stack := make([]*NFA, 0, len(postfix))

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
	return finalNFA, nil
}
