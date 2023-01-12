module json

go 1.19

replace (
	bnf => ./../../pkg/bnf
	bnf/fromcst => ./../../pkg/bnf/fromcst
	bnf/tablegen => ./../../pkg/bnf/tablegen
	cst => ./../../pkg/cst
	parser => ./../../pkg/parser
	stringpeekmover => ./../../pkg/stringpeekmover
)

require (
	bnf v0.0.0-00010101000000-000000000000
	bnf/fromcst v0.0.0-00010101000000-000000000000
	bnf/tablegen v0.0.0-00010101000000-000000000000
	parser v0.0.0-00010101000000-000000000000
)

require cst v0.0.0-00010101000000-000000000000 // indirect
