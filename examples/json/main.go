// JSON parsing example
//
// Although you can parse JSON using this approach you should not. Parsing
// 'sentence' in such simple 'language' as JSON to it's CST is total overkill.
// Resulting CST of 6 MiB JSON Parsing end up taking around 1 GiB of RAM
// and took about 6s. Attempt to output it to .txt file created 30+ GiB file
// (I don't know precise size because it was all free space I had at the moment)
// It is proof of concept and nothing more to it.
//
// NOTE: JSON bnf not support all whitespace options described at json.org.
//       As well as any unicode symbol as a string character. It is so due
//       of how BNF Grammar made.
//
package main

import (
	"fmt"
	"os"

	"github.com/TooManySugar/ll1parser/pkg/bnf"
	"github.com/TooManySugar/ll1parser/pkg/bnf/fromcst"
	"github.com/TooManySugar/ll1parser/pkg/bnf/tablegen"
	"github.com/TooManySugar/ll1parser/pkg/cst"
	"github.com/TooManySugar/ll1parser/pkg/parser"
)

const jsonBnfFileName = "json.bnf"
const jsonDataFileName = "test.json"

func readFileText(fileName string) (string, error) {
	jBNFbuf, err := os.ReadFile(fileName)
	return string(jBNFbuf), err
}

func fatalError(a ...any) {
	fmt.Fprint(os.Stderr, "ERROR: ")
	fmt.Fprintln(os.Stderr, a...)
	os.Exit(1)
}

func main() {

	jsonBnfText, err := readFileText(jsonBnfFileName)
	if err != nil {
		fatalError("can't read JSON BNF text:", err.Error())
	}

	jsonText, err := readFileText(jsonDataFileName)
	if err != nil {
		fatalError("can't read JSON file: ", err.Error())
	}

	bnfParserTable, bnfParserTableNames, err :=
		tablegen.FromGrammar(bnf.SelfGrammar())
	if err != nil {
		fatalError("can't build parser table from BNF grammar:", err.Error())
	}

	// Parser to parse JSON's BNF syntax
	bnfParser := parser.NewLL1Parser(*bnfParserTable, *bnfParserTableNames)
	jsonGrammarCst, _, err := bnfParser.Parse(jsonBnfText)
	if err != nil {
		fatalError("can't parse JSON BNF text:", err.Error())
	}

	// Create JSON BNF Grammar in BNF AST form
	jsonGrammar, err := fromcst.SelfCSTtoASTBindings().
	                        ToAST(jsonGrammarCst, jsonBnfText)
	if err != nil {
		fatalError("can't build JSON BNF AST from JSON BNF CST:", err.Error())
	}

	jsonParserTable, jsonParserTableNames, err := tablegen.FromGrammar(
	                                                 *jsonGrammar)
	if err != nil {
		fatalError("can't build parser table from JSON grammar:", err.Error())
	}
	fmt.Println("Created table for JSON parser")

	// JSON Parser!
	jsonParser := parser.NewLL1Parser(*jsonParserTable, *jsonParserTableNames)

	// From this point it can be used to parse JSON data
	jsonCST, _, err := jsonParser.Parse(jsonText)
	if err != nil {
		fatalError("can't parse JSON text:", err.Error())
	}
	fmt.Println("Parsed successfully")

	// Print all key-values of JSON
	v := NewKeyPrintingVisitor(jsonText, os.Stdout)
	fmt.Println("JSON key structure:")
	cst.Walk(v, jsonCST)
}

type keyPrintingVisitor struct {
	level int
	f *os.File
	str *string
}

func NewKeyPrintingVisitor(str string, f *os.File) (v keyPrintingVisitor)  {
	return keyPrintingVisitor{
		level: 0,
		f: f,
		str: &str,
	}
}

func (v keyPrintingVisitor) Visit(node cst.Node) (w cst.Visitor) {
	if node == nil {
		return nil
	}

	nt := node.Type()

	if nt > 10 {
		return nil
	}

	if nt == 6 {
		s := node.Childs()[0].Childs()[1]
		for i := 0; i < v.level; i++ {
			fmt.Fprint(v.f, "  ")
		}
		fmt.Fprintf(v.f, "%s\n", (*v.str)[s.Pos():s.End()])
	}

	if nt == 9 || nt == 10 {
		v.level++
	}

	return v
}
