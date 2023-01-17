module tablegen

go 1.17

replace (
    bnf => ./../
    parser => ./../../parser
)

require (
    bnf v0.0.0-00010101000000-000000000000
    parser v0.0.0-00010101000000-000000000000
)
