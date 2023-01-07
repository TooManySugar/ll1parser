package main

import (
	"strings"
	"fmt"
	"os"
	"os/exec"
)

type testDir struct {
	name string
	tests []string
}

var tests = []testDir{
	{
		name: "bnftablegen",
		tests: []string {
			"TestLinearMultiRuleGrammarToParsingTableResTable",
			"TestOrGrammarToParsingTableResTable",
			// "TestLRecursiveGrammarToParsingTableResTable",
			"TestResolvedLRecursiveGrammarToParsingTableResTable",
			"TestBNFGrammarToParsingTableResNamingMap",
			"TestBNFGrammarToParsingTableResTable",
		},
	},
	{
		name: "cst",
		tests: []string {
			"TestFprintTreeNamed",
		},
	},
	{
		name: "parser",
		tests: []string {
			"TestCustomParserCanParse",
			"TestLinearMultiRuleParserParse",
			"TestOrParserParse",
			"TestResolvedLRecursiveParserParse",
			"TestBNFParserCanParse1",
			"TestParserParseNamingMap",
		},
	},
}


// run using
// go run executetests.go | grep -Ea 'PASS|FAIL|_test'
func main() {

	for _, td := range tests {
		err := os.Chdir(td.name)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: can't open test folder: %s", err.Error())
			os.Exit(1)
		}

		sb := strings.Builder{}
		for i, testName := range td.tests {
			if i > 0 {
				sb.WriteString("|")
			}
			sb.WriteString(testName)
		}

		cmd := exec.Command("go", "test", "-v", "-run", sb.String())
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		if err != nil {
			fmt.Println(err)
		}
		os.Chdir("..")
	}
}
