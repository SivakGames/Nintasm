package parser

import (
	"errors"
	"fmt"
	"misc/nintasm/assemble/errorHandler"
	enumErrorCodes "misc/nintasm/constants/enums/errorCodes"
	enumTokenTypes "misc/nintasm/constants/enums/tokenTypes"
	"misc/nintasm/interpreter"
	"misc/nintasm/interpreter/environment/unresolvedTable"
	"misc/nintasm/interpreter/operandFactory"
	"misc/nintasm/romBuilder"
	"strings"
)

const ASSIGNMENT_MIN_OPERANDS = 1
const ASSIGNMENT_MAX_OPERANDS = 64
const ASSIGNMENT_MANAULLY_EVALS = false

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
			return err // ❌ Fails
		}
	}

	switch operationType {
	case enumTokenTypes.IDENTIFIER:
		if !isLocal {
			interpreter.OverwriteParentLabel(operationLabel)
		}

		labelAssignNode := operandFactory.CreateAssignLabelNode(operationLabel, romBuilder.GetOrg())
		_, err := interpreter.EvaluateNode(labelAssignNode)
		if err != nil {
			return err // ❌ Fails
		}
		return nil

	case enumTokenTypes.ASSIGN_EQU, enumTokenTypes.ASSIGN_simple:
		if p.lookaheadType == enumTokenTypes.None {
			return errorHandler.AddNew(enumErrorCodes.AssignmentMissingOperand)
		}

		operandList, err := p.GetOperandList(
			ASSIGNMENT_MIN_OPERANDS, ASSIGNMENT_MAX_OPERANDS, ASSIGNMENT_MANAULLY_EVALS,
			nil,
		)
		if err != nil {
			return err // ❌ Fails
		}

		if len(operandList) == 1 {
			identifierNode := operandFactory.CreateIdentifierNode(operationLabel)
			assignmentNode := operandFactory.CreateAssignmentNode(identifierNode, operandList[0])
			unresolvedAssignNode, err := interpreter.EvaluateNode(assignmentNode)
			if err != nil {
				//Sees if this is unresolvable only...
				err := errorHandler.CheckErrorContinuesUpwardPropagation(err, enumErrorCodes.UnresolvedIdentifier)
				if err != nil {
					//Sees if this is a fatal error
					err := errorHandler.CheckErrorContinuesUpwardPropagation(err, enumErrorCodes.Error)
					if err != nil {
						return err // ❌❌ CONTINUES Failing!
					}
				} else {
					unresolvedTable.AddUnresolvedSymbol(unresolvedAssignNode)
				}
			}
		} else {
			fmt.Println(operandList)
			return errors.New("Solve multi assignment operandz")
		}
		return nil

	default:
		panic("🛑 BAD LABEL OPERATION TYPE!!!")
	}
}
