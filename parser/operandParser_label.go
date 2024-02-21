package parser

import (
	"errors"
	"misc/nintasm/assemble/blockStack"
	enumTokenTypes "misc/nintasm/constants/enums/tokenTypes"
	"misc/nintasm/interpreter"
	"misc/nintasm/interpreter/operandFactory"
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
			useParentLabel := p.generateParentLabel()
			if useParentLabel == "" {
				return errors.New("No parent label! Cannot use local label yet!")
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
			return err // ❌ Fails
		}
		return nil

	case enumTokenTypes.ASSIGN_EQU, enumTokenTypes.ASSIGN_simple:
		if p.lookaheadType == enumTokenTypes.None {
			return errors.New("\x1b[31mMissing assignment operand!\x1b[0m")
		}

		operandList, err := p.GetOperandList(1, 64, false)
		if err != nil {
			return err // ❌ Fails
		}

		if isLocal {
			useParentLabel := p.generateParentLabel()
			if useParentLabel == "" {
				return errors.New("No parent label! Cannot use local label for assignment!")
			}
			operationLabel = useParentLabel + operationLabel
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

func (p *LabelOperandParser) generateParentLabel() string {
	useParentLabel := blockStack.GetTemporaryOverwritingParentLabel()
	if useParentLabel == "" {
		useParentLabel = p.parentLabel
	}
	return useParentLabel
}

func (p *LabelOperandParser) GetParentLabel() string {
	return p.parentLabel
}
