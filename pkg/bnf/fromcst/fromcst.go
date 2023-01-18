package fromcst

import (
	"bnf"
	"cst"
	"parser"
	"errors"
	"fmt"
	"strings"
	"sort"
)

type flexVisitor struct {
	isTypeIgnored func(int) bool
	isSearchedType func(int) bool
	fOnNode func(cst.Node) error

	err *error
}

func (v flexVisitor) Visit(node cst.Node) (w cst.Visitor) {
	if *v.err != nil {
		return nil
	}

	if v.isTypeIgnored(node.Type()) {
		return nil
	}

	if v.isSearchedType(node.Type()) {
		(*v.err) = v.fOnNode(node)
		return nil
	}

	return v
}

// depth-first cst search
// When meets isSearchedType(Node.type) == true executes fOnNode on node
// If fOnNode returned error immediately returns with fOnNode's error
// Childs of node on which fOnNode executed are ignored
// When Traverse complete returns nil
// Nodes on isTypeIgnored(Node.type) == true are ignored in traverse
// This function used by all functions in chain from ToAST
func lrTraverse(root cst.Node,
                isTypeIgnored func(int) bool,
                isSearchedType func(int) bool,
                fOnNode func(cst.Node) error) error {

	var err error = nil
	var v flexVisitor
	v.err = &err
	v.isTypeIgnored = isTypeIgnored
	v.isSearchedType = isSearchedType
	v.fOnNode = fOnNode

	cst.Walk(v, root)

	return err
}

func nodeName(node cst.Node, str string) string {
	return str[node.Pos():node.End()]
}

type BNFCSTtoASTBindings struct {
	IgnoreNodeTypes []int

	// nonterminal containing full singular rule syntax
	// (i.e. single meaningfull BNF line)
	RuleType int

	// first accurance of this nonterminal in RuleType will be threated as
	// head for rule
	// (i.e. full text inside < > (not including them) to the left from '::=')
	RuleHeadType int

	// first accurance of this nonterminal in RuleType will be threated as
	// tail for rule
	// (i.e. full meaningfull expression to the right from '::=')
	RuleTailType int

	ExpressionSequencesType int

	SequencesSymbolType int

	SymbolTerminalTypes []int

	SymbolNonTerminalType int

	// Content of this type (single character striung) will be used to replace
	// with corresponding value from EscapeMapping
	EscapeCharacterType int

	// Map character to it's corresponding value
	EscapeMapping map[string]string
}


// Handiy wrapper around lrTraverse
// When meets Node.type == untilType executes f
// If f returned error immediately returns with f's error
// Ignores Nodes with Node.Type ∈ b.IgnoreNodeTypes
// When Traverse complete returns nil
// This function used by all BNF's constructors
func (b BNFCSTtoASTBindings) lrTraverse(root cst.Node,
	                                    untilType int,
	                                    f func(cst.Node) error,
	                                    ) error {

	isNodeIgnored := func(nodeType int) bool {
		// Equivalent to sorted slice contains
		sr := sort.SearchInts(b.IgnoreNodeTypes, nodeType);
		return (sr < len(b.IgnoreNodeTypes) && b.IgnoreNodeTypes[sr] == nodeType)
	}

	isSearchedType := func(nodeType int) bool {
		return nodeType == untilType
	}

	return lrTraverse(root, isNodeIgnored, isSearchedType, f)
}

func (b BNFCSTtoASTBindings) parseEscape(escSeq cst.Node, str string) ([]byte, error) {
	var res []byte
	searchComplete := errors.New("headSearchComplete")
	doOnEscapedChar := func(escChar cst.Node) error {
		escStr := str[escChar.Pos():escChar.End()]
		resStr, ok := b.EscapeMapping[escStr]
		if !ok {
			return fmt.Errorf("unkknown escape character: %s", escStr)
		}
		res = []byte(resStr)
		return searchComplete
	}

	err := b.lrTraverse(escSeq, b.EscapeCharacterType, doOnEscapedChar)
	if err == nil {
		return nil, fmt.Errorf("not found escape character")
	}
	if err == searchComplete {
		return res, nil
	}
	return nil, err
}

// manually construct terminal literal
func (b BNFCSTtoASTBindings) parseTerminalName(symbol cst.Node, str string) (string, error) {
	sb := strings.Builder{}

	fOnNode := func(node cst.Node) error {
		if node.Type() == parser.BuiltinTerminal {
			sb.Write([]byte(str[node.Pos():node.End()]))
			return nil
		}

		escape, err := b.parseEscape(node, str)
		if err != nil {
			return err
		}
		sb.Write(escape)
		return nil
	}

	err := lrTraverse(
		symbol,
		func(int) bool {return false},
		func(n int) bool {return n == 16 || n == parser.BuiltinTerminal},
		fOnNode)

	return sb.String(), err
}

func (b BNFCSTtoASTBindings) parseSymbol(symbol cst.Node, str string) (*bnf.Symbol, error) {
	var res bnf.Symbol

	searchComplete := errors.New("headSearchComplete")

	doOnTerminalName := func(termNameNode cst.Node) error {
		name, err := b.parseTerminalName(termNameNode, str);
		if err != nil {
			return err
		}
		if len(name) == 0 {
			res = bnf.SymbolNothing{}
			return searchComplete
		}
		res = bnf.SymbolTerminal{Name: name}
		return searchComplete
	}

	for _, symbolTermType := range b.SymbolTerminalTypes {
		err := b.lrTraverse(symbol, symbolTermType, doOnTerminalName)
		if err == searchComplete {
			return &res, nil
		}
		if err != nil {
			return nil, err
		}
	}

	doOnNonTerminalName := func(nontermNameNode cst.Node) error {
		res = bnf.SymbolNonTerminal{Name: nodeName(nontermNameNode, str)}
		return searchComplete
	}

	err := b.lrTraverse(symbol, b.SymbolNonTerminalType, doOnNonTerminalName)
	if err == searchComplete {
		return &res, nil
	}

	return nil, fmt.Errorf("could not determine terminal type of `%s`",
	                       nodeName(symbol, str))
}

func (b BNFCSTtoASTBindings) parseSequence(sequence cst.Node, str string,
										   ) (*bnf.Sequence, error) {
	var res bnf.Sequence

	doOnSymbol := func(symbolNode cst.Node) error {
		symbol, err := b.parseSymbol(symbolNode, str)
		if err != nil {
			return err
		}
		res.Symbols = append(res.Symbols, *symbol)
		return nil
	}
	err := b.lrTraverse(sequence, b.SequencesSymbolType, doOnSymbol)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

func (b BNFCSTtoASTBindings) parseSubstitution(expression cst.Node,
                                               str string,
                                               ) (*bnf.Substitution, error) {
	var res bnf.Substitution

	doOnSequence := func(sequenceNode cst.Node) error {
		sequence, err := b.parseSequence(sequenceNode, str)
		if err != nil {
			return err
		}
		res.Sequences = append(res.Sequences, *sequence)
		return nil
	}

	err := b.lrTraverse(expression, b.ExpressionSequencesType, doOnSequence)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

func (b BNFCSTtoASTBindings) parseRule(rule cst.Node, str string) (*bnf.Rule, error) {
	var res bnf.Rule

	// str := nodeName(rule, str)
	// res.Head.Name = str

	searchComplete := errors.New("headSearchComplete")

	doOnHead := func(ruleNameNode cst.Node) error {
		res.Head.Name = nodeName(ruleNameNode, str)
		return searchComplete
	}

	err := b.lrTraverse(rule, b.RuleHeadType, doOnHead)
	if err != searchComplete {
		// a single way of how lrTravese could return an error is via f's error
		// so if err != f's error it is nil doOnHead only returns single type
		// of error

		return nil, fmt.Errorf("could not find rule name in `%s`",
		                       nodeName(rule, str))
	}

	doOnSubstitution := func(ruleNode cst.Node) error {
		substitution, err := b.parseSubstitution(ruleNode, str)
		if err != nil {
			return err
		}
		res.Tail = *substitution
		return searchComplete
	}

	err = b.lrTraverse(rule, b.RuleTailType, doOnSubstitution)
	if err != searchComplete {
		if err == nil {
			return nil, fmt.Errorf("could not find rule substitution in `%s`",
			                       nodeName(rule, str))
		}

		return nil, err
	}

	return &res, nil
}

func (b BNFCSTtoASTBindings) ToAST(root cst.Node, str string) (*bnf.Grammar, error) {

	sort.Sort(sort.IntSlice(b.IgnoreNodeTypes))

	var res bnf.Grammar

	doOnRule := func(ruleNode cst.Node) error {
		rule, err := b.parseRule(ruleNode, str)
		if err != nil {
			return err
		}
		res.Rules = append(res.Rules, *rule)
		return nil
	}

	err := b.lrTraverse(root, b.RuleType, doOnRule)

	if err != nil {
		return nil, err
	}

	return &res, nil
}

func SelfCSTtoASTBindings() BNFCSTtoASTBindings {
	return BNFCSTtoASTBindings{
		IgnoreNodeTypes: []int{22},
		RuleType: 4,
		RuleHeadType: 18,
		RuleTailType: 5,
		ExpressionSequencesType: 7,
		SequencesSymbolType: 9,
		SymbolTerminalTypes: []int {11, 12},
		SymbolNonTerminalType: 18,
		EscapeCharacterType: 17,
		EscapeMapping: map[string]string{
			 "t": "\t",
			 "n": "\n",
			 "r": "\r",
			"\"": "\"",
			"\\": "\\",
		},
	}
}
