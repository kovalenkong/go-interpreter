package go_interpreter

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

type Func func(args ...any) (any, error)
type Interpreter struct {
	variables map[string]any
	functions map[string]Func
	node      Node
}

func NewInterpreter(variables map[string]any, functions map[string]Func) *Interpreter {
	return &Interpreter{
		variables: variables,
		functions: functions,
	}
}

func (e *Interpreter) SetVar(name string, value any) {
	e.variables[name] = value
}

func (e *Interpreter) ClearVars() {
	e.variables = map[string]any{}
}

func (e *Interpreter) SetFunction(name string, function Func) {
	e.functions[name] = function
}

// Execute method run node and returns result of Interpreter input function.
func (e *Interpreter) Execute(formula string) (any, error) {
	lexer := NewLexer()
	tokens, err := lexer.Lex(strings.NewReader(formula))
	if err != nil {
		return nil, err
	}
	parser := NewParser()
	node, err := parser.Parse(tokens)
	if err != nil {
		return nil, err
	}
	e.node = node

	return e.execute(e.node)
}

func (e *Interpreter) execute(node Node) (any, error) {
	switch n := node.(type) {
	case *BinaryExpr:
		return e.evalBinaryExpr(n)
	case *Literal:
		return e.evalLiteral(n)
	case *Ident:
		return e.evalIdent(n)
	case *Function:
		return e.evalFunction(n)
	case *UnaryExpr:
		return e.evalUnary(n)
	case *Comparison:
		return e.evalComparison(n)
	default:
		return nil, fmt.Errorf("unknown node type: %T", node)
	}
}

func (e *Interpreter) evalBinaryExpr(node *BinaryExpr) (any, error) {
	left, err := e.execute(node.Left)
	if err != nil {
		return nil, err
	}
	right, err := e.execute(node.Right)
	if err != nil {
		return nil, err
	}
	l, ok := left.(float64)
	if !ok {
		return nil, fmt.Errorf("expected float64, got %T", left)
	}
	r, ok := right.(float64)
	if !ok {
		return nil, fmt.Errorf("expected float64, got %T", right)
	}
	switch node.Op {
	case ADD:
		return l + r, nil
	case SUB:
		return l - r, nil
	case MUL:
		return l * r, nil
	case DIV:
		if r == 0 {
			return nil, fmt.Errorf("zero division error")
		}
		return l / r, nil
	case EXP:
		return math.Pow(l, r), nil
	default:
		return nil, fmt.Errorf("unknown binary operation: %d", node.Op)
	}
}

func (e *Interpreter) evalLiteral(node *Literal) (any, error) {
	switch node.Kind {
	case NUMBER:
		return strconv.ParseFloat(strings.Replace(node.Value, ",", ".", -1), 10)
	case STRING:
		return node.Value, nil
	default:
		return nil, fmt.Errorf("unknown literal type: %d", node.Kind)
	}
}

func (e *Interpreter) evalIdent(node *Ident) (any, error) {
	name := node.Name
	value, ok := e.variables[name]
	if !ok {
		return nil, fmt.Errorf("variable '%s' not found", name)
	}
	return value, nil
}

func (e *Interpreter) evalFunction(node *Function) (any, error) {
	funcName := node.Name
	function, ok := e.functions[funcName]
	if !ok {
		return nil, fmt.Errorf("function '%s' not found", funcName)
	}
	args := make([]any, len(node.Args))
	for i, arg := range node.Args {
		argument, err := e.execute(arg)
		if err != nil {
			return nil, err
		}
		args[i] = argument
	}
	return function(args...)
}

func (e *Interpreter) evalUnary(node *UnaryExpr) (any, error) {
	res, err := e.execute(node.Left)
	if err != nil {
		return nil, err
	}
	val, ok := res.(float64)
	if !ok {
		return nil, fmt.Errorf("expected float64, got %T", res)
	}
	switch node.Op {
	case ADD:
		return val, nil
	case SUB:
		return -val, nil
	default:
		return nil, fmt.Errorf("unknown unary operator: %d", node.Op)
	}
}

func (e *Interpreter) evalComparison(node *Comparison) (any, error) {
	left, err := e.execute(node.Left)
	if err != nil {
		return nil, err
	}
	right, err := e.execute(node.Right)
	if err != nil {
		return nil, err
	}

	if node.Op == EQ {
		return left == right, nil
	}
	switch l := left.(type) {
	case float64:
		r, ok := right.(float64)
		if !ok {
			return nil, fmt.Errorf("can't compare float64 and %T", right)
		}
		return compare(l, r, node.Op)
	case string:
		r, ok := right.(string)
		if !ok {
			return nil, fmt.Errorf("can't compare string and %T", right)
		}
		return compare(l, r, node.Op)
	default:
		return nil, fmt.Errorf("unknown comparable type %T", left)
	}
}

func compare[T float64 | string](left, right T, op TokenType) (bool, error) {
	switch op {
	case LT:
		return left < right, nil
	case GT:
		return left > right, nil
	case LTE:
		return left <= right, nil
	case GTE:
		return left >= right, nil
	}
	return false, fmt.Errorf("unexpected comparison token: %d", op)
}
