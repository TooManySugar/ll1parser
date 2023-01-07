package bnftablegen

import (
	"testing"
	"bnf"
	"bnf/tablegen"
	tc "testcommon"
	tg "testgrammars"
	"bytes"
	"fmt"
	"os"
)

func TestLinearMultiRuleGrammarToParsingTableResTable(t *testing.T) {
	f, err := os.Open("refTestLinearMultiRuleGrammarToParsingTableResTable.bin")
	if err != nil {
		t.Errorf("Can't open reference file: %s", err.Error())
		t.Fail()
		return
	}
	defer f.Close()


	grammar := tg.LinearMultiRule()

	table, _, ok := tablegen.ToParserTable(grammar)
	if !ok {
		t.Errorf("Failed to parse grammar")
		t.Fail()
		return
	}

	sb := bytes.Buffer{}
	fmt.Fprintf(&sb, "%#v", *table)

	tc.ReaderContentMustBeEqual(t, f, &sb)
}

func TestOrGrammarToParsingTableResTable(t *testing.T) {

	f, err := os.Open("refTestOrGrammarToParsingTableResTable.bin")
	if err != nil {
		t.Errorf("Can't open reference file: %s", err.Error())
		t.Fail()
		return
	}
	defer f.Close()

	grammar := tg.Or()

	table, _, ok := tablegen.ToParserTable(grammar)
	if !ok {
		t.Errorf("Failed to parse grammar")
		t.Fail()
		return
	}

	sb := bytes.Buffer{}
	fmt.Fprintf(&sb, "%#v", *table)

	tc.ReaderContentMustBeEqual(t, f, &sb)
}

func TestLRecursiveGrammarToParsingTableResErr(t *testing.T) {
	grammar := tg.LRecursive()

	_, _, ok := tablegen.ToParserTable(grammar)
	if ok {
		fmt.Println("this must fail")
		t.Fail()
		return
	}
}

func TestResolvedLRecursiveGrammarToParsingTableResTable(t *testing.T) {

	f, err := os.Open("refTestResolvedLRecursiveGrammarToParsingTableResTable.bin")
	if err != nil {
		t.Errorf("Can't open reference file: %s", err.Error())
		t.Fail()
		return
	}
	defer f.Close()

	grammar := tg.ResolvedLRecursive()

	table, _, ok := tablegen.ToParserTable(grammar)
	if !ok {
		t.Errorf("Failed to parse grammar")
		t.Fail()
		return
	}

	sb := bytes.Buffer{}
	fmt.Fprintf(&sb, "%#v", *table)

	tc.ReaderContentMustBeEqual(t, f, &sb)
}

func TestBNFGrammarToParsingTableResNamingMap(t *testing.T) {
	ref := map[int]string{
		 0: "syntax",
		 1: "more-lines",
		 2: "line",
		 3: "opt-content",
		 4: "content",
		 5: "expression",
		 6: "expression-tail",
		 7: "list",
		 8: "list-tail",
		 9: "term",
		10: "literal",
		11: "text1",
		12: "text2",
		13: "character1",
		14: "character2",
		15: "character",
		16: "rule-name",
		17: "rule-name-tail",
		18: "rule-char",
		19: "letter",
		20: "digit",
		21: "symbol",
		22: "opt-whitespace",
	}

	grammar := bnf.SelfGrammar()

	_, res, ok := tablegen.ToParserTable(grammar)
	if !ok {
		t.Errorf("Failed to parse grammar")
		t.Fail()
		return
	}

	tc.MapsIntStringMustBeEqual(t, ref, *res)
}

func TestBNFGrammarToParsingTableResTable(t *testing.T) {

	f, err := os.Open("refTestBNFGrammarToParsingTableResTable.bin")
	if err != nil {
		t.Errorf("Can't open reference file: %s", err.Error())
		t.Fail()
		return
	}
	defer f.Close()

	grammar := bnf.SelfGrammar()

	table, _, ok := tablegen.ToParserTable(grammar)
	if !ok {
		t.Errorf("Failed to parse grammar")
		t.Fail()
		return
	}

	sb := bytes.Buffer{}
	fmt.Fprintf(&sb, "%#v", *table)

	tc.ReaderContentMustBeEqual(t, f, &sb)
}
