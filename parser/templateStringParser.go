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
func (p *Parser) parseTemplateString(templateString string, tokenType enumTokenTypes.Def) (string, error) {
	templateLabel := templateString[2 : len(templateString)-1]
	p.startAndAdvanceToNext(templateLabel, "templateString")
	dynamicString := ""

	for p.lookaheadType != enumTokenTypes.None {
		if p.lookaheadType == enumTokenTypes.DELIMITER_leftCurlyBrace {
			o, err := p.getTemplateOperand()
			if err != nil {
				return templateString, err
			}
			dynamicString += o
		} else {
			dynamicString += p.lookaheadValue
			p.eatFreelyAndAdvance(p.lookaheadType)
		}
	}

	if dynamicString == "" {
		return dynamicString, errorHandler.AddNew(enumErrorCodes.ParserTemplateStringEmpty)
	}

	if tokenType == enumTokenTypes.DYNAMIC_LABEL {
		identifierToValidate := dynamicString
		if strings.HasPrefix(identifierToValidate, ".") {
			identifierToValidate = identifierToValidate[1:]
		}

		if !p.tokenizer.IsTokenIdentifierLikeWithParent(identifierToValidate) {
			return templateString, errorHandler.AddNew(enumErrorCodes.ParserTemplateStringNotIdentifier, dynamicString)
		}
	}

	return dynamicString, nil
}

func (p *Parser) getTemplateOperand() (string, error) {
	var err error

	err = p.eatFreelyAndAdvance(enumTokenTypes.DELIMITER_leftCurlyBrace)
	if err != nil {
		return "", err
	}

	operandString := ""
	for p.lookaheadType != enumTokenTypes.DELIMITER_rightCurlyBrace && p.lookaheadType != enumTokenTypes.None {
		operandString += p.lookaheadValue
		err = p.eatFreelyAndAdvance(p.lookaheadType)
		if err != nil {
			return operandString, err
		}
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
