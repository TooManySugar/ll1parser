package bnf

// AST for defining Backus–Naur form expression set
//
// Taken from Wikipedia and modified to satisfy essential conditions
// to generate LL(1) Parser.
// Rule order important (as for current parser table gen algo) but not unique.
//
// This AST in BNF:
//
// Note: \ used to express line continuation on the next line,
//       Grammar not support multiline rules as it seen
//
//     <syntax>          ::= <opt-whitespace> <content> <more-lines>
//     <more-lines>      ::= "" | <EOL> <line> <more-lines>
//     <line>            ::= <opt-whitespace> <opt-content>
//
//     <opt-content>     ::= <content> | ""
//     <content>         ::= "<" <rule-name> ">" <opt-whitespace> "::=" \
//                           <opt-whitespace> <expression>
//
//     <expression>      ::= <list> <expression-tail>
//     <expression-tail> ::= "" | "|" <opt-whitespace> <list> <expression-tail>
//     <list>            ::= <term> <opt-whitespace> <list-tail>
//     <list-tail>       ::= "" | <term> <opt-whitespace> <list-tail>
//     <term>            ::= <literal> | "<" <rule-name> ">"
//     <literal>         ::= '"' <text1> '"' | "'" <text2> "'"
//     <text1>           ::= "" | <character1> <text1>
//     <text2>           ::= "" | <character2> <text2>
//     <character1>      ::= "'" | <character>
//     <character2>      ::= '"' | <character>
//     <character>       ::= <letter> | <digit> | <symbol> | <escape-sequence>
//     <escape-sequence> ::= "\\" <escaped-char>
//     <escaped-char>    ::= "t" | "n" | "r" | '"' | "\\"
//     <rule-name>       ::= <letter> <rule-name-tail>
//     <rule-name-tail>  ::= "" | <rule-char> <rule-name-tail>
//     <rule-char>       ::= <letter> | <digit> | "-"
//
//     <letter>          ::= "A" | "B" | "C" | "D" | "E" | "F" | "G" | "H" | \
//                           "I" | "J" | "K" | "L" | "M" | "N" | "O" | "P" | \
//                           "Q" | "R" | "S" | "T" | "U" | "V" | "W" | "X" | \
//                           "Y" | "Z" | "a" | "b" | "c" | "d" | "e" | "f" | \
//                           "g" | "h" | "i" | "j" | "k" | "l" | "m" | "n" | \
//                           "o" | "p" | "q" | "r" | "s" | "t" | "u" | "v" | \
//                           "w" | "x" | "y" | "z"
//
//     <digit>           ::= "0" | "1" | "2" | "3" | "4" | "5" | "6" | "7" | \
//                           "8" | "9"
//
//     <symbol>          ::= "|" | " " | "!" | "#" | "$" | "%" | "&" | "(" | \
//                           ")" | "*" | "+" | "," | "-" | "." | "/" | ":" | \
//                           ";" | ">" | "=" | "<" | "?" | "@" | "[" | "\" | \
//                           "]" | "^" | "_" | "`" | "{" | "}" | "~"
//
//     <opt-whitespace>  ::= " " <opt-whitespace> | ""
//     <EOL>             ::= "\n" | "\r\n"
//
func SelfGrammar() Grammar {
	return Grammar{
	   Rules: []Rule{
		   { //  0 <syntax>
			   Head: SymbolNonTerminal{
				   Name: "syntax",
			   },
			   Tail: Substitution{
				   Sequences: []Sequence{
					   {
						   Symbols: []Symbol{
							   SymbolNonTerminal{
								   Name: "opt-whitespace",
							   },
							   SymbolNonTerminal{
								   Name: "content",
							   },
							   SymbolNonTerminal{
								   Name: "more-lines",
							   },
						   },
					   },
				   },
			   },
		   },
		   { //  1 <more-lines>
			   Head: SymbolNonTerminal{
				   Name: "more-lines",
			   },
			   Tail: Substitution{
				   Sequences: []Sequence{
					   {
						   Symbols: []Symbol{
							   SymbolNothing {},
						   },
					   },
					   {
						   Symbols: []Symbol{
							   SymbolNonTerminal {
								   Name: "EOL",
							   },
							   SymbolNonTerminal {
								   Name: "line",
							   },
							   SymbolNonTerminal{
								   Name: "more-lines",
							   },
						   },
					   },
				   },
			   },
		   },
		   { //  2 <line>
			   Head: SymbolNonTerminal{
				   Name: "line",
			   },
			   Tail: Substitution{
				   Sequences: []Sequence{
					   {
						   Symbols: []Symbol{
							   SymbolNonTerminal {
								   Name: "opt-whitespace",
							   },
							   SymbolNonTerminal{
								   Name: "opt-content",
							   },
							   // SymbolNonTerminal {
							   // 	Name: "opt-whitespace",
							   // },
						   },
					   },
				   },
			   },
		   },
		   { //  3 <opt-content>
			   Head: SymbolNonTerminal{
				   Name: "opt-content",
			   },
			   Tail: Substitution{
				   Sequences: []Sequence{
					   {
						   Symbols: []Symbol{
							   SymbolNonTerminal {
								   Name: "content",
							   },
						   },
					   },
					   {
						   Symbols: []Symbol{
							   SymbolNothing {},
						   },
					   },
				   },
			   },
		   },
		   { //  4 <content>
			   Head: SymbolNonTerminal{
				   Name: "content",
			   },
			   Tail: Substitution{
				   Sequences: []Sequence{
					   {
						   Symbols: []Symbol{
							   SymbolTerminal {
								   Name: "<",
							   },
							   SymbolNonTerminal {
								   Name: "rule-name",
							   },
							   SymbolTerminal {
								   Name: ">",
							   },
							   SymbolNonTerminal {
								   Name: "opt-whitespace",
							   },
							   SymbolTerminal {
								   Name: "::=",
							   },
							   SymbolNonTerminal {
								   Name: "opt-whitespace",
							   },
							   SymbolNonTerminal {
								   Name: "expression",
							   },
						   },
					   },
				   },
			   },
		   },
		   { //  5 <expression>
			   Head: SymbolNonTerminal{
				   Name: "expression",
			   },
			   Tail: Substitution{
				   Sequences: []Sequence{
					   {
						   Symbols: []Symbol{
							   SymbolNonTerminal{
								   Name: "list",
							   },
							   SymbolNonTerminal{
								   Name: "expression-tail",
							   },
						   },
					   },
				   },
			   },
		   },
		   { //  6 <expression-tail>
			   Head: SymbolNonTerminal{
				   Name: "expression-tail",
			   },
			   Tail: Substitution{
				   Sequences: []Sequence{
					   {
						   Symbols: []Symbol{
							   SymbolNothing{},
						   },
					   },
					   {
						   Symbols: []Symbol{
							   SymbolTerminal{
								   Name: "|",
							   },
							   SymbolNonTerminal{
								   Name: "opt-whitespace",
							   },
							   SymbolNonTerminal{
								   Name: "list",
							   },
							   SymbolNonTerminal{
								   Name: "expression-tail",
							   },
						   },
					   },
				   },
			   },
		   },
		   { //  7 <list>
			   Head: SymbolNonTerminal{
				   Name: "list",
			   },
			   Tail: Substitution{
				   Sequences: []Sequence{
					   {
						   Symbols: []Symbol{
							   SymbolNonTerminal{
								   Name: "term",
							   },
							   SymbolNonTerminal{
								   Name: "opt-whitespace",
							   },
							   SymbolNonTerminal{
								   Name: "list-tail",
							   },
						   },
					   },
				   },
			   },
		   },
		   { //  8 <list-tail>
			   Head: SymbolNonTerminal{
				   Name: "list-tail",
			   },
			   Tail: Substitution{
				   Sequences: []Sequence{
					   {
						   Symbols: []Symbol{
							   SymbolNothing{},
						   },
					   },
					   {
						   Symbols: []Symbol{
							   SymbolNonTerminal{
								   Name: "term",
							   },
							   SymbolNonTerminal{
								   Name: "opt-whitespace",
							   },
							   SymbolNonTerminal{
								   Name: "list-tail",
							   },
						   },
					   },
				   },
			   },
		   },
		   { //  9 <term>
			   Head: SymbolNonTerminal{
				   Name: "term",
			   },
			   Tail: Substitution{
				   Sequences: []Sequence{
					   {
						   Symbols: []Symbol{
							   SymbolNonTerminal{
								   Name: "literal",
							   },
						   },
					   },{
						   Symbols: []Symbol{
							   SymbolTerminal{
								   Name: "<",
							   },
							   SymbolNonTerminal{
								   Name: "rule-name",
							   },
							   SymbolTerminal{
								   Name: ">",
							   },
						   },
					   },
				   },
			   },
		   },
		   { // 10 <literal>
			   Head: SymbolNonTerminal{
				   Name: "literal",
			   },
			   Tail: Substitution{
				   Sequences: []Sequence{
					   {
						   Symbols: []Symbol{
							   SymbolTerminal{
								   Name: "\"",
							   },
							   SymbolNonTerminal{
								   Name: "text1",
							   },
							   SymbolTerminal{
								   Name: "\"",
							   },
						   },
					   },
					   {
						   Symbols: []Symbol{
							   SymbolTerminal{
								   Name: "'",
							   },
							   SymbolNonTerminal{
								   Name: "text2",
							   },
							   SymbolTerminal{
								   Name: "'",
							   },
						   },
					   },
				   },
			   },
		   },
		   { // 11 <text1>
			   Head: SymbolNonTerminal{
				   Name: "text1",
			   },
			   Tail: Substitution{
				   Sequences: []Sequence{
					   {
						   Symbols: []Symbol{
							   SymbolNothing{},
						   },
					   },
					   {
						   Symbols: []Symbol{
							   SymbolNonTerminal{
								   Name: "character1",
							   },
							   SymbolNonTerminal{
								   Name: "text1",
							   },
						   },

					   },
				   },
			   },
		   },
		   { // 12 <text2>
			   Head: SymbolNonTerminal{
				   Name: "text2",
			   },
			   Tail: Substitution{
				   Sequences: []Sequence{
					   {
						   Symbols: []Symbol{
							   SymbolNothing{},
						   },
					   },
					   {
						   Symbols: []Symbol{
							   SymbolNonTerminal{
								   Name: "character2",
							   },
							   SymbolNonTerminal{
								   Name: "text2",
							   },
						   },

					   },
				   },
			   },
		   },
		   { // 13 <character1>
			   Head: SymbolNonTerminal{
				   Name: "character1",
			   },
			   Tail: Substitution{
				   Sequences: []Sequence{
					   {
						   Symbols: []Symbol{
							   SymbolTerminal{
								   Name: "'",
							   },
						   },
					   },
					   {
						   Symbols: []Symbol{
							   SymbolNonTerminal{
								   Name: "character",
							   },
						   },

					   },
				   },
			   },
		   },
		   { // 14 <character1>
			   Head: SymbolNonTerminal{
				   Name: "character2",
			   },
			   Tail: Substitution{
				   Sequences: []Sequence{
					   {
						   Symbols: []Symbol{
							   SymbolTerminal{
								   Name: "\"",
							   },
						   },
					   },
					   {
						   Symbols: []Symbol{
							   SymbolNonTerminal{
								   Name: "character",
							   },
						   },

					   },
				   },
			   },
		   },
		   { // 15 <character>
			   Head: SymbolNonTerminal{
				   Name: "character",
			   },
			   Tail: Substitution{
				   Sequences: []Sequence{
					   {
						   Symbols: []Symbol{
							   SymbolNonTerminal{
								   Name: "letter",
							   },
						   },
					   },
					   {
						   Symbols: []Symbol{
							   SymbolNonTerminal{
								   Name: "digit",
							   },
						   },
					   },
					   {
						   Symbols: []Symbol{
							   SymbolNonTerminal{
								   Name: "symbol",
							   },
						   },
					   },
						{
							Symbols: []Symbol{
								SymbolNonTerminal{
									Name: "escape-sequence",
								},
							},
						},
				   },
			   },
		   },
			{ // 16 <escape-sequence>
				Head: SymbolNonTerminal{
					Name: "escape-sequence",
				},
				Tail: Substitution{
					Sequences: []Sequence{
						{
							Symbols: []Symbol{
								SymbolTerminal{
									Name: "\\",
								},
								SymbolNonTerminal{
									Name: "escaped-char",
								},
							},
						},
					},
				},
			},
			{ // 17 <escaped-char>
				Head: SymbolNonTerminal{
					Name: "escaped-char",
				},
				Tail: Substitution{
					Sequences: []Sequence{
						{
							Symbols: []Symbol{
								SymbolTerminal{
									Name: "t", // ASCII 0x09 horizontal tab
								},
							},
						},
						{
							Symbols: []Symbol{
								SymbolTerminal{
									Name: "n", // ASCII 0x0A new line
								},
							},
						},
						{
							Symbols: []Symbol{
								SymbolTerminal{
									Name: "r", // ASCII 0x0D carriage return
								},
							},
						},
						{
							Symbols: []Symbol{
								SymbolTerminal{
									Name: "\"", // ASCII 0x22 quotation mark
								},
							},
						},
						{
							Symbols: []Symbol{
								SymbolTerminal{
									Name: "\\", // ASCII 0x5C reverse solidus
								},
							},
						},
					},
				},
			},
		   { // 18 <rule-name>
			   Head: SymbolNonTerminal{
				   Name: "rule-name",
			   },
			   Tail: Substitution{
				   Sequences: []Sequence{
					   {
						   Symbols: []Symbol{
							   SymbolNonTerminal{
								   Name: "letter",
							   },
							   SymbolNonTerminal{
								   Name: "rule-name-tail",
							   },
						   },
					   },
				   },
			   },
		   },
		   { // 19 <rule-name-tail>
			   Head: SymbolNonTerminal{
				   Name: "rule-name-tail",
			   },
			   Tail: Substitution{
				   Sequences: []Sequence{
					   {
						   Symbols: []Symbol{
							   SymbolNothing{},
						   },
					   },
					   {
						   Symbols: []Symbol{
							   SymbolNonTerminal{
								   Name: "rule-char",
							   },
							   SymbolNonTerminal{
								   Name: "rule-name-tail",
							   },
						   },
					   },
				   },
			   },
		   },
		   { // 20 <rule-char>
			   Head: SymbolNonTerminal{
				   Name: "rule-char",
			   },
			   Tail: Substitution{
				   Sequences: []Sequence{
					   {
						   Symbols: []Symbol{
							   SymbolNonTerminal{
								   Name: "letter",
							   },
						   },
					   },
					   {
						   Symbols: []Symbol{
							   SymbolNonTerminal{
								   Name: "digit",
							   },
						   },
					   },
					   {
						   Symbols: []Symbol{
							   SymbolTerminal{
								   Name: "-",
							   },
						   },
					   },
				   },
			   },
		   },
		   { // 21 <letter>
			   Head: SymbolNonTerminal{
				   Name: "letter",
			   },
			   Tail: Substitution{
				   Sequences: []Sequence{
					   {
						   Symbols: []Symbol{
							   SymbolTerminal{
								   Name: "A",
							   },
						   },
					   },
					   {
						   Symbols: []Symbol{
							   SymbolTerminal{
								   Name: "B",
							   },
						   },
					   },
					   {
						   Symbols: []Symbol{
							   SymbolTerminal{
								   Name: "C",
							   },
						   },
					   },
					   {
						   Symbols: []Symbol{
							   SymbolTerminal{
								   Name: "D",
							   },
						   },
					   },
					   {
						   Symbols: []Symbol{
							   SymbolTerminal{
								   Name: "E",
							   },
						   },
					   },
					   {
						   Symbols: []Symbol{
							   SymbolTerminal{
								   Name: "F",
							   },
						   },
					   },
					   {
						   Symbols: []Symbol{
							   SymbolTerminal{
								   Name: "G",
							   },
						   },
					   },
					   {
						   Symbols: []Symbol{
							   SymbolTerminal{
								   Name: "H",
							   },
						   },
					   },
					   {
						   Symbols: []Symbol{
							   SymbolTerminal{
								   Name: "I",
							   },
						   },
					   },
					   {
						   Symbols: []Symbol{
							   SymbolTerminal{
								   Name: "J",
							   },
						   },
					   },
					   {
						   Symbols: []Symbol{
							   SymbolTerminal{
								   Name: "K",
							   },
						   },
					   },
					   {
						   Symbols: []Symbol{
							   SymbolTerminal{
								   Name: "L",
							   },
						   },
					   },
					   {
						   Symbols: []Symbol{
							   SymbolTerminal{
								   Name: "M",
							   },
						   },
					   },
					   {
						   Symbols: []Symbol{
							   SymbolTerminal{
								   Name: "N",
							   },
						   },
					   },
					   {
						   Symbols: []Symbol{
							   SymbolTerminal{
								   Name: "O",
							   },
						   },
					   },
					   {
						   Symbols: []Symbol{
							   SymbolTerminal{
								   Name: "P",
							   },
						   },
					   },
					   {
						   Symbols: []Symbol{
							   SymbolTerminal{
								   Name: "Q",
							   },
						   },
					   },
					   {
						   Symbols: []Symbol{
							   SymbolTerminal{
								   Name: "R",
							   },
						   },
					   },
					   {
						   Symbols: []Symbol{
							   SymbolTerminal{
								   Name: "S",
							   },
						   },
					   },
					   {
						   Symbols: []Symbol{
							   SymbolTerminal{
								   Name: "T",
							   },
						   },
					   },
					   {
						   Symbols: []Symbol{
							   SymbolTerminal{
								   Name: "U",
							   },
						   },
					   },
					   {
						   Symbols: []Symbol{
							   SymbolTerminal{
								   Name: "V",
							   },
						   },
					   },
					   {
						   Symbols: []Symbol{
							   SymbolTerminal{
								   Name: "W",
							   },
						   },
					   },
					   {
						   Symbols: []Symbol{
							   SymbolTerminal{
								   Name: "X",
							   },
						   },
					   },
					   {
						   Symbols: []Symbol{
							   SymbolTerminal{
								   Name: "Y",
							   },
						   },
					   },
					   {
						   Symbols: []Symbol{
							   SymbolTerminal{
								   Name: "Z",
							   },
						   },
					   },
					   {
						   Symbols: []Symbol{
							   SymbolTerminal{
								   Name: "a",
							   },
						   },
					   },
					   {
						   Symbols: []Symbol{
							   SymbolTerminal{
								   Name: "b",
							   },
						   },
					   },
					   {
						   Symbols: []Symbol{
							   SymbolTerminal{
								   Name: "c",
							   },
						   },
					   },
					   {
						   Symbols: []Symbol{
							   SymbolTerminal{
								   Name: "d",
							   },
						   },
					   },
					   {
						   Symbols: []Symbol{
							   SymbolTerminal{
								   Name: "e",
							   },
						   },
					   },
					   {
						   Symbols: []Symbol{
							   SymbolTerminal{
								   Name: "f",
							   },
						   },
					   },
					   {
						   Symbols: []Symbol{
							   SymbolTerminal{
								   Name: "g",
							   },
						   },
					   },
					   {
						   Symbols: []Symbol{
							   SymbolTerminal{
								   Name: "h",
							   },
						   },
					   },
					   {
						   Symbols: []Symbol{
							   SymbolTerminal{
								   Name: "i",
							   },
						   },
					   },
					   {
						   Symbols: []Symbol{
							   SymbolTerminal{
								   Name: "j",
							   },
						   },
					   },
					   {
						   Symbols: []Symbol{
							   SymbolTerminal{
								   Name: "k",
							   },
						   },
					   },
					   {
						   Symbols: []Symbol{
							   SymbolTerminal{
								   Name: "l",
							   },
						   },
					   },
					   {
						   Symbols: []Symbol{
							   SymbolTerminal{
								   Name: "m",
							   },
						   },
					   },
					   {
						   Symbols: []Symbol{
							   SymbolTerminal{
								   Name: "n",
							   },
						   },
					   },
					   {
						   Symbols: []Symbol{
							   SymbolTerminal{
								   Name: "o",
							   },
						   },
					   },
					   {
						   Symbols: []Symbol{
							   SymbolTerminal{
								   Name: "p",
							   },
						   },
					   },
					   {
						   Symbols: []Symbol{
							   SymbolTerminal{
								   Name: "q",
							   },
						   },
					   },
					   {
						   Symbols: []Symbol{
							   SymbolTerminal{
								   Name: "r",
							   },
						   },
					   },
					   {
						   Symbols: []Symbol{
							   SymbolTerminal{
								   Name: "s",
							   },
						   },
					   },
					   {
						   Symbols: []Symbol{
							   SymbolTerminal{
								   Name: "t",
							   },
						   },
					   },
					   {
						   Symbols: []Symbol{
							   SymbolTerminal{
								   Name: "u",
							   },
						   },
					   },
					   {
						   Symbols: []Symbol{
							   SymbolTerminal{
								   Name: "v",
							   },
						   },
					   },
					   {
						   Symbols: []Symbol{
							   SymbolTerminal{
								   Name: "w",
							   },
						   },
					   },
					   {
						   Symbols: []Symbol{
							   SymbolTerminal{
								   Name: "x",
							   },
						   },
					   },
					   {
						   Symbols: []Symbol{
							   SymbolTerminal{
								   Name: "y",
							   },
						   },
					   },
					   {
						   Symbols: []Symbol{
							   SymbolTerminal{
								   Name: "z",
							   },
						   },
					   },
				   },
			   },
		   },
		   { // 22 <digit>
			   Head: SymbolNonTerminal{
				   Name: "digit",
			   },
			   Tail: Substitution{
				   Sequences: []Sequence{
					   {
						   Symbols: []Symbol{
							   SymbolTerminal{
								   Name: "0",
							   },
						   },
					   },
					   {
						   Symbols: []Symbol{
							   SymbolTerminal{
								   Name: "1",
							   },
						   },
					   },
					   {
						   Symbols: []Symbol{
							   SymbolTerminal{
								   Name: "2",
							   },
						   },
					   },
					   {
						   Symbols: []Symbol{
							   SymbolTerminal{
								   Name: "3",
							   },
						   },
					   },
					   {
						   Symbols: []Symbol{
							   SymbolTerminal{
								   Name: "4",
							   },
						   },
					   },
					   {
						   Symbols: []Symbol{
							   SymbolTerminal{
								   Name: "5",
							   },
						   },
					   },
					   {
						   Symbols: []Symbol{
							   SymbolTerminal{
								   Name: "6",
							   },
						   },
					   },
					   {
						   Symbols: []Symbol{
							   SymbolTerminal{
								   Name: "7",
							   },
						   },
					   },
					   {
						   Symbols: []Symbol{
							   SymbolTerminal{
								   Name: "8",
							   },
						   },
					   },
					   {
						   Symbols: []Symbol{
							   SymbolTerminal{
								   Name: "9",
							   },
						   },
					   },
				   },
			   },
		   },
		   { // 23 <symbol>
			   Head: SymbolNonTerminal{
				   Name: "symbol",
			   },
			   Tail: Substitution{
				   Sequences: []Sequence{
					   {
						   Symbols: []Symbol{
							   SymbolTerminal{
								   Name: "|",
							   },
						   },
					   },
					   {
						   Symbols: []Symbol{
							   SymbolTerminal{
								   Name: " ",
							   },
						   },
					   },
					   {
						   Symbols: []Symbol{
							   SymbolTerminal{
								   Name: "!",
							   },
						   },
					   },
					   {
						   Symbols: []Symbol{
							   SymbolTerminal{
								   Name: "#",
							   },
						   },
					   },
					   {
						   Symbols: []Symbol{
							   SymbolTerminal{
								   Name: "$",
							   },
						   },
					   },
					   {
						   Symbols: []Symbol{
							   SymbolTerminal{
								   Name: "%",
							   },
						   },
					   },
					   {
						   Symbols: []Symbol{
							   SymbolTerminal{
								   Name: "&",
							   },
						   },
					   },
					   {
						   Symbols: []Symbol{
							   SymbolTerminal{
								   Name: "(",
							   },
						   },
					   },
					   {
						   Symbols: []Symbol{
							   SymbolTerminal{
								   Name: ")",
							   },
						   },
					   },
					   {
						   Symbols: []Symbol{
							   SymbolTerminal{
								   Name: "*",
							   },
						   },
					   },
					   {
						   Symbols: []Symbol{
							   SymbolTerminal{
								   Name: "+",
							   },
						   },
					   },
					   {
						   Symbols: []Symbol{
							   SymbolTerminal{
								   Name: ",",
							   },
						   },
					   },
					   {
						   Symbols: []Symbol{
							   SymbolTerminal{
								   Name: "-",
							   },
						   },
					   },
					   {
						   Symbols: []Symbol{
							   SymbolTerminal{
								   Name: ".",
							   },
						   },
					   },
					   {
						   Symbols: []Symbol{
							   SymbolTerminal{
								   Name: "/",
							   },
						   },
					   },
					   {
						   Symbols: []Symbol{
							   SymbolTerminal{
								   Name: ":",
							   },
						   },
					   },
					   {
						   Symbols: []Symbol{
							   SymbolTerminal{
								   Name: ";",
							   },
						   },
					   },
					   {
						   Symbols: []Symbol{
							   SymbolTerminal{
								   Name: ">",
							   },
						   },
					   },
					   {
						   Symbols: []Symbol{
							   SymbolTerminal{
								   Name: "=",
							   },
						   },
					   },
					   {
						   Symbols: []Symbol{
							   SymbolTerminal{
								   Name: "<",
							   },
						   },
					   },
					   {
						   Symbols: []Symbol{
							   SymbolTerminal{
								   Name: "?",
							   },
						   },
					   },
					   {
						   Symbols: []Symbol{
							   SymbolTerminal{
								   Name: "@",
							   },
						   },
					   },
					   {
						   Symbols: []Symbol{
							   SymbolTerminal{
								   Name: "[",
							   },
						   },
					   },
					   {
						   Symbols: []Symbol{
							   SymbolTerminal{
								   Name: "]",
							   },
						   },
					   },
					   {
						   Symbols: []Symbol{
							   SymbolTerminal{
								   Name: "^",
							   },
						   },
					   },
					   {
						   Symbols: []Symbol{
							   SymbolTerminal{
								   Name: "_",
							   },
						   },
					   },
					   {
						   Symbols: []Symbol{
							   SymbolTerminal{
								   Name: "`",
							   },
						   },
					   },
					   {
						   Symbols: []Symbol{
							   SymbolTerminal{
								   Name: "{",
							   },
						   },
					   },
					   {
						   Symbols: []Symbol{
							   SymbolTerminal{
								   Name: "}",
							   },
						   },
					   },
					   {
						   Symbols: []Symbol{
							   SymbolTerminal{
								   Name: "~",
							   },
						   },
					   },

				   },
			   },
		   },
		   { // 24 <opt-whitespace>
			   Head: SymbolNonTerminal{
				   Name: "opt-whitespace",
			   },
			   Tail: Substitution{
				   Sequences: []Sequence{
					   {
						   Symbols: []Symbol{
							   SymbolTerminal{
								   Name: " ",
							   },
							   SymbolNonTerminal{
								   Name: "opt-whitespace",
							   },
						   },
					   },
					   {
						   Symbols: []Symbol{
							   SymbolNothing{},
						   },
					   },
				   },
			   },
		   },
			{ // 25 <EOL>
				Head: SymbolNonTerminal{
					Name: "EOL",
				},
				Tail: Substitution{
					Sequences: []Sequence{
						{
							Symbols: []Symbol{
								SymbolTerminal{
									Name: "\n",
								},
							},
						},
						{
							Symbols: []Symbol{
								SymbolTerminal{
									Name: "\r\n",
								},
							},
						},
					},
				},
			},
	   },
   }
}
