// Collection of different grammars used by tests
package testgrammars

import (
	"bnf"
)

// <Expr> ::= "A" "::=" <T>
// <T>    ::= "C"
func LinearMultiRule() bnf.Grammar {
	return bnf.Grammar {
		Rules: []bnf.Rule {
			{
				Head: bnf.SymbolNonTerminal{
					Name: "Expr",
				},
				Tail: bnf.Substitution{
					Sequences: []bnf.Sequence{
						{
							Symbols: []bnf.Symbol{
								bnf.SymbolTerminal {
									Name: "A",
								},
								bnf.SymbolTerminal {
									Name: "::=",
								},
								bnf.SymbolNonTerminal {
									Name: "T",
								},
							},
						},
					},
				},

			},
			{
				Head: bnf.SymbolNonTerminal{
					Name: "T",
				},
				Tail: bnf.Substitution{
					Sequences: []bnf.Sequence{
						{
							Symbols: []bnf.Symbol{
								bnf.SymbolTerminal {
									Name: "C",
								},
							},
						},
					},
				},
			},
		},
	}
}

// <A> ::= '"' <T> "^"
// <T> ::= "B" | "C" | ""
func Or() bnf.Grammar {
	return bnf.Grammar{
		Rules: []bnf.Rule{
			{
				Head: bnf.SymbolNonTerminal{
					Name: "A",
				},
				Tail: bnf.Substitution{
					Sequences: []bnf.Sequence{
						{
							Symbols: []bnf.Symbol{
								bnf.SymbolTerminal{
									Name: "\"",
								},
								bnf.SymbolNonTerminal{
									Name: "T",
								},
								bnf.SymbolTerminal{
									Name: "^",
								},
							},
						},
					},
				},

			},
			{
				Head: bnf.SymbolNonTerminal{
					Name: "T",
				},
				Tail: bnf.Substitution{
					Sequences: []bnf.Sequence{
						{
							Symbols: []bnf.Symbol{
								bnf.SymbolTerminal {
									Name: "B",
								},
							},
						},
						{
							Symbols: []bnf.Symbol{
								bnf.SymbolTerminal {
									Name: "C",
								},
							},
						},
						{
							Symbols: []bnf.Symbol{
								bnf.SymbolNothing {},
							},
						},
					},
				},
			},
		},
	}
}

// <T> ::= "_" | <T> "B"
func LRecursive() bnf.Grammar {
	return bnf.Grammar{
		Rules: []bnf.Rule{
			{
				Head: bnf.SymbolNonTerminal{
					Name: "T",
				},
				Tail: bnf.Substitution{
					Sequences: []bnf.Sequence{
						{
							Symbols: []bnf.Symbol{
								bnf.SymbolTerminal{
									Name: "_",
								},
							},
						},
						{
							Symbols: []bnf.Symbol{
								bnf.SymbolNonTerminal{
									Name: "T",
								},
								bnf.SymbolTerminal{
									Name: "B",
								},
							},
						},
					},
				},

			},
		},
	}
}

// <T>  ::= <T> "B" | "_"
//
//        |
//        V
//
// <T>  ::= "_" <Ta>
// <Ta> ::= "" | "B" <Ta>
func ResolvedLRecursive() bnf.Grammar {

	return bnf.Grammar{
		Rules: []bnf.Rule{
			{
				Head: bnf.SymbolNonTerminal{
					Name: "T",
				},
				Tail: bnf.Substitution{
					Sequences: []bnf.Sequence{
						{
							Symbols: []bnf.Symbol{
								bnf.SymbolTerminal{
									Name: "_",
								},
								bnf.SymbolNonTerminal{
									Name: "Ta",
								},
							},
						},
					},
				},

			},
			{
				Head: bnf.SymbolNonTerminal{
					Name: "Ta",
				},
				Tail: bnf.Substitution{
					Sequences: []bnf.Sequence{
						{
							Symbols: []bnf.Symbol{
								bnf.SymbolNothing {},
							},
						},
						{
							Symbols: []bnf.Symbol{
								bnf.SymbolTerminal {
									Name: "B",
								},
								bnf.SymbolNonTerminal{
									Name: "Ta",
								},
							},
						},
					},
				},
			},
		},
	}
}

