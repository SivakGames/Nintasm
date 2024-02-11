package parser

import (
	"fmt"
	enumTokenTypes "misc/nintasm/enums/tokenTypes"
)

type BlockType struct {
}

type CaptureBlockOperandParser struct {
	OperandParser
}

var BlockStack = []string{}

func NewCaptureBlockOperandParser() CaptureBlockOperandParser {
	return CaptureBlockOperandParser{}
}

// Main directive parser
func (p *CaptureBlockOperandParser) Process(operationType tokenEnum, operationValue string) error {
	var err error
	if operationType == enumTokenTypes.DIRECTIVE_blockStart {
		fmt.Println("HeyyyyyyYY")
	} else {
		return err
	}

	return nil
}

// For repeat, upon hitting an endrepeat, the block at the top of the stack would need to be expanded.
// All instances of the top iterator (if present) will be filled
