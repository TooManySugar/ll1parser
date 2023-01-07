module cst_test

go 1.17

replace (
	cst => ./../../pkg/cst
	testcommon => ./../testcommon
)

require (
	cst v0.0.0-00010101000000-000000000000
	testcommon v0.0.0-00010101000000-000000000000
)
