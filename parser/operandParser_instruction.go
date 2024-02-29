package parser

import (
	"misc/nintasm/assemble/errorHandler"
	enumErrorCodes "misc/nintasm/constants/enums/errorCodes"
	"misc/nintasm/parser/instructionHandler"
	"strings"
)

type InstructionOperandParser struct {
	OperandParser
}

// Constructor
func NewInstructionOperandParser() InstructionOperandParser {
	return InstructionOperandParser{}
}

func (p *InstructionOperandParser) Process(operationValue string) error {
	instructionName := strings.ToUpper(operationValue)
	operandList, err := p.GetOperandList(0, 1, false, []string{"instruction"})
	if err != nil {
		return err // ❌ Fails
	}

	instructionXYIndex := p.instructionXYIndex
	instructionMode := p.instructionMode

	err = instructionHandler.EvaluateInstruction(instructionName, &operandList, instructionMode, instructionXYIndex)
	if err != nil {
		err := errorHandler.CheckErrorContinuesUpwardPropagation(err, enumErrorCodes.Error)
		if err != nil {
			return err // ❌❌ CONTINUES Failing!
		}
	}
	return nil
}
