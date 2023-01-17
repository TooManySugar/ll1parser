package tablegen

import (
	"bnf"
	"errors"
	"fmt"
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

// enumerate
// creates maping from string slice:
// arr[i] <=> arr[slice[i]] = i
func enumerate(arr []string) map[string]int {
	res := map[string]int{}
	for i := range arr {
		res[arr[i]] = i
	}
	return res
}

func terminalFirsts(t bnf.SymbolTerminal) byteSet {
	return newByteSet(t.Name[0])
}

type tableGenerator struct {
	g         bnf.Grammar
	ruleCount int
	// rule name to it's index in grammar
	ruleMap map[string]int
	firsts  []byteSet
	follows []byteSet
}

func (tg *tableGenerator,
) sequenceFirsts(sequence *bnf.Sequence) (byteSet, error) {
	res := newByteSet()
	isEmptyStrFound := true
	for _, symbol := range (*sequence).Symbols {
		switch v := symbol.(type) {
		case bnf.SymbolTerminal:
			tRes := terminalFirsts(v)
			// TODO SetsCombine
			for _, k := range tRes.ToSlice() {
				res.Add(k)
			}
			return res, nil

		case bnf.SymbolNothing:
			if len((*sequence).Symbols) > 1 {
				return res, errors.New(
					"empty string is not a single sequence symbol")
			}
			res.Add(EOS)
			return res, nil

		case bnf.SymbolNonTerminal:
			ruleIndex, ok := tg.ruleMap[v.Name]
			if !ok {
				switch v.Name {
				case "EOL":
					res.Add('\r')
					res.Add('\n')
					return res, nil

				default:
					return res,
						fmt.Errorf("no rules defined for non terminal <%s>",
							v.Name)

				}
			}

			// TODO: left recursion check

			ruleFirsts := tg.firsts[ruleIndex]
			isEmptyStrFound = false
			for _, k := range ruleFirsts.ToSlice() {
				if k == EOS {
					isEmptyStrFound = true
					continue
				}
				res.Add(k)
			}

			if !isEmptyStrFound {
				return res, nil
			}

			continue

		default:
			panic(fmt.Sprintf("symbol of unknown type %T", v))
		}
	}

	if isEmptyStrFound {
		res.Add(EOS)
	}

	return res, nil
}

func (tg *tableGenerator,
) substitutuionFirsts(substitution *bnf.Substitution) (byteSet, error) {
	res := newByteSet()
	for _, sequence := range (*substitution).Sequences {
		seqFirsts, err := tg.sequenceFirsts(&sequence)
		if err != nil {
			return res, err
		}

		for _, k := range seqFirsts.ToSlice() {
			res.Add(k)
		}
	}

	return res, nil
}

func (tg *tableGenerator) ruleFirsts(ruleIndex int,
	firstsChanged *bool) error {
	if ruleIndex >= tg.ruleCount {
		return errors.New("ruleIndex exceeding grammar rules")
	}

	tempRes, err := tg.substitutuionFirsts(&tg.g.Rules[ruleIndex].Tail)
	if err != nil {
		return err
	}

	// add byteSet tempRes to byteSet firsts[ruleIndex]
	// if tempRes contains new value set *firtsChanged to true
	for _, k := range tempRes.ToSlice() {
		if tg.firsts[ruleIndex].Contains(k) {
			continue
		}
		*firstsChanged = true
		tg.firsts[ruleIndex].Add(k)
	}
	return nil
}

func (tg *tableGenerator) findFirsts() error {
	var err error
	firstsChanged := true
	for firstsChanged {
		firstsChanged = false
		for ruleIndex := range tg.g.Rules {
			err = tg.ruleFirsts(ruleIndex, &firstsChanged)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// symbolFirsts
// Helper function to get firsts for symbol assuming
// firsts sets already calculated
func (tg *tableGenerator) symbolFirsts(symbol *bnf.Symbol) (byteSet, error) {
	var ret byteSet
	switch v := (*symbol).(type) {
	case bnf.SymbolTerminal:
		if len(v.Name) == 0 {
			return ret, fmt.Errorf("terminal with empty string check grammar")
		}
		return terminalFirsts(v), nil

	case bnf.SymbolNothing:
		return newByteSet(EOS), nil

	case bnf.SymbolNonTerminal:
		ruleIndex, ok := tg.ruleMap[v.Name]
		if !ok {
			switch v.Name {
			case "EOL":
				return newByteSet('\n', '\r'), nil
			default:
				return ret,
					fmt.Errorf("no rules defined for non terminal <%s>", v.Name)
			}
		}
		return tg.firsts[ruleIndex], nil

	default:
		return ret, fmt.Errorf("symbol of unknown type: %T", v)
	}
}

func (tg *tableGenerator,
) ruleFollowsInSequence(
	ruleIndex int,
	sequence *bnf.Sequence,
	sequenceRuleIndex int,
	followsChanged *bool) error {

	symbolsLen := len((*sequence).Symbols)
	if symbolsLen == 0 {
		return fmt.Errorf("empty sequencese")
	}

	// Last symbol - special treatment

	// skip self repeating production in follows calculation
	// same follow sets
	for ruleIndex != sequenceRuleIndex {
		symbol := (*sequence).Symbols[symbolsLen-1]

		v, ok := symbol.(bnf.SymbolNonTerminal)
		if !ok {
			break
		}
		// last symbol in the end of sequence is non-terminal

		if v.Name != tg.g.Rules[ruleIndex].Head.Name {
			break
		}
		// processed rule is in the end of this sequence:
		//   add sequence's follows to rule's follows

		// rule 3
		// add last non terminal follows to follows[ruleIndex]
		// if it's follows contains new value set *firtsChanged to true
		ruleFollows := tg.follows[sequenceRuleIndex]
		for _, k := range ruleFollows.ToSlice() {
			if tg.follows[ruleIndex].Contains(k) {
				continue
			}
			*followsChanged = true
			tg.follows[ruleIndex].Add(k)
		}
		break
	}

	// iterate over rest of symbols searching for processed rule
	for i := symbolsLen - 2; i >= 0; i-- {
		symbol := (*sequence).Symbols[i]

		v, ok := symbol.(bnf.SymbolNonTerminal)
		if !ok {
			continue
		}

		if v.Name != tg.g.Rules[ruleIndex].Head.Name {
			continue
		}

		// rule 2
		// get firsts for next symbol in sequance (can be itself)
		nextFirsts, err := tg.symbolFirsts(&(*sequence).Symbols[i+1])
		if err != nil {
			return err
		}

		// rule 4
		// if next symbol (q) 's firsts contain empty string (Є)
		// add to target rule's follows it's firsts w/o empty
		// and follow set for sequence's rule (A)
		//
		// { FIRST(q) – Є } U FOLLOW(A)
		//
		// skip self repeating production in follows calculation:
		// same follow sets

		// add sequence's rule follows to follows[ruleIndex]
		// if it's follows contains new value set *firtsChanged to
		// true
		if nextFirsts.Contains(EOS) && ruleIndex != sequenceRuleIndex {
			// TODO Merge Sets
			ruleFollows := tg.follows[sequenceRuleIndex]
			for _, j := range ruleFollows.ToSlice() {
				if tg.follows[ruleIndex].Contains(j) {
					continue
				}
				*followsChanged = true
				tg.follows[ruleIndex].Add(j)
			}
		}

		for _, k := range nextFirsts.ToSlice() {
			if k == EOS {
				continue
			}

			if tg.follows[ruleIndex].Contains(k) {
				continue
			}
			*followsChanged = true
			tg.follows[ruleIndex].Add(k)
		}
	}

	return nil
}

func (tg *tableGenerator) ruleFollowsInRule(ruleIndex int,
	againstRuleIndex int,
	followsChanged *bool) error {
	var err error
	if ruleIndex >= tg.ruleCount {
		return fmt.Errorf("ruleNum exceeding grammar rules")
	}

	if againstRuleIndex >= tg.ruleCount {
		return fmt.Errorf("targetRuleNum exceeding grammar rules")
	}

	for _, sequence := range tg.g.Rules[againstRuleIndex].Tail.Sequences {
		err = tg.ruleFollowsInSequence(ruleIndex, &sequence,
			againstRuleIndex, followsChanged)
		if err != nil {
			return err
		}
	}

	return nil
}

func (tg *tableGenerator) findFollows() error {
	// rule 1
	// Assuming starting rule is allways at index 0
	tg.follows[0].Add(EOS)
	followsChanged := true
	var err error
	for followsChanged {
		followsChanged = false
		for ruleIndex := range tg.g.Rules {
			for againstRuleIndex := range tg.g.Rules {
				err = tg.ruleFollowsInRule(ruleIndex, againstRuleIndex,
					&followsChanged)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (tg *tableGenerator,
) makeTableRow(ruleIndex int) (map[byte][]parser.ParserOp, error) {
	rule := tg.g.Rules[ruleIndex]

	res := map[byte][]parser.ParserOp{}

	for _, sequence := range rule.Tail.Sequences {
		seqFirsts, err := tg.sequenceFirsts(&sequence)
		if err != nil {
			return res, err
		}

		for _, term := range seqFirsts.ToSlice() {
			if term == EOS {
				// For each term in ruleFollows add empty []ParserOp
				// if term in ruleFollow is Epsilon(TypeNothing)
				//     add empty []ParserOp ?to EOS?
				seqFollows := tg.follows[ruleIndex]

				for _, followTerm := range seqFollows.ToSlice() {
					if v, ok := res[byte(followTerm)]; ok {
						fmt.Println(v)
						return res,
							fmt.Errorf("grammar lead to multiple parser op " +
								"sets per table cell")
					}
					res[byte(followTerm)] = []parser.ParserOp{}
				}
				continue
			}

			for _, symbol := range sequence.Symbols {
				switch v := symbol.(type) {
				case bnf.SymbolTerminal:
					res[byte(term)] = append(res[byte(term)],
						parser.OpTerminal(v.Name))

				case bnf.SymbolNonTerminal:
					ruleIndex, ok := tg.ruleMap[v.Name]
					if !ok {
						switch v.Name {
						case "EOL":
							ruleIndex = parser.BuiltinEOL

						default:
							return res,
								fmt.Errorf("no rules defined for non terminal "+
									"<%s>", v.Name)
						}
					}

					res[byte(term)] = append(res[byte(term)],
						parser.OpNonTerminal(ruleIndex))

				case bnf.SymbolNothing:
					// This must be uncreachable:
					//
					// Nothing can only be a single symbol in sequence
					// This means sequance's firsts only contains a empty string
					// so iteration would be only over one item (empty string)
					// and it guarded off at if guard before reaching this loop
					panic("uncreachable")

				default:
					panic(fmt.Sprintf("symbol of unknown type %T", v))
				}
			}
		}
	}

	return res, nil
}

func (tg *tableGenerator,
) makeTable() (table *map[int]map[byte][]parser.ParserOp,
	err error) {
	tableV := map[int]map[byte][]parser.ParserOp{}
	for ruleIndex := range tg.g.Rules {
		tableV[ruleIndex], err = tg.makeTableRow(ruleIndex)
		if err != nil {
			return nil, err
		}
	}
	return &tableV, nil
}

func (tg *tableGenerator) Run() (table *map[int]map[byte][]parser.ParserOp,
	err error) {
	err = tg.findFirsts()
	if err != nil {
		return nil, err
	}

	err = tg.findFollows()
	if err != nil {
		return nil, err
	}

	return tg.makeTable()
}

func FromGrammar(g bnf.Grammar) (table *map[int]map[byte][]parser.ParserOp,
	rowNames *map[int]string,
	err error) {

	ruleHeads := collectRuleHeads(g)

	var tablegen tableGenerator
	// newTableGenerator
	{
		ruleHeads := collectRuleHeads(g)
		ruleCount := len(ruleHeads)

		firsts := make([]byteSet, ruleCount)
		for i := 0; i < ruleCount; i++ {
			firsts[i] = newByteSet()
		}

		follows := make([]byteSet, ruleCount)
		for i := 0; i < ruleCount; i++ {
			follows[i] = newByteSet()
		}

		tablegen = tableGenerator{
			g:         g,
			ruleCount: ruleCount,
			ruleMap:   enumerate(ruleHeads),
			firsts:    firsts,
			follows:   follows,
		}
	}

	table, err = tablegen.Run()
	if err != nil {
		return nil, nil, err
	}

	rowNamesV := make(map[int]string, len(ruleHeads))
	for k, v := range tablegen.ruleMap {
		rowNamesV[v] = k
	}

	return table, &rowNamesV, nil
}
