package parser

import (
	"errors"
	"fmt"
	"misc/nintasm/instructionData"
	"misc/nintasm/tokenizer/tokenizerSpec"
	"strings"
)

type InstructionOperandParser struct {
	OperandParser
}

func NewInstructionOperandParser() InstructionOperandParser {
	return InstructionOperandParser{}
}

func (p *InstructionOperandParser) Process(operationValue string) {
	var err error = nil
	var instructionMode instructionData.InstructionModes

	operandList := []Node{}
	// The instructionName itself (in upper case)
	instructionName := strings.ToUpper(operationValue)
	// What modes this instruction can use
	allowedModesForInstruction := instructionData.AllowedModes[instructionName]

	isBranch := false
	instructionXYIndex := tokenizerSpec.None

	switch p.lookaheadType {
	//Try no operand
	case tokenizerSpec.None:
		instructionMode = instructionData.IMPL
	//Try A (for shifts)
	case tokenizerSpec.REGISTER_A:
		instructionMode = instructionData.A
		err = p.eatFreelyAndAdvance(tokenizerSpec.REGISTER_A)
		if err != nil {
			fmt.Println(err) // ❌ Fails
			return
		}
	default:
		isBranch = p.checkIfBranchInstruction(&allowedModesForInstruction.Modes)
		operandList, err = p.GetOperandList()
		if err != nil {
			fmt.Println(err)
			return
		}
		instructionXYIndex = p.instructionXYIndex
		instructionMode = p.instructionMode
	}

	_ = isBranch
	_ = instructionMode

	//If index is present, see if it's usable with desired mode

	if instructionXYIndex == tokenizerSpec.REGISTER_X || instructionXYIndex == tokenizerSpec.REGISTER_Y {
		instructionMode, err = p.checkModeSupportsXY(instructionMode, instructionXYIndex)
		if err != nil {
			fmt.Println(err)
			return
		}
	} else if instructionXYIndex != tokenizerSpec.None {
		fmt.Println("MAJOR ERROR!!!")
	}

	fmt.Println(instructionData.InstructionModeDetails[instructionMode])
	fmt.Println(operandList)

	// TODO: Convert to ZP?
	// TODO: Check if instruction itself supports mode

	//fmt.Println("PARSING SUCCESS", nodez[0])

	return
}

// +++++++++++++++++

// Helper to see if instruction is a branch instruction. Only branch instructions contain REL mode
func (p *InstructionOperandParser) checkIfBranchInstruction(modes *[]instructionData.InstructionModes) bool {
	for _, v := range *modes {
		if v == instructionData.REL {
			return true
		}
	}
	return false
}

// +++++++++++++++++
// Operand wants to use an X/Y index. See if index is used properly and eat it
func (p *InstructionOperandParser) checkModeSupportsXY(instructionMode instructionData.InstructionModes, instructionIndex tokenizerSpec.TokenType) (instructionData.InstructionModes, error) {
	xyMode, ok := instructionData.ModesWithXYIndexes[instructionMode]

	if ok {
		if instructionIndex == tokenizerSpec.REGISTER_X {
			return xyMode.X, nil
		}
		return xyMode.Y, nil
	}

	return instructionMode, errors.New("X or Y indexes cannot be used with target mode")
}
