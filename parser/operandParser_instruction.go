package parser

import (
	"errors"
	"fmt"
	enumInstructionModes "misc/nintasm/enums/instructionModes"
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
	var instructionMode instModes
	var operand *Node = nil

	instructionName := strings.ToUpper(operationValue)

	isBranch := false
	instructionXYIndex := tokenizerSpec.None
	opcodesAndSupportedModes := instructionData.OpcodesAndSupportedModes[instructionName]
	operandList := []Node{}
	useInstructionMode, useInstructionZPMode := enumInstructionModes.None, enumInstructionModes.None

	switch p.lookaheadType {
	//Try no operand
	case tokenizerSpec.None:
		instructionMode = enumInstructionModes.IMPL
	//Try A (for shifts)
	case tokenizerSpec.REGISTER_A:
		instructionMode = enumInstructionModes.A
		err = p.eatFreelyAndAdvance(tokenizerSpec.REGISTER_A)
		if err != nil {
			return err // ❌ Fails
		}
		if p.lookaheadType != tokenizerSpec.None {
			return errors.New("No tokens can follow the A")
		}

	default:
		isBranch = checkIfBranchInstruction(instructionName)
		operandList, err = p.GetOperandList()
		if err != nil {
			return err // ❌ Fails
		}
		instructionXYIndex = p.instructionXYIndex
		instructionMode = p.instructionMode
		if isBranch && instructionMode == enumInstructionModes.ABS {
			instructionMode = enumInstructionModes.REL
		}
	}

	//If an index is present, see if it's usable with desired mode and set it
	if instructionXYIndex == tokenizerSpec.REGISTER_X || instructionXYIndex == tokenizerSpec.REGISTER_Y {
		instructionMode, err = checkModeSupportsXY(instructionMode, instructionXYIndex)
		if err != nil {
			return err // ❌ Fails
		}
	} else if instructionXYIndex != tokenizerSpec.None {
		panic("MAJOR ERROR with XY index checks!!!")
	}

	if len(operandList) == 1 {
		operand = &operandList[0]
	}

	//Used for auto ZP convert if possible
	instructionZPModeEquivalent := getZeroPageEquivalent(instructionMode)

	// Check if instruction itself supports mode
	for _, m := range *opcodesAndSupportedModes.Modes {
		if m == instructionMode {
			useInstructionMode = m
			continue
		}
		if m == instructionZPModeEquivalent && operand != nil && operand.Resolved {
			if operand.AsNumber <= 255 && operand.AsNumber >= 0 {
				useInstructionZPMode = m
				fmt.Println("Yep")
			}
		}
	}

	if useInstructionMode == enumInstructionModes.None {
		return errors.New("Mode is not supported by instruction!") // ❌ Fails
	}
	//Overwrite mode with ZP version if possible
	if useInstructionZPMode != enumInstructionModes.None {
		instructionMode = useInstructionZPMode
	}

	return nil
}

// +++++++++++++++++

// Helper to see if instruction is a branch instruction. Only branch instructions contain REL mode
func checkIfBranchInstruction(instructionName string) bool {
	// What modes this instruction can use
	opcodesAndSupportedModes := instructionData.OpcodesAndSupportedModes[instructionName]

	for _, v := range *opcodesAndSupportedModes.Modes {
		if v == enumInstructionModes.REL {
			return true
		}
	}
	return false
}

// +++++++++++++++++

// Operand wants to use an X/Y index. See if index is used properly and eat it
func checkModeSupportsXY(instructionMode instModes, instructionIndex tokenizerSpec.TokenType) (instModes, error) {
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
func getZeroPageEquivalent(instructionMode instModes) instModes {
	zpMode, ok := instructionData.ABStoZP[instructionMode]
	if ok {
		return zpMode
	}
	return enumInstructionModes.None
}
