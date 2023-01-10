// Parser Table constuctor from BNF's AST
package tablegen

import (
	"fmt"
	"bnf"
	"parser"
)

const EOS = byte(0)

func collectRuleHeads(g bnf.Grammar) []string {
	res := make([]string, len(g.Rules))
	for i := range g.Rules {
		res[i] = g.Rules[i].Head.Name
	}
	return res
}

func enumerate(arr []string) map[string]int {
	res := map[string]int{}
	for i := range arr {
		res[arr[i]] = i
	}
	return res
}

func firstsForTerminal(t bnf.SymbolTerminal) byteSet {
	return newByteSet(t.Name[0])
}

func firstsForSymbol(firsts *map[int]byteSet,
									nonTerminalsMap map[string]int,
									symbol *bnf.Symbol) byteSet {
	symbolType := (*symbol).Type()

	if symbolType == bnf.TypeTerminal {

		t, ok := (*symbol).(bnf.SymbolTerminal)
		if !ok {
			panic("can't cast terminal to it's native type")
		}

		return firstsForTerminal(t)
	} else if symbolType == bnf.TypeNonTerminal {
		nt, ok := (*symbol).(bnf.SymbolNonTerminal)
		if !ok {
			panic("can't cast non terminal to it's native type")
		}

		ntRuleNum, ok := nonTerminalsMap[nt.Name]
		if !ok {
			panic(fmt.Sprint("no rules defined for terminal ", nt.Name))
		}

		return (*firsts)[ntRuleNum]
	} else if symbolType == bnf.TypeNothing {

		res := newByteSet(EOS)

		return res
	}

	panic(fmt.Sprint("unknown symbol type: ", symbolType))
}

func firstsForSequence(g bnf.Grammar,
                       firsts *map[int]byteSet,
                       nonTerminalsMap map[string]int,
                       reviewedRules *map[int]bool,
                       sequence *bnf.Sequence) byteSet {
	res := newByteSet()

	isEmptyFound := true

	for _, symbol := range (*sequence).Symbols {

		symbol_type := symbol.Type()

		if symbol_type == bnf.TypeTerminal {
			t, ok := symbol.(bnf.SymbolTerminal)
			if !ok {
				panic("can't cast terminal to it's native type")
			}
			tRes := firstsForTerminal(t)
			// TODO SetsCombine
			for _, k := range tRes.ToSlice() {
				res.Add(k)
			}
			return res
		} else if symbol_type == bnf.TypeNothing {
			if len((*sequence).Symbols) > 1 {
				panic("empty string not a single sequence symbol")
			}
			res.Add(EOS)
			return res
		} else if symbol_type == bnf.TypeNonTerminal {
			// firsts rules 3


			nt, ok := symbol.(bnf.SymbolNonTerminal)
			if !ok {
				panic("can't cast terminal to it's native type")
			}

			ntRuleNum, ok := nonTerminalsMap[nt.Name]
			if !ok {
				if nt.Name == "EOL" {
					res.Add('\r')
					res.Add('\n')
					return res
				} else {
					panic(fmt.Sprintf("no rules defined for non terminal <%s>", nt.Name))
				}
			}

			// Check if self refering check other substitutuion without self refering
			if _, ok := (*reviewedRules)[ntRuleNum]; ok {
				panic("left recursive expression")
				// res[EOS] = true
				// return res
			}


			reviewedRulesCopy := map[int]bool{}
			for k := range (*reviewedRules) {
				reviewedRulesCopy[k] = true
			}
			reviewedRulesCopy[ntRuleNum] = true

			firstsForRule := firstsForSubstitutuion(g, firsts,  nonTerminalsMap, &reviewedRulesCopy, &g.Rules[ntRuleNum].Tail)


			// fmt.Println("firstsForRule", firstsForRule)

			isEmptyFound = false
			for _, k := range firstsForRule.ToSlice() {
				if k == EOS {
					isEmptyFound = true
					continue
				}
				res.Add(k)
			}

			if !isEmptyFound {
				return res
			}

			continue
		}

		fmt.Println("Must not reach here on test run")
		panic("")
	}

	if isEmptyFound {
		res.Add(EOS)
	}

	return res
}


func firstsForSubstitutuion(g bnf.Grammar,
                            firsts *map[int]byteSet,
                            nonTerminalsMap map[string]int,
                            reviewedRules *map[int]bool,
                            substitution *bnf.Substitution) byteSet {
	// substitution
	res := newByteSet()
	for _, sequence := range (*substitution).Sequences {
		tRes := firstsForSequence(g, firsts, nonTerminalsMap, reviewedRules, &sequence)

		// TODO Merge Sets
		for _, k := range tRes.ToSlice() {
			res.Add(k)
		}
	}

	return res
}

func firstsForRule(g bnf.Grammar,
                   firsts *map[int]byteSet,
                   nonTerminalsMap map[string]int,
                   ruleNum int) byteSet {

	if v, ok := (*firsts)[ruleNum]; ok {
		return v
	}

	if ruleNum >= len(g.Rules) {
		panic("ruleNum exceeding grammar rules")
	}

	// Must itterate through productions and then sequence
	reviewedRules := map[int]bool{}
	reviewedRules[ruleNum] = true
	(*firsts)[ruleNum] = firstsForSubstitutuion(g, firsts,  nonTerminalsMap, &reviewedRules, &g.Rules[ruleNum].Tail)
	return (*firsts)[ruleNum]
}

func calcFirsts(g bnf.Grammar,
                nonTermNameToRuleIndex map[string]int) map[int]byteSet {

	firsts := map[int]byteSet{}

	for i := range g.Rules {
		firstsForRule(g, &firsts, nonTermNameToRuleIndex, i)
	}

	return firsts
}

func followsForRuleInSequence(g bnf.Grammar,
                              firsts map[int]byteSet,
                              follows *map[int]byteSet,
                              nonTerminalsMap map[string]int,
                              ruleNum int,
                              sequence *bnf.Sequence,
                              sequenceRuleNum int) {

	symbolsLen := len((*sequence).Symbols)
	if symbolsLen == 0 {
		panic("empty sequenceses non allowed")
	}

	{
		symbol := (*sequence).Symbols[symbolsLen - 1]
		symbol_type := symbol.Type()
		if symbol_type == bnf.TypeNonTerminal {

			nt, ok := symbol.(bnf.SymbolNonTerminal)
			if !ok {
				panic("can't cast non terminal to it's native type")
			}

			ntRuleNum, ok := nonTerminalsMap[nt.Name]
			if !ok {
				panic(fmt.Sprint("no rules defined for terminal ", nt.Name))
			}

			if ntRuleNum == ruleNum {
				if ruleNum == sequenceRuleNum {
					// fmt.Println("self repeating production - skip in follows calculation - same follow sets")
					return
				}

				if sequenceRuleNum > ruleNum {
					// valid bnf. error
					panic("can't determine follow set for subsequent rule")
				}

				// rule 3
				// TODO Merge Sets
				ruleFollows := (*follows)[sequenceRuleNum]
				for _, k := range ruleFollows.ToSlice() {
					(*follows)[ruleNum].Add(k)
				}
			}
		}
	}


	for i := symbolsLen - 2; i >= 0; i-- {
		symbol := (*sequence).Symbols[i]
		symbol_type := symbol.Type()
		if symbol_type == bnf.TypeTerminal {
			continue
		} else if symbol_type == bnf.TypeNonTerminal {
			nt, ok := symbol.(bnf.SymbolNonTerminal)
			if !ok {
				panic("can't cast non terminal to it's native type")
			}

			ntRuleNum, ok := nonTerminalsMap[nt.Name]
			if !ok {
				if nt.Name == "EOL" {
					ntRuleNum = -2
				} else {
					panic(fmt.Sprint("no rules defined for terminal ", nt.Name))
				}
			}

			if ntRuleNum != ruleNum {
				continue
			}


			// rule 2
			nextFirsts := firstsForSymbol(&firsts, nonTerminalsMap, &(*sequence).Symbols[i+1])
			for _, k := range nextFirsts.ToSlice() {

				if k == EOS {
					// rule 4
					if sequenceRuleNum > ruleNum {
						// valid bnf. error
						panic("can't determine follow set for subsequent rule")
					}

					// { FIRST(q) – Є } U FOLLOW(A)
					if ruleNum == sequenceRuleNum {
						// fmt.Println("self repeating production - skip in follows calculation - same follow sets")
						continue
					}

					// TODO Merge Sets
					ruleFollows := (*follows)[sequenceRuleNum]
					for _, j := range ruleFollows.ToSlice() {
						(*follows)[ruleNum].Add(j)
					}

					continue
				}

				(*follows)[ruleNum].Add(k)
			}
		}
	}
}

func followsForRuleinRule(g bnf.Grammar,
                          firsts map[int]byteSet,
                          follows *map[int]byteSet,
                          ruleNum int,
                          nonTerminalsMap map[string]int,
                          targetRuleNum int) {

	if ruleNum >= len(g.Rules) {
		panic("ruleNum exceeding grammar rules")
	}

	if targetRuleNum >= len(g.Rules) {
		panic("targetRuleNum exceeding grammar rules")
	}

	for _, sequence := range g.Rules[targetRuleNum].Tail.Sequences {
		followsForRuleInSequence(g, firsts, follows, nonTerminalsMap, ruleNum, &sequence, targetRuleNum)
	}

	return
}

func calcFollows(g bnf.Grammar,
                 firsts map[int]byteSet,
                 nonTermNameToRuleIndex map[string]int) map[int]byteSet {

	follows := map[int]byteSet{}
	for i := range g.Rules {
		follows[i] = newByteSet()
	}

	// Assuming starting rule is allways at index 0
	follows[0].Add(EOS)

	for ruleIndex := range g.Rules {
		// fmt.Println("follows for rule", g.Rules[ruleIndex].Head.Name)
		for againstRuleIndex := range g.Rules {
			followsForRuleinRule(g, firsts, &follows, ruleIndex, nonTermNameToRuleIndex, againstRuleIndex)
		}
	}

	return follows
}

func calcTableRow(g bnf.Grammar,
                  firsts *map[int]byteSet,
                  follows *map[int]byteSet,
                  nonTerminalsMap map[string]int,
                  ruleIndex int) map[byte][]parser.ParserOp {
	rule := g.Rules[ruleIndex]

    res := map[byte][]parser.ParserOp{}

	for _, sequence := range rule.Tail.Sequences {

		reviewedRules := map[int]bool{ruleIndex: true}
		seqFirsts := firstsForSequence(g, firsts, nonTerminalsMap, &reviewedRules, &sequence)

		for _, term := range seqFirsts.ToSlice() {
			if term == EOS {
				seqFollows, ok := (*follows)[ruleIndex]
				if !ok {
					panic("follows not calculated yet")
				}



				for _, followTerm := range seqFollows.ToSlice() {

					if v, ok := res[byte(followTerm)]; ok {
						fmt.Println(v)
						panic("multiple")
					}
					res[byte(followTerm)] = []parser.ParserOp{}
				}
				// panic("Epsilon can't go here yet")
				// For each term in ruleFollows add empty []ParserOp
				// if term in ruleFollow is Epsilon(TypeNothing) add empty []ParserOp ?to EOS?
				continue
			}

			for _, symbol := range sequence.Symbols {
				switch symbol.Type() {
				case bnf.TypeTerminal: {
					bTerm, ok := symbol.(bnf.SymbolTerminal)
					if !ok {
						panic("can't cast bnf.TypeTerminal symbol to it's native type")
					}

					res[byte(term)] = append(res[byte(term)], parser.OpTerminal(bTerm.Name))
				}
				case bnf.TypeNonTerminal: {
					bNonTerm, ok := symbol.(bnf.SymbolNonTerminal)
					if !ok {
						panic("can't cast bnf.TypeNonTerminal symbol to it's native type")
					}

					ntIndex, ok := nonTerminalsMap[bNonTerm.Name]
					if !ok {
						if bNonTerm.Name == "EOL" {
							ntIndex = parser.BuiltinEOL
						} else {
							panic("HERE")
						}
					}

					res[byte(term)] = append(res[byte(term)],
						parser.OpNonTerminal(ntIndex))
				}
				case bnf.TypeNothing: {
					if len(sequence.Symbols) > 1 {
						panic("bnf. nothing type in sequence")
					}

					panic("A?")
				}
				default:
					panic(fmt.Sprint("unknown bnf. Symbol Type: ", symbol.Type()))
				}
			}
		}
	}

	return res
}

func calcTable(g bnf.Grammar,
               firsts *map[int]byteSet,
               follows *map[int]byteSet,
               nonTerminalsMap map[string]int) map[int]map[byte][]parser.ParserOp {
	res := map[int]map[byte][]parser.ParserOp{}
	for ruleNum := range g.Rules {
		res[ruleNum] = calcTableRow(g, firsts, follows, nonTerminalsMap, ruleNum)
	}
	return res
}

func ToParserTable(g bnf.Grammar) (*map[int]map[byte][]parser.ParserOp, *map[int]string, bool) {

	// TODO: condtional checks here

	nonTerminals := collectRuleHeads(g)
	nonTerminalsMap := enumerate(nonTerminals)

	firsts := calcFirsts(g, nonTerminalsMap)
	// for i := range g.Rules {
	// 	fmt.Println("FIRSTS(<" + g.Rules[i].Head.Name + ">):", firsts[i])
	// }

	follows := calcFollows(g, firsts, nonTerminalsMap)
	// for i := range g.Rules {
	// 	fmt.Println("FOLLOWS(<" + g.Rules[i].Head.Name + ">):", follows[i])
	// }

	res := calcTable(g, &firsts, &follows, nonTerminalsMap)

	namingMap := map[int]string{}
	for k, v := range nonTerminalsMap {
		namingMap[v] = k
	}

	return &res, &namingMap, true
}
