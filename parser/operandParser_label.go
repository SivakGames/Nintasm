package parser

import (
	"fmt"
	enumTokenTypes "misc/nintasm/enums/tokenTypes"
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

	case enumTokenTypes.ASSIGN_EQU, enumTokenTypes.ASSIGN_simple:
		//fmt.Println("assigning stuff")
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
