module parser_test

go 1.17

replace (
	bnf => ./../../pkg/bnf
	bnf/tablegen => ./../../pkg/bnf/tablegen
	cst => ./../../pkg/cst
	parserimpl => ./../../pkg/parserimpl
	parserinterfaces => ./../../pkg/parserinterfaces
	stringpeekmover => ./../../pkg/stringpeekmover
	testcommon => ./../testcommon
	testgrammars => ./../testgrammars
)

require (
	bnf v0.0.0-00010101000000-000000000000
	bnf/tablegen v0.0.0-00010101000000-000000000000
	cst v0.0.0-00010101000000-000000000000
	parserimpl v0.0.0-00010101000000-000000000000
	stringpeekmover v0.0.0-00010101000000-000000000000
	testcommon v0.0.0-00010101000000-000000000000
	testgrammars v0.0.0-00010101000000-000000000000
)

require parserinterfaces v0.0.0-00010101000000-000000000000
