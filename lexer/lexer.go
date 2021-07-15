package lexer

import (
	"interpreter/token"
)

type Lexer struct {
	input        string
	position     int  // current position in input (points to current char)
	readPosition int  // current reading position in input (after current char)
	ch           byte // character under examination
}

func New(input string) *Lexer {

	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++

}

func (l *Lexer) NextToken() token.Token {
	var t token.Token

	l.skipWhiteSpaces()

	switch l.ch {
	case '=':
		if l.peakChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			t = token.Token{Type: token.EQ, Literal: literal}
		} else {
			t = newToken(token.ASSIGN, l.ch)
		}

	case '!':
		if l.peakChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			t = token.Token{Type: token.NOT_EQ, Literal: literal}
		} else {
			t = newToken(token.BANG, l.ch)
		}
	case '+':
		t = newToken(token.PLUS, l.ch)
	case '-':
		t = newToken(token.MINUS, l.ch)
	case '*':
		t = newToken(token.ASTERISK, l.ch)
	case '/':
		t = newToken(token.SLASH, l.ch)
	case '(':
		t = newToken(token.LPAREN, l.ch)
	case ')':
		t = newToken(token.RPAREN, l.ch)
	case '{':
		t = newToken(token.LBRACE, l.ch)
	case '}':
		t = newToken(token.RBRACE, l.ch)
	case '<':
		t = newToken(token.LT, l.ch)
	case '>':
		t = newToken(token.GT, l.ch)
	case ';':
		t = newToken(token.SEMICOLON, l.ch)
	case ',':
		t = newToken(token.COMMA, l.ch)
	case '"':
		t.Type = token.STRING
		t.Literal = l.readString()
	case '[':
		t = newToken(token.LBRACKET, l.ch)
	case ']':
		t = newToken(token.RBRACKET, l.ch)
	case ':':
		t = newToken(token.COLON,l.ch)
	case 0:
		t.Literal = ""
		t.Type = token.EOF
	default:
		if isLetter(l.ch) {

			t.Literal = l.readIdentifier()
			t.Type = token.LookupIdentifier(t.Literal)
			return t

		} else if isDigit(l.ch) {
			t.Type = token.INT
			t.Literal = l.readNumber()
			return t
		} else {
			t = newToken(token.ILLEGAL, l.ch)
		}

	}
	l.readChar()
	return t
}

/*


 */

//newToken is a helper function to initialize the token.Token
func newToken(tokenType token.TokenType, ch byte) token.Token {
	t := token.Token{
		Type:    tokenType,
		Literal: string(ch),
	}
	return t
}

//readIdentifier reads the input and advances the position until a non letter is encountered
//it returns the identifier with input[position]
func (l *Lexer) readIdentifier() string {
	position := l.position

	for isLetter(l.ch) {
		l.readChar()
	}

	return l.input[position:l.position]
}

//isLetter returns true if the passed byte is a permitted letter, we can include special characters *,?, etc just like "_"
func isLetter(ch byte) bool {
	return ch >= 'a' && ch <= 'z' || ch >= 'A' && ch <= 'Z' || ch == '_'
}

func (l *Lexer) readNumber() string {
	position := l.position

	for isDigit(l.ch) {
		l.readChar()
	}

	return l.input[position:l.position]
}

func isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

func (l *Lexer) skipWhiteSpaces() {

	for l.ch == ' ' || l.ch == '\n' || l.ch == '\t' || l.ch == '\r' {
		l.readChar()
	}
}

func (l *Lexer) peakChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition]
	}

}

func (l *Lexer) readString() string {

	position := l.position + 1
	for {
		l.readChar()
		if l.ch == '"' || l.ch == 0 {
			break
		}
	}
	return l.input[position:l.position]
}
