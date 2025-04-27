package regex

import "fmt"

func ToPostfix(infix []Token) ([]Token, error) {
	var output []Token
	var stack []Token

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

				if (top.Type != LPAREN) && (precedence(top) >= precedence(tok)) {
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

func precedence(tok Token) int {
	switch tok.Type {
	case KLEENE, PLUS, QUESTION:
		return 3
	case CONCAT:
		return 2
	case OR:
		return 1
	default:
		return 0
	}
}
