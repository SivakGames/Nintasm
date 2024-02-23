package parser

import (
	"errors"
	enumTokenTypes "misc/nintasm/constants/enums/tokenTypes"
	"misc/nintasm/interpreter"
	"misc/nintasm/interpreter/operandFactory"
	"misc/nintasm/romBuilder"
	"strconv"
	"strings"
)

type LabelOperandParser struct {
	OperandParser
}

func NewLabelOperandParser() LabelOperandParser {
	return LabelOperandParser{}
}

func (p *LabelOperandParser) Process(operationType tokenEnum, operationValue string, operationLabel string) error {
	isLocal := strings.HasPrefix(operationLabel, ".")
	if isLocal {
		_, err := interpreter.GetParentLabel()
		if err != nil {
			return err
		}
	}

	switch operationType {
	case enumTokenTypes.IDENTIFIER:
		if !isLocal {
			interpreter.OverwriteParentLabel(operationLabel)
		}

		org := romBuilder.GetOrg()
		identifierNode := operandFactory.CreateIdentifierNode(operationType, operationLabel)
		numberNode := operandFactory.CreateNumericLiteralNode(enumTokenTypes.NUMBER_decimal, strconv.Itoa(org), org)
		assignmentNode := operandFactory.CreateAssignmentNode(identifierNode, numberNode)
		_, err := interpreter.EvaluateNode(assignmentNode)
		if err != nil {
			return err // ❌ Fails
		}
		return nil

	case enumTokenTypes.ASSIGN_EQU, enumTokenTypes.ASSIGN_simple:
		if p.lookaheadType == enumTokenTypes.None {
			return errors.New("\x1b[31mMissing assignment operand!\x1b[0m")
		}

		operandList, err := p.GetOperandList(1, 64, false, nil)
		if err != nil {
			return err // ❌ Fails
		}

		if len(operandList) == 1 {
			identifierNode := operandFactory.CreateIdentifierNode(operationType, operationLabel)
			assignmentNode := operandFactory.CreateAssignmentNode(identifierNode, operandList[0])
			_, err := interpreter.EvaluateNode(assignmentNode)
			if err != nil {
				return err // ❌ Fails
			}
		} else {
			return errors.New("Solve multi assignment operandz")
		}
		return nil

	default:
		return errors.New("BAD LABEL OPERATION TYPE!!!")
	}
}
