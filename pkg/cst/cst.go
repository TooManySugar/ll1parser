package cst

type Node struct {
	Type int
	Name string
	Childs []Node
}

type Visitor interface {
	Visit(node *Node) (w Visitor)
}

func Walk(v Visitor, node Node) {
	w := v.Visit(&node)
	if w == nil {
		return
	}
	for _, child := range node.Childs {
		Walk(w, child)
	}
}
