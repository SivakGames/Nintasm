package parser

import (
	"fmt"
	enumTokenTypes "misc/nintasm/constants/enums/tokenTypes"
)

// --------------------------------------------------

var templateStringParser = newParser()
var templateStringOperandParser = newOperandParser()

// ==================================================

// Template string parsing
func (p *Parser) parseTemplateString(templateString string) error {
	templateLabel := templateString[2 : len(templateString)-1]
	p.startAndAdvanceToNext(templateLabel, "templateString")
	identifierString := ""

	for p.lookaheadType != enumTokenTypes.None {
		if p.lookaheadType == enumTokenTypes.DELIMITER_leftCurlyBrace {
			o, err := p.getTemplateOperand()
			if err != nil {
				return err
			}
			identifierString += o
		} else {
			identifierString += p.lookaheadValue
			p.eatFreelyAndAdvance(p.lookaheadType)
		}
	}
	fmt.Println(identifierString)
	return nil
}

func (p *Parser) getTemplateOperand() (string, error) {
	p.eatFreelyAndAdvance(enumTokenTypes.DELIMITER_leftCurlyBrace)
	operandString := ""
	for p.lookaheadType != enumTokenTypes.DELIMITER_rightCurlyBrace {
		operandString += p.lookaheadValue
		p.eatFreelyAndAdvance(p.lookaheadType)
	}
	p.eatFreelyAndAdvance(enumTokenTypes.DELIMITER_rightCurlyBrace)

	templateStringOperandParser.startAndAdvanceToNext(operandString, "operand")
	node, err := templateStringOperandParser.GetOperandList(1, 1, false, []string{})
	if err != nil {
		return "", err
	}

	fmt.Println(node)

	return node[0].NodeValue, nil
}
