package bnf_test

import (
	"testing"
	"fmt"
	"os"
	"bytes"
	"runtime"
	"strings"

	"cst"
	"bnf"
	"bnf/tablegen"
	"parserinterfaces"
	"parserimpl"
	spm "stringpeekmover"
	tc "testcommon"
	tg "testgrammars"
)

func TestCustomParserCanParse(t *testing.T) {

	// <A> ::= '"' <T> "^" <EOL>
	// <T> ::= "B" | "C"


	//     "      B   C   ^
	// <A> "T^EOL
	// <T>        B   C

	table :=  map[int]map[byte][]parserimpl.ParserOp {
		0: {
			'"': []parserimpl.ParserOp{
				parserimpl.OpTerminal("\""),
				parserimpl.OpNonTerminal(1),
				parserimpl.OpTerminal("^"),
				parserimpl.OpNonTerminal(parserimpl.BuiltinEOL),
			},
		},
		1: {
			'B': []parserimpl.ParserOp{
				parserimpl.OpTerminal("B"),
			},
			'C': []parserimpl.ParserOp{
				parserimpl.OpTerminal("C"),
			},
		},
	}

	// map[int]string{0: "A", 1: "T"

	parser := parserimpl.NewLL1Parser(table, map[int]string{})

	pm := spm.NewSimplePeekMover("\"B^\n")

	_, _, err := parser.Parse(pm)
	if err != nil {
		t.Errorf("Failed to parse input: %s", err.Error())
		t.Fail()
		return
	}
}

type ParserTestCaseOutput struct {
	parseTreeFile string
	namingMapFile string
	err error
}

type ParserTestCase struct {
	input string
	output ParserTestCaseOutput
}

func parserParseString(parser parserinterfaces.LL1Parser,
                       in string) (parseTree *cst.Node,
                                   namingMap *map[int]string,
                                   err error) {
	pm := spm.NewSimplePeekMover(in)
	return parser.Parse(pm)
}

func performTestsOnParser(t *testing.T,
                          parser parserinterfaces.LL1Parser,
                          tests []ParserTestCase) {
	for i, test := range tests {

		var reft, refnt *os.File

		if test.output.err == nil {
			var err error
			reft, err = os.Open(test.output.parseTreeFile)
			if err != nil {
				t.Errorf("Can't open reference file: %s", err.Error())
				t.Fail()
				return
			}
			defer reft.Close()

			refnt, err = os.Open(test.output.namingMapFile)
			if err != nil {
				t.Errorf("Can't open reference file: %s", err.Error())
				t.Fail()
				return
			}
			defer refnt.Close()
		}



		rest, resnt, reserr := parserParseString(parser, test.input)
		errorMismatchError := func() {
			t.Errorf("TEST %d: Expected error %v got error %v", i, test.output.err, reserr)
		}
		if reserr != test.output.err {
			if reserr == nil || test.output.err == nil {
				errorMismatchError()
				t.Fail()
				return
			}

			if reserr.Error() != test.output.err.Error() {
				errorMismatchError()
				t.Fail()
				return
			}
		}
		if test.output.err != nil {
			continue
		}

		sb := bytes.Buffer{}

		fmt.Fprintf(&sb, "%#v", *rest)
		tc.ReaderContentMustBeEqual(t, reft, &sb)

		sb.Reset()

		fmt.Fprintf(&sb, "%#v", *resnt)
		tc.ReaderContentMustBeEqual(t, refnt, &sb)
	}
}

func functionName() string {
	counter, _, _, success := runtime.Caller(1)
	if !success {
		panic("functionName: runtime.Caller: failed")
	}
	return strings.Split(runtime.FuncForPC(counter).Name(), ".")[1]
}

func TestLinearMultiRuleParserParse(t *testing.T) {

	grammar := tg.LinearMultiRule()

	table, tableNames, ok := tablegen.ToParserTable(grammar)
	if !ok {
		t.Errorf("Failed to parse grammar")
		t.Fail()
		return
	}

	tests := []ParserTestCase{
		{
			input: "A::=C",
			output: ParserTestCaseOutput{
				parseTreeFile: "refTetsLinearMultiRuleParserParse10.bin",
				namingMapFile: "refTetsLinearMultiRuleParserParse11.bin",
				err: nil,
			},
		},
		{
			input: "!::=C",
			output: ParserTestCaseOutput{
				err: fmt.Errorf(
					"no rules for '!' (33) and non terminal op <Expr>"),
			},
		},
		{
			input: "A:!=C",
			output: ParserTestCaseOutput{
				err: fmt.Errorf(
					"expected char ':' (58), got '!' (33)"),
			},
		},
		{
			input: "A::=!",
			output: ParserTestCaseOutput{
				err: fmt.Errorf(
					"no rules for '!' (33) and non terminal op <T>"),
			},
		},
		{
			input: "A::=C!",
			output: ParserTestCaseOutput{
				err: fmt.Errorf("expected end of input got '!' (33)"),
			},
		},
	}

	parser := parserimpl.NewLL1Parser(*table, *tableNames)

	performTestsOnParser(t, parser, tests)
}

func TestOrParserParse(t *testing.T) {

	grammar := tg.Or()

	table, tableNames, ok := tablegen.ToParserTable(grammar)
	if !ok {
		t.Errorf("Failed to parse grammar")
		t.Fail()
		return
	}

	tests := []ParserTestCase{
		{
			input: "\"B^",
			output: ParserTestCaseOutput{
				parseTreeFile: "refTestOrParserParse10.bin",
				namingMapFile: "refTestOrParserParse11.bin",
				err: nil,
			},
		},
		{
			input: "\"^",
			output: ParserTestCaseOutput{
				parseTreeFile: "refTestOrParserParse20.bin",
				namingMapFile: "refTestOrParserParse21.bin",
				err: nil,
			},
		},
	}

	parser := parserimpl.NewLL1Parser(*table, *tableNames)

	performTestsOnParser(t, parser, tests)
}

func TestResolvedLRecursiveParserParse(t *testing.T) {

	grammar := tg.ResolvedLRecursive()

	table, tableNames, ok := tablegen.ToParserTable(grammar)
	if !ok {
		t.Errorf("Failed to parse grammar")
		t.Fail()
		return
	}

	tests := []ParserTestCase{
		{
			input: "_B",
			output: ParserTestCaseOutput{
				parseTreeFile: "refTestResolvedLRecursiveParserParse10.bin",
				namingMapFile: "refTestResolvedLRecursiveParserParse11.bin",
				err: nil,
			},
		},
		{
			input: "_",
			output: ParserTestCaseOutput{
				parseTreeFile: "refTestResolvedLRecursiveParserParse20.bin",
				namingMapFile: "refTestResolvedLRecursiveParserParse21.bin",
				err: nil,
			},
		},
		{
			input: "_BBBBBBB",
			output: ParserTestCaseOutput{
				parseTreeFile: "refTestResolvedLRecursiveParserParse30.bin",
				namingMapFile: "refTestResolvedLRecursiveParserParse31.bin",
				err: nil,
			},
		},
	}

	parser := parserimpl.NewLL1Parser(*table, *tableNames)

	// for _, test := range tests {
	// 	t, nt, err := parserParseString(parser, test.input)
	// 	if err != nil {
	// 		panic(err)
	// 	}

	// 	f, _ := os.Create(test.output.parseTreeFile)
	// 	fmt.Fprintf(f, "%#v", *t)
	// 	f.Close()
	// 	f, _ = os.Create(test.output.namingMapFile)
	// 	fmt.Fprintf(f, "%#v", *nt)
	// 	f.Close()
	// }

	performTestsOnParser(t, parser, tests)
}

func TestBNFParserCanParse1(t *testing.T) {

	grammar := bnf.SelfGrammar()

	table, tableNames, ok := tablegen.ToParserTable(grammar)
	if !ok {
		t.Errorf("Failed to parse grammar")
		t.Fail()
		return
	}

	parser := parserimpl.NewLL1Parser(*table, *tableNames)

	pm1 := spm.NewSimplePeekMover(
		" <A-2-B> ::= \"A000123\" \"B\" | \"B\"| \"A\" |\"A\"|\"A\"   \n" +
		" <A2-B-2-B>::= <B> ",)

	_, _, err := parser.Parse(pm1)
	if err != nil {
		t.Errorf("Failed to parse input: %s", err.Error())
		t.Fail()
		return
	}
}

func TestParserParseNamingMap(t *testing.T) {

	grammar := bnf.SelfGrammar()

	table, tableNames, ok := tablegen.ToParserTable(grammar)
	if !ok {
		t.Errorf("Failed to parse grammar")
		t.Fail()
		return
	}

	ref := make(map[int]string, len(*tableNames) + 3)

	ref[-1] = "_literal"
	ref[-2] = "_endofline"
	ref[-3] = "_nothing"

	for k, v := range *tableNames {
		if k < 0 {
			t.Errorf("grammar.ToParserTable returned names with negative keys")
			t.Fail()
			return
		}
		ref[k] = v
	}

	parser := parserimpl.NewLL1Parser(*table, *tableNames)

	pm1 := spm.NewSimplePeekMover(
		" <A-2-B> ::= \"A000123\" \"B\" | \"B\"| \"A\" |\"A\"|\"A\"   \n" +
		" <A2-B-2-B>::= <B> ",)

	_, res, err := parser.Parse(pm1)
	if err != nil {
		t.Errorf("Failed to parse input")
		t.Fail()
		return
	}

	tc.MapsIntStringMustBeEqual(t, ref, *res)
}
