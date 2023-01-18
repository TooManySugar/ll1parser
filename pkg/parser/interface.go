package parser

import (
	"github.com/TooManySugar/ll1parser/pkg/cst"
)

type LL1Parser interface {
	// Returns highly abstract parse tree and map to read it's nodes names
	// In case of error returns nil, nil, err
	Parse(src any) (parseTree cst.Node,
	                namingMap *map[int]string,
	                err error)
}
