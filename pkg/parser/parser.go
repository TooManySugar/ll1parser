// Implementation of LL1Parser interface
package parser

import (
	"fmt"
	"bytes"
	"io"
	"cst"
)

// Parser operands types
const (
	opTypeNonTerminal int = iota
	opTypeTerminal
	opTypeFunction
	opTypeEOS
	opTypeChar
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
	return opTypeNonTerminal
}

func (nt nonTerminal_t) Name() int {
	return nt.name
}


type terminal_t struct {
	value string
}

func OpTerminal(s string) ParserOp {
	return terminal_t{ value: s }
}

func (t terminal_t) parserOpType() int {
	return opTypeTerminal
}

func (t terminal_t) Value() string {
	return t.value
}


type function_t struct {
	name int
	pos int
	amount int
}

func opFunction(name int, pos int, amount int) ParserOp {
	return function_t {
		name: name,
		pos: pos,
		amount: amount,
	}
}

func (f function_t) parserOpType() int {
	return opTypeFunction
}

func (f function_t) Name() int {
	return f.name
}

func (f function_t) Pos() int {
	return f.pos
}

func (f function_t) Amount() int {
	return f.amount
}


type opEOS_t struct {
}

func opEOS() ParserOp {
	return opEOS_t{}
}

func (e opEOS_t) parserOpType() int {
	return opTypeEOS
}


type opChar_t struct {
	value byte
}

func opChar(value byte) ParserOp {
	return opChar_t{
		value: value,
	}
}

func (c opChar_t) parserOpType() int {
	return opTypeChar
}

func (c opChar_t) Value() byte {
	return c.value
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
	              names map[int]string) LL1Parser {
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

type ll1parserScanner struct {
	src []byte
	offset int
	// lineOffset
	// onLineOffset
}

func (s *ll1parserScanner) peek() byte {
	if s.offset >= len(s.src) {
		return byte(0)
	}
	return s.src[s.offset]
}

func (s *ll1parserScanner) next() {
	s.offset++
}

func (s *ll1parserScanner) aPos() int {
	return s.offset
}

func charCode(char byte) string {
	return fmt.Sprintf("'%c' (%d)", char, char)
}

type realParser struct {
	table map[int]map[byte][]ParserOp
	names map[int]string
	scanner ll1parserScanner
	opStack ll1parserOpStack
	prodStack ll1parserProdStack
}

func (p *realParser) processTableNonTerminal(name int) error {
	nodeTypeName := func (name_id int) string {
		val, ok := p.names[name_id]
		if !ok {
			return fmt.Sprintf("Unknown_%d", name_id)
		}
		return val
	}

	ruleMap, ok := p.table[name]
	if !ok {
		// Table error
		return fmt.Errorf("no rules for non terminal: %s", nodeTypeName(name))
	}

	opsToPush, ok := ruleMap[p.scanner.peek()]
	if !ok {
		// Parsing error
		return fmt.Errorf("no rules for %s and non terminal op <%s>",
			           charCode(p.scanner.peek()),
			           nodeTypeName(name))
	}

	p.opStack.Push(
		opFunction(name, p.scanner.aPos(), len(opsToPush)))

	for i := len(opsToPush) - 1; i >= 0; i -- {
		p.opStack.Push(opsToPush[i])
	}

	return nil
}

func (p *realParser) processBuiltinNonTerminal(name int) error {
	switch name {
	case -2:
		input := p.scanner.peek()

		switch input {
		case '\n':
			p.opStack.Push(
				opFunction(name, p.scanner.aPos(), 1))
			p.opStack.Push(OpTerminal("\n"))
			return nil
		case '\r':
			p.opStack.Push(
				opFunction(name, p.scanner.aPos(), 2))
			p.opStack.Push(OpTerminal("\n"))
			p.opStack.Push(OpTerminal("\r"))
			return nil
		default:
			// Parsing error
			return fmt.Errorf("no rules for %s and builtin terminal op <EOL>",
			                  charCode(p.scanner.peek()))
		}
	}
	return fmt.Errorf("unknown built in type: %d", name)
}

func (p *realParser) processNonTerminal(nt nonTerminal_t) error {
	name := nt.Name()
	if name >= 0 {
		return p.processTableNonTerminal(name)
	} else {
		return p.processBuiltinNonTerminal(name)
	}
}

func (p *realParser) processTerminal(t terminal_t) {
	tValue := t.Value()

	p.opStack.Push(
		opFunction(builtinTerminal, p.scanner.aPos(), len(tValue)))
	for i := len(tValue) - 1; i >= 0; i -- {
		p.opStack.Push(opChar(tValue[i]))
	}
}

func (p *realParser) processFunction(f function_t) {
	name := f.Name()
	amount := f.Amount()
	if amount == 0 {
		p.prodStack.Push(cst.NewNode(name,
		                             p.scanner.aPos(),
		                             p.scanner.aPos(),
		                             []cst.Node{
		                                 cst.NewNode(builtinNothing,
		                                     p.scanner.aPos(),
		                                     p.scanner.aPos(),
		                                     nil),
		                             }))
		return
	}

	var childs []cst.Node

	discard := (name == builtinTerminal || name == builtinEOL)
	if !discard {
		childs = make([]cst.Node, amount)
	}

	for i := amount - 1; i >= 0; i-- {
		node, ok := p.prodStack.Pop()
		if !ok {
			panic("trying to pop from empty stack")
		}

		if !discard {
			childs[i] = node
			continue
		}

		if node.Type() != builtinTerminal {
			panic("Tring to combine chars of terminal from non chars type")
		}
	}

	p.prodStack.Push(cst.NewNode(name, f.Pos(), p.scanner.aPos(), childs))
}

func (p *realParser) processEOS() (cst.Node, *map[int]string, error) {
	if p.scanner.peek() != byte(0) {
		// fmt.Println("expected end of input got: ", in.Peek())
		// Parsing error Unexpected EOF
		return nil, nil, fmt.Errorf("expected end of input got %s",
									charCode(p.scanner.peek()))
	}
	// fmt.Println("Parsed successfully")

	n, ok := p.prodStack.Pop()
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

func (p *realParser) processChar(c opChar_t) error {
	if c.Value() != p.scanner.peek() {
		// Parsing error
		return fmt.Errorf("expected char %s, got %s",
		                  charCode(c.Value()),
		                  charCode(p.scanner.peek()))
	}

	p.prodStack.Push(
		cst.NewNode(builtinTerminal,
		            p.scanner.aPos(),
		            p.scanner.aPos() + 1,
		            nil))
	p.scanner.next()
	return nil
}

func (p realParser) parse() (cst.Node, *map[int]string, error) {

	p.opStack.Push(opEOS())
	p.opStack.Push(OpNonTerminal(0))

	for p.opStack.Len() > 0 {
		// fmt.Println(opStack.stack, fmt.Sprintf("`%c`", in.Peek()))

		op, _ := p.opStack.Pop()

		// fmt.Println( (*op).ParserOpType() )

		switch (*op).parserOpType() {
		case opTypeNonTerminal: {
			nt, ok := (*op).(nonTerminal_t)
			if !ok {
				panic("can't cast nonTerminal op to it's type")
			}

			err := p.processNonTerminal(nt)
			if err != nil {
				return nil, nil, err
			}
		}
		case opTypeTerminal: {
			t, ok := (*op).(terminal_t)
			if !ok {
				panic("can't cast terminal op to it's type")
			}

			p.processTerminal(t)
		}
		case opTypeFunction: {
			f, ok := (*op).(function_t)
			if !ok {
				panic("can't cast opTypeFunction op to it's type")
			}

			p.processFunction(f)
		}
		case opTypeEOS: {
			return p.processEOS()
		}
		case opTypeChar: {
			c, ok := (*op).(opChar_t)
			if !ok {
				panic("can't cast char op to it's type")
			}

			err := p.processChar(c)
			if err != nil {
				return nil, nil, err
			}
		}
		default: {
			panic(fmt.Sprint("unknown terminal type:", (*op).parserOpType()))
		}}

	}

	panic("unreachable")
}

func readSource(src any) ([]byte, error) {
	if src == nil {
		return nil, fmt.Errorf("empty source")
	}
	switch s := src.(type) {
	case string:
		return []byte(s), nil
	case []byte:
		return s, nil
	case *bytes.Buffer:
		// is io.Reader, but src is already available in []byte form
		if s != nil {
			return s.Bytes(), nil
		}
	case io.Reader:
		return io.ReadAll(s)
	}
	return nil, fmt.Errorf("invalid source")
}

func (p ll1parser_t) Parse(src any) (cst.Node, *map[int]string, error) {

	if len(p.table) == 0 {
		return nil, nil, fmt.Errorf("empty parsing table")
	}

	if _, found := p.table[0]; !found {
		return nil, nil,
			fmt.Errorf("can't start parsing: no rule for base entry point - 0")
	}

	text, err := readSource(src)
	if err != nil {
		return nil, nil, err
	}

	var rp realParser
	rp.table = p.table
	rp.names = p.names
	rp.scanner = ll1parserScanner{
		src: text,
		offset: 0,
	}

	return rp.parse()
}
