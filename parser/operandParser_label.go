package parser

import (
	"fmt"
	"misc/nintasm/tokenizer/tokenizerSpec"
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

func (p *LabelOperandParser) Process(operationType tokenizerSpec.TokenType, operationValue string, operationLabel string) {
	isLocal := strings.HasPrefix(operationLabel, ".")

	switch operationType {
	case tokenizerSpec.IDENTIFIER:
		if isLocal {
			if p.parentLabel == "" {
				fmt.Println("No parent label!")
			}
		} else {
			p.parentLabel = operationLabel
		}

	case tokenizerSpec.ASSIGN_EQU, tokenizerSpec.ASSIGN_simple:
		//fmt.Println("assigning stuff")
	case tokenizerSpec.DIRECTIVE_labeled:
		//fmt.Println("label dir")
	case tokenizerSpec.DIRECTIVE_labeledBlockStart:
		fmt.Println("label dir st")
	case tokenizerSpec.DIRECTIVE_labeledBlockEnd:
		fmt.Println("label dir en")
	default:
		fmt.Println("BAD LABEL OPERATION TYPE!!!")
	}
}
