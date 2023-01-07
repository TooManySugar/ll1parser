module parserimpl

go 1.17

replace (
	bnf => ./../bnf
	cst => ./../cst
	parserinterfaces => ./../parserinterfaces
)

require (
	cst v0.0.0-00010101000000-000000000000
	parserinterfaces v0.0.0-00010101000000-000000000000
)

require bnf v0.0.0-00010101000000-000000000000 // indirect
