package fromcst

import (
	"fmt"
	"sort"
	"cst"
	"bnf"
)

// Stack based tree LR traverse traversing cst
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

	stack := []cst.Node{}
	stackPush := func(node cst.Node) {
		stack = append(stack, node)
	}

	stackPop := func() (cst.Node, bool) {
		lenStack := len(stack)
		if lenStack == 0 {
			return nil, false
		}

		res := stack[lenStack-1]
		stack = stack[:lenStack-1]
		return res, true
	}

	stackIsEmpty := func() bool {
		return len(stack) == 0
	}

	stackPush(root)

	for !stackIsEmpty() {
		node, _ := stackPop()

		lenChilds := len(node.Childs())

		if isSearchedType(node.Type()) {
			err := fOnNode(node)
			if err != nil {
				return err
			}
			continue
		}

		for i := lenChilds - 1; i >= 0; i-- {
			nodeChild := node.Childs()[i]

			if isTypeIgnored(nodeChild.Type()) {
				continue
			}

			stackPush(nodeChild)
		}
	}

	return nil
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

	BuiltInLiteralType int
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

func (b BNFCSTtoASTBindings) parseSymbol(symbol cst.Node, str string) (*bnf.Symbol, error) {
	var res bnf.Symbol

	searchComplete := fmt.Errorf("headSearchComplete")

	doOnTerminalName := func(termNameNode cst.Node) error {
		if name := nodeName(termNameNode, str); len(name) == 0 {
			res = bnf.SymbolNothing{}
		} else {
			res = bnf.SymbolTerminal{Name: name}
		}
		return searchComplete
	}

	for _, symbolTermType := range b.SymbolTerminalTypes {
		err := b.lrTraverse(symbol, symbolTermType, doOnTerminalName)
		if err != nil {
			return &res, nil
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

	searchComplete := fmt.Errorf("headSearchComplete")

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
		RuleHeadType: 16,
		RuleTailType: 5,
		ExpressionSequencesType: 7,
		SequencesSymbolType: 9,
		SymbolTerminalTypes: []int {11, 12},
		SymbolNonTerminalType: 16,
		BuiltInLiteralType: -1,
	}
}