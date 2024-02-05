package parser

import (
	"errors"
	"fmt"
	enumTokenTypes "misc/nintasm/enums/tokenTypes"
	"misc/nintasm/interpreter"
	"misc/nintasm/parser/operandFactory"
	"misc/nintasm/romBuilder"
	"strconv"
	"strings"
)

type LabelOperandParser struct {
	OperandParser
	parentLabel string
}

func NewLabelOperandParser() LabelOperandParser {
	return LabelOperandParser{
		parentLabel: "",
	}
}

func (p *LabelOperandParser) Process(operationType tokenEnum, operationValue string, operationLabel string) error {
	isLocal := strings.HasPrefix(operationLabel, ".")

	switch operationType {
	case enumTokenTypes.IDENTIFIER:
		if isLocal {
			if p.parentLabel == "" {
				return errors.New("No parent label!")
			}
		} else {
			p.parentLabel = operationLabel
		}

		org := romBuilder.GetOrg()
		identifierNode := operandFactory.CreateIdentifierNode(operationType, operationLabel)
		numberNode := operandFactory.CreateNumericLiteralNode(enumTokenTypes.NUMBER_decimal, strconv.Itoa(org), org)
		assignmentNode := operandFactory.CreateAssignmentNode(identifierNode, numberNode)
		_, err := interpreter.EvaluateNode(assignmentNode)
		if err != nil {
			return err
		}
		return nil

	case enumTokenTypes.ASSIGN_EQU, enumTokenTypes.ASSIGN_simple:
		if p.lookaheadType == enumTokenTypes.None {
			return errors.New("\x1b[31mMissing assignment operand!\x1b[0m")
		}

		operandList, err := p.GetOperandList()
		if err != nil {
			return err
		}
		if len(operandList) == 1 {
			identifierNode := operandFactory.CreateIdentifierNode(operationType, operationLabel)
			assignmentNode := operandFactory.CreateAssignmentNode(identifierNode, operandList[0])
			_, err := interpreter.EvaluateNode(assignmentNode)
			if err != nil {
				return err
			}
		} else {
			return errors.New("Solve multi assignment operandz")
		}
		return nil

	case enumTokenTypes.DIRECTIVE_labeled:
		//fmt.Println("label dir")
		return nil
	case enumTokenTypes.DIRECTIVE_labeledBlockStart:
		fmt.Println("label dir st")
		return nil
	case enumTokenTypes.DIRECTIVE_labeledBlockEnd:
		fmt.Println("label dir en")
		return nil
	default:
		return errors.New("BAD LABEL OPERATION TYPE!!!")
	}
}

func (p *LabelOperandParser) GetParentLabel() string {
	return p.parentLabel
}
