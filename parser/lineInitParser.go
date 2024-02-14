package parser

import (
	enumTokenTypes "misc/nintasm/enums/tokenTypes"
	"strings"
)

const INITIAL_LINE_TARGET_TOKENIZER = "initial"

// Parser for getting/formatting content on new lines
type InitialLineParser struct {
	Parser
}

// Create helper
func NewInitialLineParser() InitialLineParser {
	return InitialLineParser{}
}

func (p *InitialLineParser) Process(line string) (string, error) {

	err := p.startAndAdvanceToNext(line, INITIAL_LINE_TARGET_TOKENIZER)
	if err != nil {
		return "", err
	}

	reformattedString := ""

	for p.hasMore {
		if p.lookaheadType == enumTokenTypes.WHITESPACE {
			reformattedString += " "
		} else {
			reformattedString += p.lookaheadValue
		}
		err = p.eatFreelyAndAdvance(p.lookaheadType)
		if err != nil {
			return "", err
		}
	}
	reformattedString = strings.TrimRight(reformattedString, " ")
	return reformattedString, nil
}
