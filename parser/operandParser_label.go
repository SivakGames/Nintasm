package parser

import (
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

func (p *LabelOperandParser) Process(operationType tokenEnum, operationValue string, operationLabel string) {
	isLocal := strings.HasPrefix(operationLabel, ".")

	switch operationType {
	case enumTokenTypes.IDENTIFIER:
		if isLocal {
			if p.parentLabel == "" {
				fmt.Println("No parent label!")
			}
		} else {
			p.parentLabel = operationLabel
		}

		org := romBuilder.GetOrg()
		identifierNode := operandFactory.CreateIdentifierNode(operationType, operationLabel)
		numberNode := operandFactory.CreateNumericLiteralNode(enumTokenTypes.NUMBER_decimal, strconv.Itoa(org), org)
		assignmentNode := operandFactory.CreateAssignmentNode(identifierNode, numberNode)
		_ = interpreter.EvaluateNode(assignmentNode)

	case enumTokenTypes.ASSIGN_EQU, enumTokenTypes.ASSIGN_simple:
		if p.lookaheadType == enumTokenTypes.None {
			fmt.Println("\x1b[31mMissing assignment operand!\x1b[0m")
			return
		}

		operandList, err := p.GetOperandList()
		if err != nil {
			fmt.Println(err)
			return
		}
		if len(operandList) == 1 {
			identifierNode := operandFactory.CreateIdentifierNode(operationType, operationLabel)
			assignmentNode := operandFactory.CreateAssignmentNode(identifierNode, operandList[0])
			_ = interpreter.EvaluateNode(assignmentNode)
		} else {
			fmt.Println("Solve multi assignment operandz")
		}

	case enumTokenTypes.DIRECTIVE_labeled:
		//fmt.Println("label dir")
	case enumTokenTypes.DIRECTIVE_labeledBlockStart:
		fmt.Println("label dir st")
	case enumTokenTypes.DIRECTIVE_labeledBlockEnd:
		fmt.Println("label dir en")
	default:
		fmt.Println("BAD LABEL OPERATION TYPE!!!")
	}
}

func addLabelToSymbolTable() {

}

func (p *LabelOperandParser) GetParentLabel() string {
	return p.parentLabel
}
