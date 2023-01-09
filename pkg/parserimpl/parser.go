// Implementation of LL1Parser interface
package parserimpl

import (
	"fmt"
	"cst"
	"parserinterfaces"
)

// Parser operands types
const (
	nonTerminal int = iota
	terminal
	function
	opEOS
	opChar
)

// builin terminal types
const (
	//
	builtinNothing int = iota-3
	// for now builtin rule <EOL> ::= "\r\n" | "\n"
	builtinEOL
	builtinTerminal
)

// make it expose to be readable from table generator(s)
// but keep lowercase in enum for consistency within
const BuiltinEOL = builtinEOL

// Parser operands inheritance work around
type ParserOp interface {
	parserOpType() int
}

type nonTerminal_t struct {
	name int
}

func OpNonTerminal(n int) ParserOp {
	return nonTerminal_t{ name: n }
}

func (nt nonTerminal_t) parserOpType() int {
	return nonTerminal
}

func (nt nonTerminal_t) GetName() int {
	return nt.name
}


type terminal_t struct {
	value string
}

func OpTerminal(s string) ParserOp {
	return terminal_t{ value: s }
}

func (t terminal_t) parserOpType() int {
	return terminal
}

func (t terminal_t) GetValue() string {
	return t.value
}



type function_t struct {
	Name int
	Pos int
	Amount int
}

func (f function_t) parserOpType() int {
	return function
}

func (f function_t) GetName() int {
	return f.Name
}

func (f function_t) GetPos() int {
	return f.Pos
}

func (f function_t) GetAmount() int {
	return f.Amount
}


type opEOS_t struct {
}

func (e opEOS_t) parserOpType() int {
	return opEOS
}


type opChar_t struct {
	Value byte
}

func (c opChar_t) parserOpType() int {
	return opChar
}

func (c opChar_t) GetValue() byte {
	return c.Value
}

type ParserOpList []ParserOp

func NewParserOpList(parserOps ...ParserOp) []ParserOp {
	return append([]ParserOp{}, parserOps...)
}

type ll1parser_t struct {
	table map[int]map[byte][]ParserOp
	names map[int]string
}

func NewLL1Parser(table map[int]map[byte][]ParserOp,
	              names map[int]string) parserinterfaces.LL1Parser {
	return ll1parser_t{table: table, names: names}
}


type ll1parserOpStack struct {
	stack []ParserOp
}

func (opStack *ll1parserOpStack) Pop() (*ParserOp, bool) {
	stackLen := len(opStack.stack)

	if stackLen == 0 {
		return nil, false
	}

	var res ParserOp
	res = opStack.stack[stackLen - 1]
	opStack.stack = opStack.stack[:stackLen - 1]
	return &res, true
}

func (opStack *ll1parserOpStack) Push(pOp ParserOp) {
	opStack.stack = append(opStack.stack, pOp)
}

func (opStack *ll1parserOpStack) Len() int {
	return len(opStack.stack)
}

type ll1parserProdStack struct {
	stack []cst.Node
}

func (prodStack *ll1parserProdStack) Pop() (cst.Node, bool) {
	stackLen := len(prodStack.stack)

	if stackLen == 0 {
		return nil, false
	}

	var res cst.Node
	res = prodStack.stack[stackLen - 1]
	prodStack.stack = prodStack.stack[:stackLen - 1]
	return res, true
}

func (prodStack *ll1parserProdStack) Push(node cst.Node) {
	prodStack.stack = append(prodStack.stack, node)
}

func (prodStack ll1parserProdStack) Len() int {
	return len(prodStack.stack)
}


func (p ll1parser_t) Parse(in parserinterfaces.LL1ParserReader) (cst.Node, *map[int]string, error) {

	if len(p.table) == 0 {
		return nil, nil, fmt.Errorf("empty parsing table")
	}

	if _, found := p.table[0]; !found {
		return nil, nil,
			fmt.Errorf("can't start parsing: no rule for base entry point - 0")
	}

	charCode := func (char byte) string {
		return fmt.Sprintf("'%c' (%d)", char, char)
	}

	nodeTypeName := func (name_id int) string {
		val, ok := p.names[name_id]
		if !ok {
			return fmt.Sprintf("Unknown_%d", name_id)
		}
		return val
	}

	var opStack ll1parserOpStack
	var prodStack ll1parserProdStack

	opStack.Push(opEOS_t{})
	opStack.Push(nonTerminal_t{name: 0})


	for opStack.Len() > 0 {
		// fmt.Println(opStack.stack, fmt.Sprintf("`%c`", in.Peek()))

		op, _ := opStack.Pop()

		// fmt.Println( (*op).ParserOpType() )

		switch (*op).parserOpType() {
		case nonTerminal: {
			nt, ok := (*op).(nonTerminal_t)
			if !ok {
				panic("can't cast nonTerminal op to it's type")
			}

			nodeType := nt.GetName()

			if nodeType > -1 {
				ruleMap, ok := p.table[nodeType]
				if !ok {
					// Table error
					return nil, nil, fmt.Errorf("no rules for non terminal: %s",
					                            nodeTypeName(nodeType))
				}

				opsToPush, ok := ruleMap[in.Peek()]
				if !ok {
					// Parsing error
					return nil, nil,
						fmt.Errorf("no rules for %s and non terminal op <%s>",
						           charCode(in.Peek()),
						           nodeTypeName(nodeType))
				}

				opStack.Push(function_t{Name: nodeType, Pos: in.Pos(), Amount: len(opsToPush)})

				for i := len(opsToPush) - 1; i >= 0; i -- {
					opStack.Push(opsToPush[i])
				}
			} else {
				switch nodeType {
				case -2:
					input := in.Peek()

					switch input {
					case '\n':
						opStack.Push(function_t{Name: nodeType, Pos: in.Pos(), Amount: 1})
						opStack.Push(OpTerminal("\n"))
					case '\r':
						opStack.Push(function_t{Name: nodeType, Pos: in.Pos(), Amount: 2})
						opStack.Push(OpTerminal("\n"))
						opStack.Push(OpTerminal("\r"))
					default:
						// Parsing error
						return nil, nil,
						fmt.Errorf("no rules for %s and builtin terminal op <EOL>",
						           charCode(in.Peek()))
					}
				default:
					return nil, nil,
						fmt.Errorf("unknown built in type: %d", nodeType)
				}
			}
		}
		case terminal: {
			t, ok := (*op).(terminal_t)
			if !ok {
				panic("can't cast terminal op to it's type")
			}

			tValue := t.GetValue()

			opStack.Push(function_t{Name: builtinTerminal, Pos: in.Pos(), Amount: len(tValue)})
			for i := len(tValue) - 1; i >= 0; i -- {
				opStack.Push(opChar_t{Value: tValue[i]})
			}
		}
		case function: {
			f, ok := (*op).(function_t)
			if !ok {
				panic("can't cast function op to it's type")
			}

			childs := []cst.Node{}

			if f.Amount == 0 {
				childs = append(childs, cst.NewNode(-3, in.Pos(), in.Pos(), nil))
				prodStack.Push(cst.NewNode(f.Name, in.Pos(), in.Pos(), childs))
				break
			}

			for i := 0; i < f.Amount; i++ {
				node, ok := prodStack.Pop()
				if !ok {
					panic("trying to pop from empty stack")
				}
				if f.Name == builtinTerminal {
					if node.Type() != builtinTerminal {
						panic("Tring to combine chars of terminal from non chars type")
					}
				} else {
					childs = append([]cst.Node{node}, childs...)
				}
			}
			if f.Name == builtinEOL {
				childs = []cst.Node{}
			}
			prodStack.Push(cst.NewNode(f.Name, f.GetPos(), in.Pos(), childs))
		}
		case opEOS: {
			if in.Peek() != byte(0) {
				// fmt.Println("expected end of input got: ", in.Peek())
				// Parsing error Unexpected EOF
				return nil, nil, fmt.Errorf("expected end of input got %s",
				                            charCode(in.Peek()))
			}
			// fmt.Println("Parsed successfully")

			n, ok := prodStack.Pop()
			if !ok {
				// Table error too
				return nil, nil, fmt.Errorf("prod stack empty")
			}


			ret_names := p.names
			ret_names[builtinTerminal] = "_literal"
			ret_names[builtinEOL]      = "_endofline"
			ret_names[builtinNothing]  = "_nothing"

			return n, &ret_names, nil
		}
		case opChar: {
			c, ok := (*op).(opChar_t)
			if !ok {
				panic("can't cast char op to it's type")
			}

			if c.GetValue() != in.Peek() {
				// Parsing error
				return nil, nil, fmt.Errorf("expected char %s, got %s", charCode(c.GetValue()), charCode(in.Peek()))
			}

			prodStack.Push(cst.NewNode(builtinTerminal, in.Pos(), in.Pos() + 1, nil))
			in.Move()
		}
		default: {
			panic(fmt.Sprint("unknown terminal type:", (*op).parserOpType()))
		}}

	}

	panic("unreachable")
}
