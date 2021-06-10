package parser

import (
	"interpreter/ast"
	"interpreter/lexer"
	"interpreter/token"

)

type Parser struct {
	l *lexer.Lexer
	curToken token.Token
	peekToken token.Token
}

func New(l *lexer.Lexer) *Parser {
	p:= &Parser{l: l}

	// read two tokens so current and next are both set
	p.NextToken()
	p.NextToken()
	return p
}

func(p *Parser) NextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func(p *Parser) ParseProgram() *ast.Program {return nil}

