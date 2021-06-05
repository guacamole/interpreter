package lexer

import "interpreter/token"

type Lexer struct {
	input string
	position int // current position in input (points to current char)
	readPosition int // current reading position in input (after current char)
	ch byte // character under examination
}

func New(input string) *Lexer {

	l := &Lexer{input:input}
	l.ReadChar()
	return l
}

func (l *Lexer) ReadChar() {
	if l.readPosition >= len(l.input){
		l.ch = 0
	}else{
		l.ch = l.input[l.readPosition]
		l.position = l.readPosition
		l.readPosition++
	}
}

func (l *Lexer) NextToken() token.Token{
	var t token.Token

	switch l.ch{
	case '=':
		t= newToken(token.ASSIGN,l.ch)
	case '+':
		t = newToken(token.PLUS,l.ch)
	case '(':
		t =newToken(token.LPAREN,l.ch)
	case ')':
		t = newToken(token.RPAREN,l.ch)
	case'{':
		t = newToken(token.LBRACE,l.ch)
	case '}':
		t = newToken(token.RBRACE,l.ch)
	case ';':
		t = newToken(token.SEMICOLON,l.ch)
	case ',':
		t = newToken(token.COMMA,l.ch)
	case 0:
		t.Literal =""
		t.Type = token.EOF
		
	}
	l.ReadChar()
	return t
	
}

func newToken(tokenType token.TokenType, ch byte) token.Token {
	t := token.Token{
		Type:    tokenType,
		Literal: string(ch),
	}
	return t
}