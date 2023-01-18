package bnftablegen

import (
	"testing"
	"bytes"
	"fmt"
	"os"

	"github.com/TooManySugar/ll1parser/pkg/bnf"
	"github.com/TooManySugar/ll1parser/pkg/bnf/tablegen"
	tc "github.com/TooManySugar/ll1parser/test/testcommon"
	tg "github.com/TooManySugar/ll1parser/test/testgrammars"
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

	table, _, err := tablegen.FromGrammar(grammar)
	if err != nil {
		t.Fatalf("Failed to parse grammar: %s", err.Error())
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

	table, _, err := tablegen.FromGrammar(grammar)
	if err != nil {
		t.Fatalf("Failed to parse grammar: %s", err.Error())
	}

	sb := bytes.Buffer{}
	fmt.Fprintf(&sb, "%#v", *table)

	tc.ReaderContentMustBeEqual(t, f, &sb)
}

func TestLRecursiveGrammarToParsingTableResErr(t *testing.T) {
	grammar := tg.LRecursive()

	_, _, err := tablegen.FromGrammar(grammar)
	if err == nil {
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

	table, _, err := tablegen.FromGrammar(grammar)
	if err != nil {
		t.Fatalf("Failed to parse grammar: %s", err.Error())
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
		16: "escape-sequence",
		17: "escaped-char",
		18: "rule-name",
		19: "rule-name-tail",
		20: "rule-char",
		21: "letter",
		22: "digit",
		23: "symbol",
		24: "opt-whitespace",
		25: "EOL",
	}

	grammar := bnf.SelfGrammar()

	_, res, err := tablegen.FromGrammar(grammar)
	if err != nil {
		t.Fatalf("Failed to parse grammar: %s", err.Error())
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

	table, _, err := tablegen.FromGrammar(grammar)
	if err != nil {
		t.Fatalf("Failed to parse grammar: %s", err.Error())
	}

	sb := bytes.Buffer{}
	fmt.Fprintf(&sb, "%#v", *table)

	tc.ReaderContentMustBeEqual(t, f, &sb)
}
