package bnf

import (
	"strings"
)

// Possible return values of Symbol.Type()
const (
	TypeNonTerminal int = iota
	TypeTerminal
	TypeNothing // Nothing
)

// Interface used to provide fake inhertance among possible BNF Type
// via type casting
type Symbol interface  {
	Type() int
	// To simplify Rule printing made Symbol by default be
	// compliant with fmt.Stringer interface
	String() string
}


type SymbolNonTerminal struct {
	Name string
}

func (nt SymbolNonTerminal) Type() int {
	return TypeNonTerminal
}

func (nt SymbolNonTerminal) String() string {
	sb := strings.Builder{}
	sb.WriteByte('<')
	sb.WriteString(nt.Name)
	sb.WriteByte('>')
	return sb.String()
}


type SymbolTerminal struct {
	Name string
}

func (nt SymbolTerminal) Type() int {
	return TypeTerminal
}

func (t SymbolTerminal) String() string {
	sb := strings.Builder{}
	res := t.Name
	res = strings.ReplaceAll(res, "\\", "\\\\")
	res = strings.ReplaceAll(res, "\t", "\\t")
	res = strings.ReplaceAll(res, "\n", "\\n")
	res = strings.ReplaceAll(res, "\r", "\\r")

	if strings.ContainsRune(res, '"') && !strings.ContainsRune(res, '\'') {
		sb.WriteByte('\'')
		sb.WriteString(res)
		sb.WriteByte('\'')
		return sb.String()
	}

	sb.WriteByte('"')

	res = strings.ReplaceAll(res, "\"", "\\\"")
	sb.WriteString(res)

	sb.WriteByte('"')
	return sb.String()
}


type SymbolNothing struct {
}

func (n SymbolNothing) Type() int {
	return TypeNothing
}

func (n SymbolNothing) String() string {
	return "\"\""
}


type Sequence struct {
	Symbols []Symbol
}


type Substitution struct {
	Sequences []Sequence
}


type Rule struct {
	Head SymbolNonTerminal
	Tail Substitution
}

func (r Rule) String() string {
	sb := strings.Builder{}
	sb.WriteString(r.Head.String())
	sb.WriteString(" ::=")
	for i, sequence := range r.Tail.Sequences {
		if i > 0 {
			sb.WriteString(" |")
		}
		for _, symbol := range sequence.Symbols {
			sb.WriteString(" ")
			sb.WriteString(symbol.String())
		}
	}

	return sb.String()
}


type Grammar struct {
	Rules []Rule
}

func (g Grammar) String() string {

	rulesCount := len(g.Rules)
	if rulesCount == 0 {
		return ""
	}

	sb := strings.Builder{}
	sb.WriteString(g.Rules[0].String())

	for i := 1; i < rulesCount; i++ {
		sb.WriteByte('\n')
		sb.WriteString(g.Rules[i].String())
	}

	return sb.String()
}
