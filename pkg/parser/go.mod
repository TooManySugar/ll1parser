module parser

go 1.19

replace (
	bnf => ./../bnf
	cst => ./../cst
)

require (
	cst v0.0.0-00010101000000-000000000000
)

require bnf v0.0.0-00010101000000-000000000000 // indirect
