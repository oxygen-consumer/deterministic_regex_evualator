package dfa

func RunDFA(dfa *DFA, testString string) bool {
	currentState := dfa.Start

	for _, symbol := range testString {
		if nextState, exists := currentState.Transitions[symbol]; exists {
			currentState = nextState
		} else {
			return false
		}
	}

	return currentState.IsAccepting
}
