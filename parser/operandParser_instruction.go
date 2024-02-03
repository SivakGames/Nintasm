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

func (p *InstructionOperandParser) Process(operationValue string) error {
	var err error = nil
	var instructionMode instructionData.InstructionModes

	operandList := []Node{}
	// The instructionName itself (in upper case)
	instructionName := strings.ToUpper(operationValue)

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
			return err // ❌ Fails
		}
		if p.lookaheadType != tokenizerSpec.None {
			fmt.Println("No tokens can follow the A")
		}

	default:
		isBranch = checkIfBranchInstruction(instructionName)
		operandList, err = p.GetOperandList()
		if err != nil {
			return err // ❌ Fails
		}
		instructionXYIndex = p.instructionXYIndex
		instructionMode = p.instructionMode
		if isBranch && instructionMode == instructionData.ABS {
			instructionMode = instructionData.REL
		}
	}

	_ = operandList

	//If an index is present, see if it's usable with desired mode and set it
	if instructionXYIndex == tokenizerSpec.REGISTER_X || instructionXYIndex == tokenizerSpec.REGISTER_Y {
		instructionMode, err = checkModeSupportsXY(instructionMode, instructionXYIndex)
		if err != nil {
			return err // ❌ Fails
		}
	} else if instructionXYIndex != tokenizerSpec.None {
		panic("MAJOR ERROR with XY index checks!!!")
	}

	//Used for auto ZP convert if possible
	instructionZPModeEquivalent := getZeroPageEquivalent(instructionMode)
	opcodesAndSupportedModes := instructionData.OpcodesAndSupportedModes[instructionName]
	useInstructionMode, useInstructionZPMode := instructionData.None, instructionData.None

	// TODO: Check if instruction itself supports mode
	for _, m := range *opcodesAndSupportedModes.Modes {
		if m == instructionMode {
			useInstructionMode = m
			continue
		}
		if m == instructionZPModeEquivalent {
			useInstructionZPMode = m
		}
	}

	if useInstructionMode == instructionData.None {
		return errors.New("Mode is not supported by instruction!") // ❌ Fails
	}
	//Overwrite mode with ZP version if possible
	if useInstructionZPMode != instructionData.None {
		instructionMode = useInstructionZPMode
	}

	//fmt.Println("PARSING SUCCESS", nodez[0])

	return nil
}

// +++++++++++++++++

// Helper to see if instruction is a branch instruction. Only branch instructions contain REL mode
func checkIfBranchInstruction(instructionName string) bool {
	// What modes this instruction can use
	opcodesAndSupportedModes := instructionData.OpcodesAndSupportedModes[instructionName]

	for _, v := range *opcodesAndSupportedModes.Modes {
		if v == instructionData.REL {
			return true
		}
	}
	return false
}

// +++++++++++++++++

// Operand wants to use an X/Y index. See if index is used properly and eat it
func checkModeSupportsXY(instructionMode instructionData.InstructionModes, instructionIndex tokenizerSpec.TokenType) (instructionData.InstructionModes, error) {
	xyMode, ok := instructionData.ModesWithXYIndexes[instructionMode]

	if ok {
		if instructionIndex == tokenizerSpec.REGISTER_X {
			return xyMode.X, nil
		}
		return xyMode.Y, nil
	}

	return instructionMode, errors.New("X or Y indexes cannot be used with target mode")
}

// If absolute mode (X,Y too) get the ZP version for auto convert. Get None otherwise
func getZeroPageEquivalent(instructionMode instructionData.InstructionModes) instructionData.InstructionModes {
	zpMode, ok := instructionData.ABStoZP[instructionMode]
	if ok {
		return zpMode
	}
	return instructionData.None
}
