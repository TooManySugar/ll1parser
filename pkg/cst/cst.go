package cst

type Node interface {
	Type() int
	Pos() int
	End() int
	Childs() []Node
}

type node struct {
	typei int
	start int
	end int
	childs []Node
}

func NewNode(nType int, start int, end int, childs []Node) Node {
	return node {
		typei: nType,
		start: start,
		end: end,
		childs: childs,
	}
}

func (n node) Type() int {
	return n.typei
}

func (n node) Pos() int {
	return n.start
}

func (n node) End() int {
	return n.end
}

func (n node) Childs() []Node {
	return n.childs
}

type Visitor interface {
	Visit(node Node) (w Visitor)
}

func Walk(v Visitor, node Node) {
	w := v.Visit(node)
	if w == nil {
		return
	}
	for _, child := range node.Childs() {
		Walk(w, child)
	}
}
