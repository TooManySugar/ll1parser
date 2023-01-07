package cst

import (
	"fmt"
	"io"
	"os"
	"strings"
)

type Node struct {
	Type int
	Name string
	Childs []Node
}

type typeName struct {
	Type int
	Name string
	Level int
}

func traverseLRM(node *Node, level int) []typeName {

	res := []typeName{}
	sb := strings.Builder{}
	sb.WriteString(node.Name)
	for _, c := range node.Childs {
		c_res := traverseLRM(&c, level+1)
		res = append(res, c_res...)
		sb.WriteString(c_res[0].Name)
	}

	res = append([]typeName{{
		Type: node.Type,
		Name: sb.String(),
		Level: level,
	}}, res...)

	return res
}

func FprintTreeNamed(w io.Writer, root Node, namingMap map[int]string) (int, error) {
	nodeTypeToName := func (name_id int) string {
		val, ok := namingMap[name_id]
		if !ok {
			return fmt.Sprintf("Unknown_%d", name_id)
		}
		return val
	}

	res := traverseLRM(&root, 0)

	sb := strings.Builder{}
	for _, node := range res {
		for i := 0; i < (node.Level); i++ {
			sb.WriteString("  ")
		}
		sb.WriteString(fmt.Sprintf("%s: string(%s)\n", nodeTypeToName(node.Type), node.Name))
	}

	return w.Write([]byte(sb.String()))
}

func PrintTreeNamed(root Node, namingMap map[int]string) (int, error) {
	return FprintTreeNamed(os.Stdout, root, namingMap)
}
