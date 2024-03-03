package parser

import (
	"misc/nintasm/assemble/errorHandler"
	enumErrorCodes "misc/nintasm/constants/enums/errorCodes"
	"misc/nintasm/parser/instructionHandler"
	"strings"
)

const INSTRUCTION_MIN_OPERANDS = 0
const INSTRUCTION_MAX_OPERANDS = 1
const INSTRUCTION_MANAULLY_EVALS = true

type InstructionOperandParser struct {
	OperandParser
}

// Constructor
func NewInstructionOperandParser() InstructionOperandParser {
	return InstructionOperandParser{}
}

// Main instruction parsing
func (p *InstructionOperandParser) Process(operationValue string) error {
	instructionName := strings.ToUpper(operationValue)
	operandList, err := p.GetOperandList(
		INSTRUCTION_MIN_OPERANDS, INSTRUCTION_MAX_OPERANDS, INSTRUCTION_MANAULLY_EVALS,
		[]string{"instruction"},
	)
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
