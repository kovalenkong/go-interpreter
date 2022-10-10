package go_interpreter

type TokenType uint

const (
	EOF    TokenType = iota
	LPAREN           // (
	RPAREN           // )

	IDENT
	DELIMITER // ;

	NUMBER // i.e. 12.34 or 23
	STRING // i.e. hello note! without quotes

	ADD // +
	SUB // -
	MUL // *
	DIV // /
	EXP // ^

	EQ  // =
	LT  // <
	GT  // >
	LTE // <=
	GTE // >=
)
