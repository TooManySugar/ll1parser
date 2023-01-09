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

type namedVisitor struct {
	level int
	names map[int]string
	sb *strings.Builder
	str *string
}

func NewNamedVisitor(names map[int]string, str string) (v namedVisitor)  {
	sb := strings.Builder{}
	return namedVisitor{
		level: 0,
		names: names,
		sb: &sb,
		str: &str,
	}
}

func (v namedVisitor) Visit(node cst.Node) (w cst.Visitor) {
	if node == nil {
		return nil
	}

	nodeTypeToName := func (name_id int) string {
		val, ok := v.names[name_id]
		if !ok {
			return fmt.Sprintf("Unknown_%d", name_id)
		}
		return val
	}
	fmt.Fprintf(v.sb, "%s%s: string(%s)\n", strings.Repeat("  ", int(v.level)), nodeTypeToName(node.Type()), (*v.str)[node.Pos():node.End()])

	return namedVisitor {
		level: v.level + 1,
		names: v.names,
		sb: v.sb,
		str: v.str,
	}
}

func FprintTreeNamed(w io.Writer, root cst.Node, namingMap map[int]string, str string) (int, error) {
	v := NewNamedVisitor(namingMap, str)
	cst.Walk(v, root)
	return w.Write([]byte(v.sb.String()))
}

func TestFprintTreeNamed(t *testing.T) {

	f, err := os.Open("refTestFprintTreeNamed.bin")
	if err != nil {
		t.Errorf("Can't open reference file: %s", err.Error())
		t.Fail()
		return
	}
	defer f.Close()

	str := " <Expr> ::= \"A\""

	tree := cst.NewNode(0, // syntax
		0, 15,
		// Name: "",
		[]cst.Node{
			cst.NewNode(1, // rule
				0, 15,
				[]cst.Node{
					cst.NewNode(2, // opt-whitespace
						0, 1,
						[]cst.Node{
							cst.NewNode(8000, // const char
								0, 1,
								nil,
							),
							cst.NewNode(2, // opt-whitespace
								1, 1,
								[]cst.Node{
									cst.NewNode(8000, // const char
										1, 1,
										nil,
									),
								},
							),
						},
					),
					cst.NewNode(8000, // const char
						1, 2,
						nil,
					),
					cst.NewNode(3, // rule-name
						2, 6,
						[]cst.Node{
							cst.NewNode(3, // rule-name
								2, 5,
								[]cst.Node{
									cst.NewNode(3, // rule-name
										2, 4,
										[]cst.Node{
											cst.NewNode(3, // rule-name
												2, 3,
												[]cst.Node{
													cst.NewNode(4, // letter
														2, 3,
														[]cst.Node{
															cst.NewNode(8000, // const char
																2, 3,
																nil,
															),
														},
													),
												},
											),
											cst.NewNode(5, // rule-char
												3, 4,
												[]cst.Node{
													cst.NewNode(4, // letter
														3, 4,
														[]cst.Node{
															cst.NewNode(8000, // const char
																3, 4,
																nil,
															),
														},
													),
												},
											),
										},
									),
									cst.NewNode(5, // rule-char
										4, 5,
										[]cst.Node{
											cst.NewNode(4, // letter
												4, 5,
												[]cst.Node{
													cst.NewNode(8000, // const char
														4, 5,
														nil,
													),
												},
											),
										},
									),
								},
							),
							cst.NewNode(5, // rule-char
								5, 6,
								[]cst.Node{
									cst.NewNode(4, // letter
										5, 6,
										[]cst.Node{
											cst.NewNode(8000, // const char
												5, 6,
												nil,
											),
										},
									),
								},
							),
						},
					),
					cst.NewNode(8000, // const char
						6, 7,
						nil,
					),
					cst.NewNode(2, // opt-whitespace
						7, 8,
						[]cst.Node{
							cst.NewNode(8000, // const char]
								7, 8,
								nil,
							),
							cst.NewNode(2, // opt-whitespace
								8, 8,
								[]cst.Node{
									cst.NewNode(8000, // const char
										8, 8,
										nil,
									),
								},
							),
						},
					),
					cst.NewNode(8000, // const char
						8, 11,
						nil,
					),
					cst.NewNode(2, // opt-whitespace
						11, 12,
						[]cst.Node{
							cst.NewNode(8000, // const char
								11, 12,
								nil,
							),
							cst.NewNode(2, // opt-whitespace
								12, 12,
								[]cst.Node{
									cst.NewNode(8000, // const char
										12, 12,
										nil,
									),
								},
							),
						},
					),
					cst.NewNode(6, // expression
						12, 15,
						[]cst.Node{
							cst.NewNode(8, // list
								12, 15,
								[]cst.Node{
									cst.NewNode(9, // term
										12, 15,
										[]cst.Node{
											cst.NewNode(10, // literal
												12, 15,
												[]cst.Node{
													cst.NewNode(8000, // const char
														12, 13,
														nil,
													),
													cst.NewNode( // text1
														11, // text1
														13, 14,
														[]cst.Node{
															cst.NewNode( // character1
																12, // character1
																13, 14,
																[]cst.Node{
																	cst.NewNode(13, // character
																		13, 14,
																		[]cst.Node{
																			cst.NewNode(4, // letter
																				13, 14,
																				[]cst.Node{
																					cst.NewNode(8000, // char
																						13, 14,
																						nil,
																					),
																				},
																			),
																		},
																	),
																},
															),
															cst.NewNode( // text1
																11, // text1
																14, 14,
																[]cst.Node{
																	cst.NewNode( // const char
																		8000, // const char
																		14, 14,
																		nil,
																	),
																},
															),
														},
													),
													cst.NewNode(8000, // const char
														14, 15,
														nil,
													),
												},
											),
										},
									),
								},
							),
						},
					),
					cst.NewNode(7, // line-end
						15, 15,
						[]cst.Node{
							cst.NewNode(2, // opt-whitespace
								15, 15,
								[]cst.Node{
									cst.NewNode(8000, // const char
										15, 15,
										nil,
									),
								},
							),
							cst.NewNode(8001, 15,
								15,
								nil,
							),
						},
					),
				},
			),
		},
	)

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
	FprintTreeNamed(&sb, tree, typeNames, str)

	tc.ReaderContentMustBeEqual(t, f, &sb)
}
