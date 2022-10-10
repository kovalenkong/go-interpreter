package go_interpreter

// Node type implemented by all Node structures.
type Node interface {
	Pos() uint
}

// A Literal node represents a literal of basic type.
type Literal struct {
	pos   uint
	Kind  TokenType
	Value string
}

// An Ident node represents an identifier.
type Ident struct {
	pos  uint
	Name string
}

// A BinaryExpr node represents a binary expression.
type BinaryExpr struct {
	pos   uint
	Left  Node
	Right Node
	Op    TokenType
}

// A UnaryExpr node represents a unary expression.
type UnaryExpr struct {
	pos  uint
	Left Node
	Op   TokenType
}

// A Function node represents a function literal.
type Function struct {
	pos  uint
	Name string
	Args []Node
}

// A Comparison node represents a comparison expression.
type Comparison struct {
	pos   uint
	Left  Node
	Right Node
	Op    TokenType
}

func (l *Literal) Pos() uint {
	return l.pos
}
func (c *Comparison) Pos() uint {
	return c.pos
}

func (i *Ident) Pos() uint {
	return i.pos
}

func (b *BinaryExpr) Pos() uint {
	return b.pos
}

func (u *UnaryExpr) Pos() uint {
	return u.pos
}

func (f *Function) Pos() uint {
	return f.pos
}
