package parser

import (
	"misc/nintasm/assemble/errorHandler"
	enumErrorCodes "misc/nintasm/constants/enums/errorCodes"
	enumTokenTypes "misc/nintasm/constants/enums/tokenTypes"
	"misc/nintasm/interpreter"
	"misc/nintasm/interpreter/environment/namespaceTable"
	"misc/nintasm/interpreter/environment/symbolAsNodeTable"
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

func (p *LabelOperandParser) Process(operationTokenEnum tokenEnum, operationLabel string) error {
	err := p.doProcess(operationTokenEnum, operationLabel)
	if err != nil {
		err := errorHandler.CheckErrorContinuesUpwardPropagation(err, enumErrorCodes.Error)
		if err != nil {
			return err // ❌❌ CONTINUES Failing!
		}
		return nil
	}
	return nil
}

func (p *LabelOperandParser) doProcess(operationTokenEnum tokenEnum, operationLabel string) error {
	if operationTokenEnum == enumTokenTypes.TEMPLATE_STRING {
		newLabel, err := p.getTemplateString(operationLabel)
		if err != nil {
			return err // ❌ Fails
		}
		operationLabel = newLabel
		operationTokenEnum = enumTokenTypes.IDENTIFIER
	}

	isLocal := strings.HasPrefix(operationLabel, ".")
	if isLocal {
		_, err := interpreter.GetParentLabel()
		if err != nil {
			return err // ❌ Fails
		}
	}

	switch operationTokenEnum {
	case enumTokenTypes.IDENTIFIER:
		if !isLocal {
			prevParent := interpreter.GetParentLabelNoError()
			interpreter.OverwriteParentLabel(operationLabel)
			if prevParent != "" {
				symbolAsNodeTable.AddIdentifierKeyToPrevLabelNextLabelTable(prevParent, operationLabel)
			}
			interpreter.ClearLocalLabel()
		} else {
			interpreter.SetLocalLabel(operationLabel)
		}

		labelAssignNode := operandFactory.CreateAssignLabelNode(operationLabel, romBuilder.GetOrg())
		_, err := interpreter.EvaluateNode(labelAssignNode)
		if err != nil {
			err := errorHandler.CheckErrorContinuesUpwardPropagation(err, enumErrorCodes.Error)
			if err != nil {
				return err // ❌❌ CONTINUES Failing!
			}
		}

		return nil

	case enumTokenTypes.ASSIGN_EQU, enumTokenTypes.ASSIGN_simple:
		if p.lookaheadType == enumTokenTypes.None {
			return errorHandler.AddNew(enumErrorCodes.AssignmentMissingOperand)
		}

		if namespaceTable.IsDefiningNamespace && !isLocal {
			return errorHandler.AddNew(enumErrorCodes.AssignmentNamespaceNotLocal)
		} else if !namespaceTable.IsDefiningNamespace && isLocal {
			return errorHandler.AddNew(enumErrorCodes.AssignmentLocalNotInNamespace)
		}

		operandList, err := p.GetOperandList(
			ASSIGNMENT_MIN_OPERANDS, ASSIGNMENT_MAX_OPERANDS, ASSIGNMENT_MANAULLY_EVALS,
			nil,
		)
		if err != nil {
			return err // ❌ Fails
		}

		if len(operandList) == 1 {
			err := doAssignment(operationLabel, &operandList[0])
			if err != nil {
				return err
			}

		} else {
			for i, o := range operandList {
				evalNode, err := interpreter.EvaluateNode(o)
				if err != nil {
					return err
				}
				operandList[i] = evalNode
			}
			multiByteNode := operandFactory.CreateMultiByteNode(operandList)
			err := doAssignment(operationLabel, &multiByteNode)
			if err != nil {
				return err
			}
		}
		return nil

	default:
		panic("🛑 BAD LABEL OPERATION TYPE!!!")
	}
}

func doAssignment(operationLabel string, operand *Node) error {
	identifierNode := operandFactory.CreateIdentifierNode(operationLabel)
	assignmentNode := operandFactory.CreateAssignmentNode(identifierNode, *operand)
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

	return nil
}
