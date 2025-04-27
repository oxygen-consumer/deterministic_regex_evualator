package regex

func Tokenize(regexString string) []Token {
	tokens := make([]Token, 0, len(regexString))

	for _, tok := range regexString {
		switch tok {
		case '*':
			tokens = append(tokens, Token{Type: KLEENE, Value: '*'})
		case '+':
			tokens = append(tokens, Token{Type: PLUS, Value: '+'})
		case '?':
			tokens = append(tokens, Token{Type: QUESTION, Value: '?'})
		case '|':
			tokens = append(tokens, Token{Type: OR, Value: '|'})
		case '(':
			tokens = append(tokens, Token{Type: LPAREN, Value: '('})
		case ')':
			tokens = append(tokens, Token{Type: RPAREN, Value: ')'})
		default:
			tokens = append(tokens, Token{Type: CHAR, Value: tok})
		}
	}

	return tokens
}
