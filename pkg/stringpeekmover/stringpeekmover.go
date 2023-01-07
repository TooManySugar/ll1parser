// Implementation of PeekMover interface over string
package stringpeekmover

import (
	pi "parserinterfaces"
)

type simplePeekMover_t struct {
	str string
	pos int
}

func NewSimplePeekMover(str string) pi.LL1ParserReader {
	return &simplePeekMover_t{str: str, pos: 0}
}

func (pm *simplePeekMover_t) Peek() byte {
	if pm.pos >= len(pm.str) {
		return byte(0)
	}
	return pm.str[pm.pos]
}

func (pm *simplePeekMover_t) Move()  {
	pm.pos++
}

func (pm *simplePeekMover_t) Reset()  {
	pm.pos = 0
}
