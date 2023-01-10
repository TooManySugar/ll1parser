// Example showing how BNF's own syntax AST repesentation can be converted to
// parser's table and string form which could be parsed using previously
// generated table by parser to CST for further conversing back to BNF's syntax
// AST repesentation equal to iniitial one
package main

import (
	"fmt"
	"os"

	"bnf"
	"bnf/tablegen"
	"bnf/fromcst"
	"parser"
)

func isNonTerminalsEqual(a bnf.SymbolNonTerminal,
                         b bnf.SymbolNonTerminal) bool {

	return a.Name == b.Name
}

func isTerminalsEqual(a bnf.SymbolTerminal, b bnf.SymbolTerminal) bool {
	return a.Name == b.Name
}

func isSymbolsEqual(a bnf.Symbol, b bnf.Symbol) bool {
	if a.Type() != b.Type() {
		return false
	}

	switch a.Type() {
	case bnf.TypeNonTerminal: {
		return isNonTerminalsEqual(a.(bnf.SymbolNonTerminal),
		                           b.(bnf.SymbolNonTerminal))
	}
	case bnf.TypeTerminal: {
		return isTerminalsEqual(a.(bnf.SymbolTerminal),
		                        b.(bnf.SymbolTerminal))
	}
	case bnf.TypeNothing: {
		return true
	}
	default: {}
	}

	fmt.Fprintln(os.Stderr, "ERROR: unknown BNF type:", a.Type())
	os.Exit(1)
	panic("unreachable")
}

func isSequencesEqual(a bnf.Sequence, b bnf.Sequence) bool {
	if len(a.Symbols) != len(b.Symbols) {
		return false
	}

	for i := range a.Symbols {
		if !isSymbolsEqual(a.Symbols[i], b.Symbols[i]) {
			return false
		}
	}

	return true
}

func isSubstitutionsEqual(a bnf.Substitution, b bnf.Substitution) bool {
	if len(a.Sequences) != len(b.Sequences) {
		return false
	}

	for i := range a.Sequences {
		if !isSequencesEqual(a.Sequences[i], b.Sequences[i]) {
			return false
		}
	}
	return true
}

func isRulesEqual(a bnf.Rule, b bnf.Rule) bool {

	if !isNonTerminalsEqual(a.Head, b.Head) {
		return false
	}

	return isSubstitutionsEqual(a.Tail, b.Tail)
}

func isGrammarsEqual(a bnf.Grammar, b bnf.Grammar) bool {
	if len(a.Rules) != len(b.Rules) {
		return false
	}

	for i := range a.Rules {
		if !isRulesEqual(a.Rules[i], b.Rules[i]) {
			return false
		}
	}

	return true
}

func main() {
	reference := bnf.SelfGrammar()
	bnfStr := reference.String()

	parserTable, parserTableNames, ok := tablegen.ToParserTable(reference)
	if !ok {
		fmt.Fprintln(os.Stderr,
		             "ERROR: can't build Parser table from grammar")
		os.Exit(1)
	}

	parser := parser.NewLL1Parser(*parserTable, *parserTableNames)

	cst, _, err := parser.Parse(bnfStr)
	if err != nil {
		fmt.Fprintln(os.Stderr,
		             "ERROR: can't parse input:", err.Error())
		os.Exit(1)
	}

	result, err := fromcst.SelfCSTtoASTBindings().ToAST(cst, bnfStr)
	if err != nil {
		fmt.Fprintln(os.Stderr,
		             "ERROR: can't build AST from CST:", err.Error())
		os.Exit(1)
	}

	if !isGrammarsEqual(reference, *result) {
		fmt.Fprintln(os.Stderr, "ERROR: final AST not equal to initial")
		os.Exit(1)
	}

	fmt.Println("Resulted AST equal to initial")

	os.Exit(0)
}
