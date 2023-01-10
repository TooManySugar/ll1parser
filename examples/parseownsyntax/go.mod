module doodling

go 1.17

replace (
	cst => ./../../pkg/cst
	bnf => ./../../pkg/bnf
	bnf/fromcst => ./../../pkg/bnf/fromcst
	bnf/tablegen => ./../../pkg/bnf/tablegen
	parserimpl => ./../../pkg/parserimpl
	stringpeekmover => ./../../pkg/stringpeekmover
)

require (
	bnf v0.0.0-00010101000000-000000000000
	bnf/fromcst v0.0.0-00010101000000-000000000000
	bnf/tablegen v0.0.0-00010101000000-000000000000
	parserimpl v0.0.0-00010101000000-000000000000
	stringpeekmover v0.0.0-00010101000000-000000000000
)

require (
	cst v0.0.0-00010101000000-000000000000 // indirect
)
