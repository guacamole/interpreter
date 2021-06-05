package token

type TokenType string

type Token struct {
	Type TokenType
	Literal string
}

const (
	ILLEGAL ="ILLEGAL"
	EOF ="EOF"

	//identifier and literals

	IDENT ="IDENT" //add,a,b,x,foo
	INT ="INT" //1,2,3

	//operators

	ASSIGN = "="
	PLUS = "+"

	//Delimiter

	COMMA = ","
	SEMICOLON = ";"
	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"

	//keywords

	FUNCTION = "FUNCTION"
	LET = "LET"

	)