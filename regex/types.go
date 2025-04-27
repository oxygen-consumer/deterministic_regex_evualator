package regex

type TokenType int

const (
	CHAR     TokenType = iota
	KLEENE             // *
	PLUS               // +
	QUESTION           // ?
	OR                 // |
	CONCAT             // .
	LPAREN             // (
	RPAREN             // )
)

type Token struct {
	Type  TokenType
	Value rune
}

func (tok Token) precedence() int {
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
