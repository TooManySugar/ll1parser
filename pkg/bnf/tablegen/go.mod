module tablegen

go 1.17

replace (
    bnf => ./../
    parserimpl => ./../pkg/parserimpl
)

require (
    bnf v0.0.0-00010101000000-000000000000
    parserimpl v0.0.0-00010101000000-000000000000
)
