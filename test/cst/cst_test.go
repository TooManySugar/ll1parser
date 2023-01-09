package cst_test

import (
	"testing"
	"cst"
	tc "testcommon"
	"bytes"
	"os"
	"io"
	"fmt"
	"strings"
)


type typeName struct {
	Type int
	Name string
	Level int
}

type namedVisitor struct {
	level int
	names map[int]string
	output *[]typeName
	visited []int
}

func NewNamedVisitor(names map[int]string) (v namedVisitor)  {
	o := make([]typeName, 0)
	return namedVisitor{
		level: 0,
		names: names,
		output: &o,
		visited: []int{},
	}
}

func (v namedVisitor) Visit(node *cst.Node) (w cst.Visitor) {
	if node == nil {
		return nil
	}

	if len(node.Name) != 0 {
		for _, oi := range v.visited {
			(*v.output)[oi].Name += node.Name
		}
	}

	*v.output = append(*v.output, typeName{
		Type: node.Type,
		Name: node.Name,
		Level: v.level,
	})
	return namedVisitor {
		level: v.level + 1,
		names: v.names,
		output: v.output,
		visited: append(v.visited, len(*v.output)-1),
	}
}

func FprintTreeNamed(w io.Writer, root cst.Node, namingMap map[int]string) (int, error) {
	nodeTypeToName := func (name_id int) string {
		val, ok := namingMap[name_id]
		if !ok {
			return fmt.Sprintf("Unknown_%d", name_id)
		}
		return val
	}

	v := NewNamedVisitor(namingMap)
	cst.Walk(v, root)


	sb := strings.Builder{}
	for _, node := range *v.output {
		for i := 0; i < (node.Level); i++ {
			sb.WriteString("  ")
		}
		sb.WriteString(fmt.Sprintf("%s: string(%s)\n", nodeTypeToName(node.Type), node.Name))
	}
	return w.Write([]byte(sb.String()))
}

func TestFprintTreeNamed(t *testing.T) {

	f, err := os.Open("refTestFprintTreeNamed.bin")
	if err != nil {
		t.Errorf("Can't open reference file: %s", err.Error())
		t.Fail()
		return
	}
	defer f.Close()

	tree := cst.Node {
		Type: 0, // syntax
		Name: "",
		Childs: []cst.Node{
			{
				Type: 1, // rule
				Name: "",
				Childs: []cst.Node{
					{
						Type: 2, // opt-whitespace
						Name: "",
						Childs: []cst.Node{
							{
								Type: 8000, // const char
								Name: " ",
								Childs: []cst.Node{},
							},
							{
								Type: 2, // opt-whitespace
								Name: "",
								Childs: []cst.Node{
									{
										Type: 8000, // const char
										Name: "",
										Childs: []cst.Node{},
									},
								},
							},
						},
					},
					{
						Type: 8000, // const char
						Name: "<",
						Childs: []cst.Node{},
					},
					{
						Type: 3, // rule-name
						Name: "",
						Childs: []cst.Node{
							{
								Type: 3, // rule-name
								Name: "",
								Childs: []cst.Node{
									{
										Type: 3, // rule-name
										Name: "",
										Childs: []cst.Node{
											{
												Type: 3, // rule-name
												Name: "",
												Childs: []cst.Node{
													{
														Type: 4, // letter
														Name: "",
														Childs: []cst.Node{
															{
																Type: 8000, // const char
																Name: "E",
																Childs: []cst.Node{},
															},
														},
													},
												},
											},
											{
												Type: 5, // rule-char
												Name: "",
												Childs: []cst.Node{
													{
														Type: 4, // letter
														Name: "",
														Childs: []cst.Node{
															{
																Type: 8000, // const char
																Name: "x",
																Childs: []cst.Node{},
															},
														},
													},
												},
											},
										},
									},
									{
										Type: 5, // rule-char
										Name: "",
										Childs: []cst.Node{
											{
												Type: 4, // letter
												Name: "",
												Childs: []cst.Node{
													{
														Type: 8000, // const char
														Name: "p",
														Childs: []cst.Node{},
													},
												},
											},
										},
									},
								},
							},
							{
								Type: 5, // rule-char
								Name: "",
								Childs: []cst.Node{
									{
										Type: 4, // letter
										Name: "",
										Childs: []cst.Node{
											{
												Type: 8000, // const char
												Name: "r",
												Childs: []cst.Node{},
											},
										},
									},
								},
							},
						},
					},
					{
						Type: 8000, // const char
						Name: ">",
						Childs: []cst.Node{},
					},
					{
						Type: 2, // opt-whitespace
						Name: "",
						Childs: []cst.Node{
							{
								Type: 8000, // const char
								Name: " ",
								Childs: []cst.Node{},
							},
							{
								Type: 2, // opt-whitespace
								Name: "",
								Childs: []cst.Node{
									{
										Type: 8000, // const char
										Name: "",
										Childs: []cst.Node{},
									},
								},
							},
						},
					},
					{
						Type: 8000, // const char
						Name: "::=",
						Childs: []cst.Node{},
					},
					{
						Type: 2, // opt-whitespace
						Name: "",
						Childs: []cst.Node{
							{
								Type: 8000, // const char
								Name: " ",
								Childs: []cst.Node{},
							},
							{
								Type: 2, // opt-whitespace
								Name: "",
								Childs: []cst.Node{
									{
										Type: 8000, // const char
										Name: "",
										Childs: []cst.Node{},
									},
								},
							},
						},
					},
					{ // expression
						Type: 6, // expression
						Name: "",
						Childs: []cst.Node{
							{ // list
								Type: 8, // list
								Name: "",
								Childs: []cst.Node{
									{ // term
										Type: 9, // term
										Name: "",
										Childs: []cst.Node{
											{ // literal
												Type: 10, // literal
												Name: "",
												Childs: []cst.Node{
													{
														Type: 8000, // const char
														Name: "\"",
														Childs: []cst.Node{},
													},
													{ // text1
														Type: 11, // text1
														Name: "",
														Childs: []cst.Node{
															{ // character1
																Type: 12, // character1
																Name: "",
																Childs: []cst.Node{
																	{ // character
																		Type: 13, // character
																		Name: "",
																		Childs: []cst.Node{
																			{ // letter
																				Type: 4, // letter
																				Name: "",
																				Childs: []cst.Node{
																					{ // const
																						Type: 8000, // char
																						Name: "A",
																						Childs: []cst.Node{},
																					},
																				},
																			},
																		},
																	},
																},
															},
															{ // text1
																Type: 11, // text1
																Name: "",
																Childs: []cst.Node{
																	{ // const char
																		Type: 8000, // const char
																		Name: "",
																		Childs: []cst.Node{},
																	},
																},
															},
														},
													},
													{
														Type: 8000, // const char
														Name: "\"",
														Childs: []cst.Node{},
													},
												},
											},
										},
									},
								},
							},

						},
					},
					{ // line-end
						Type: 7, // line-end
						Name: "",
						Childs: []cst.Node{
							{
								Type: 2, // opt-whitespace
								Name: "",
								Childs: []cst.Node{
									{
										Type: 8000, // const char
										Name: "",
										Childs: []cst.Node{},
									},
								},
							},
							{ // EOL
								Type: 8001,
								Name: "",
								Childs: []cst.Node{},
							},
						},
					},
				},
			},
		},
	}

	typeNames := map[int]string {
		0: "syntax",
		1: "rule",
		2: "opt-whitespace",
		3: "rule-name",
		4: "letter",
		5: "rule-char",
		6: "expression",
		7: "line-end",
		8: "line-end-repeat",
		9: "list",
		10: "term",
		11: "literal",
		12: "text1",
		13: "character1",
		14: "character",

		8000: "_literal",
		8001: "EOL",
		-3:   "_nothing",
	}

	sb := bytes.Buffer{}
	FprintTreeNamed(&sb, tree, typeNames)

	tc.ReaderContentMustBeEqual(t, f, &sb)
}
