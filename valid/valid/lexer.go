package valid

const (
	scanOther = iota
	scanSymbol
	scanNumber
)

const (
	tokenUnknown = iota
	tokenSymbol
	tokenNumber
	tokenSep
	tokenEquals
	tokenArgOpen
	tokenArgClose
)

type Token struct {
	Type int
	Id   string
}

type Lexer struct {
}

func (l *Lexer) scanType(c rune) int {
	if l.isSymbol(c) {
		return scanSymbol
	} else if l.isNumber(c) {
		return scanNumber
	}

	return scanOther
}

func (l *Lexer) isSymbol(c rune) bool {
	if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '@' {
		return true
	}

	return false
}

func (l *Lexer) isNumber(c rune) bool {
	if (c >= '0' && c <= '9') || c == '.' {
		return true
	}

	return false
}

func (l *Lexer) Tokenize(code string) []*Token {
	tokens := make([]*Token, 0, 255)
	j := 0

	for i := 0; i < len(code); i++ {
		switch l.scanType(rune(code[i])) {
		case scanSymbol:
			for j = i; j < len(code); j++ {
				if !l.isSymbol(rune(code[j])) {
					break
				}
			}

			tok := &Token{tokenSymbol, code[i:j]}
			tokens = append(tokens, tok)

			/* Update position */
			i = j - 1
		case scanNumber:
			for j = i; j < len(code); j++ {
				if !l.isNumber(rune(code[j])) {
					break
				}
			}

			tok := &Token{tokenNumber, code[i:j]}
			tokens = append(tokens, tok)

			/* Update position */
			i = j - 1
		default:
			if code[i] == ' ' || code[i] == '\n' {
				continue
			}

			tok := &Token{}
			tok.Id = code[i : i+1]
			switch code[i] {
			case ',':
				tok.Type = tokenSep
			case '=':
				tok.Type = tokenEquals
			case '(':
				tok.Type = tokenArgOpen
			case ')':
				tok.Type = tokenArgClose
			default:
				tok.Type = tokenUnknown
			}

			tokens = append(tokens, tok)
		}
	}

	return tokens
}
