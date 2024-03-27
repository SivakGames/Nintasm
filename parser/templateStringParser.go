package parser

import (
	"misc/nintasm/assemble/errorHandler"
	enumErrorCodes "misc/nintasm/constants/enums/errorCodes"
	enumTokenTypes "misc/nintasm/constants/enums/tokenTypes"
	"misc/nintasm/interpreter"
	"misc/nintasm/interpreter/environment"
	"strings"
)

// --------------------------------------------------

var templateStringParser = newParser()
var templateStringOperandParser = newOperandParser()

// ==================================================

// Template string parsing
func (p *Parser) parseTemplateString(templateString string) (string, error) {
	templateLabel := templateString[2 : len(templateString)-1]
	p.startAndAdvanceToNext(templateLabel, "templateString")
	identifierString := ""

	for p.lookaheadType != enumTokenTypes.None {
		if p.lookaheadType == enumTokenTypes.DELIMITER_leftCurlyBrace {
			o, err := p.getTemplateOperand()
			if err != nil {
				return templateString, err
			}
			identifierString += o
		} else {
			identifierString += p.lookaheadValue
			p.eatFreelyAndAdvance(p.lookaheadType)
		}
	}

	if identifierString == "" {
		return identifierString, errorHandler.AddNew(enumErrorCodes.ParserTemplateStringEmpty)
	}

	identifierToValidate := identifierString
	if strings.HasPrefix(identifierToValidate, ".") {
		identifierToValidate = identifierToValidate[1:]
	}

	if !p.tokenizer.IsTokenIdentifierLikeWithParent(identifierToValidate) {
		return templateString, errorHandler.AddNew(enumErrorCodes.ParserTemplateStringNotIdentifier, identifierString)
	}

	return identifierString, nil
}

func (p *Parser) getTemplateOperand() (string, error) {
	var err error

	err = p.eatFreelyAndAdvance(enumTokenTypes.DELIMITER_leftCurlyBrace)
	if err != nil {
		return "", err
	}

	operandString := ""
	for p.lookaheadType != enumTokenTypes.DELIMITER_rightCurlyBrace {
		operandString += p.lookaheadValue
		p.eatFreelyAndAdvance(p.lookaheadType)
	}
	err = p.eatFreelyAndAdvance(enumTokenTypes.DELIMITER_rightCurlyBrace)
	if err != nil {
		return operandString, err
	}

	err = templateStringOperandParser.startAndAdvanceToNext(operandString, "operand")
	if err != nil {
		return operandString, err
	}

	operandList, err := templateStringOperandParser.GetOperandList(1, 1, true, []string{})
	if err != nil {
		return operandString, err
	}

	unresolvedFlag := environment.GetUnresolvedSilentErrorFlag()
	environment.ClearUnresolvedSilentErrorFlag()
	defer environment.SetUnresolvedSilentErrorFlagTo(unresolvedFlag)

	evaluatedNode, err := interpreter.EvaluateNode(operandList[0])
	if err != nil {
		return operandString, err
	}

	return evaluatedNode.NodeValue, nil
}
