module fromcst

go 1.17

replace (
	bnf => ./../
	cst => ./../../cst
)

require (
	bnf v0.0.0-00010101000000-000000000000
	cst v0.0.0-00010101000000-000000000000
)
