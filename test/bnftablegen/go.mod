module bnfasttotable_test

go 1.17

replace (
	bnf => ./../../pkg/bnf
	bnf/tablegen => ./../../pkg/bnf/tablegen
	cst => ./../../pkg/cst
	parser => ./../../pkg/parser
	testcommon => ./../testcommon
	testgrammars => /../testgrammars
)

require (
	bnf v0.0.0-00010101000000-000000000000
	bnf/tablegen v0.0.0-00010101000000-000000000000
	testcommon v0.0.0-00010101000000-000000000000
	testgrammars v0.0.0-00010101000000-000000000000
)

require (
	cst v0.0.0-00010101000000-000000000000 // indirect
	parser v0.0.0-00010101000000-000000000000 // indirect
)
