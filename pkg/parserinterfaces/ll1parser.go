package parserinterfaces

import (
	"cst"
)

type LL1ParserReader interface {
	// Get current symbol returns 0 if end reached
	Peek() byte
	// Move to next symbol
	Move()
	//
	Pos() int
}

type LL1Parser interface {
	// Returns highly abstract parse tree and map to read it's nodes names
	// In case of error returns nil, nil, err
	Parse(in LL1ParserReader) (parseTree cst.Node,
	                           namingMap *map[int]string,
	                           err error)
}
