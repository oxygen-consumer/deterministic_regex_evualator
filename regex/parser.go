package regex

import "fmt"

func Parse(infix []Token) ([]Token, error) {
	postfix, err := toPostfix(insertConcat(infix))
	if err != nil {
		return nil, err
	}

	return postfix, nil
}

func toPostfix(infix []Token) ([]Token, error) {
	output := make([]Token, 0, len(infix))
	stack := make([]Token, 0, len(infix))

	for _, tok := range infix {
		switch tok.Type {
		case CHAR:
			output = append(output, tok)
		case LPAREN:
			stack = append(stack, tok)
		case RPAREN:
			foundLeftParen := false

			for len(stack) > 0 {
				top := stack[len(stack)-1]
				stack = stack[:len(stack)-1]

				if top.Type == LPAREN {
					foundLeftParen = true
					break
				}

				output = append(output, top)
			}

			if !foundLeftParen {
				return nil, fmt.Errorf("missmatched parantheses")
			}

		case KLEENE, PLUS, QUESTION, OR, CONCAT:
			for len(stack) > 0 {
				top := stack[len(stack)-1]

				if (top.Type != LPAREN) && (top.precedence() >= tok.precedence()) {
					output = append(output, top)
					stack = stack[:len(stack)-1]
				} else {
					break
				}
			}

			stack = append(stack, tok)
		default:
			return nil, fmt.Errorf("unknown token: %v", tok)
		}
	}

	for len(stack) > 0 {
		top := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		if top.Type == LPAREN || top.Type == RPAREN {
			return nil, fmt.Errorf("missmatched parantheses")
		}

		output = append(output, top)
	}

	return output, nil
}

func insertConcat(tokens []Token) []Token {
	result := make([]Token, 0, len(tokens)*2)

	for i := range tokens {
		result = append(result, tokens[i])

		if i+1 < len(tokens) {
			current := tokens[i]
			next := tokens[i+1]

			if (current.Type == CHAR || current.Type == RPAREN || current.Type == KLEENE ||
				current.Type == PLUS || current.Type == QUESTION) &&
				(next.Type == CHAR || next.Type == LPAREN) {
				result = append(result, Token{Type: CONCAT, Value: '.'})
			}
		}
	}

	return result
}
