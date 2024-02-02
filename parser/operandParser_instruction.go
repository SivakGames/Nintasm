package parser

import (
	"errors"
	"fmt"
	"misc/nintasm/instructionData"
	"misc/nintasm/interpreter"
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
	// The instruction itself (in upper case)
	instruction := strings.ToUpper(operationValue)
	// What modes this instruction can use
	allowedModesForInstruction := instructionData.AllowedModes[instruction]

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

	// TODO: Convert to ZP?
	// TODO: Check if instruction itself supports mode

	nodez := interpreter.InterpretOperands(operandList)
	_ = nodez
	//fmt.Println("PARSING SUCCESS", nodez[0])

	return
}

// See what token (if any) precedes the operand. This will determine general mode...
//func (p *InstructionOperandParser) getOperandModeByLeadToken(isBranch bool) (instructionData.InstructionModes, []Node, tokenizerSpec.TokenType, error) {
//instructionMode := instructionData.None
//instructionXYindex := tokenizerSpec.None
//var err error = nil
//var operandList []Node

//switch p.lookaheadType {

// --------------------------------------------------------
// [ for indirect addressing.
// *This is more complex in assuming an operand wrapped in braces
//  so most handling is done here.  This will also directly return the operand

//	case tokenizerSpec.DELIMITER_leftSquareBracket:
//		instructionMode = instructionData.IND
//		var operand Node
//
//		err = p.eatFreelyAndAdvance(tokenizerSpec.DELIMITER_leftSquareBracket)
//		if err != nil {
//			return instructionMode, operandList, instructionXYindex, err // ❌ Fails
//		}
//
//		//Get the operand
//		operand, err = p.Statement()
//		if err != nil {
//			return instructionMode, operandList, instructionXYindex, err // ❌ Fails
//		}
//		operandList = append(operandList, operand)
//
//		switch p.lookaheadType {
//		case tokenizerSpec.None:
//			return instructionMode, operandList, instructionXYindex, errors.New("Indirect End of input") // ❌ Fails
//
//		// xxxxxxxxxxxxxxx
//		// Indirect X
//		case tokenizerSpec.DELIMITER_comma:
//			err = p.checkValidXYIndexes(tokenizerSpec.REGISTER_X)
//			if err != nil {
//				return instructionMode, operandList, instructionXYindex, err // ❌ Fails
//			}
//			err = p.eatFreelyAndAdvance(tokenizerSpec.DELIMITER_rightSquareBracket)
//			if err != nil {
//				return instructionMode, operandList, instructionXYindex, err // ❌ Fails
//			}
//			instructionXYindex = tokenizerSpec.REGISTER_X
//			//instructionMode = instructionData.IND_X
//			return instructionMode, operandList, instructionXYindex, nil // 🟢 Indirect X Succeeds
//
//		// yyyyyyyyyyyyyy
//		// Indirect only or Indirect Y
//		case tokenizerSpec.DELIMITER_rightSquareBracket:
//			err = p.eatFreelyAndAdvance(tokenizerSpec.DELIMITER_rightSquareBracket)
//			if err != nil {
//				return instructionMode, operandList, instructionXYindex, err // ❌ Fails
//			}
//			if p.lookaheadType != tokenizerSpec.None {
//				err = p.checkValidXYIndexes(tokenizerSpec.REGISTER_Y)
//				if err != nil {
//					return instructionMode, operandList, instructionXYindex, err // ❌ Fails
//				}
//				instructionXYindex = tokenizerSpec.REGISTER_Y
//			}
//			return instructionMode, operandList, instructionXYindex, nil // 🟢 Indirect or Indirect Y Succeeds
//
//		default:
//			return instructionMode, operandList, instructionXYindex, errors.New("Unknown token for indirect op") // ❌ Fails
//		}

// --------------------------------------------------------
//# for immediate mode

//	case tokenizerSpec.DELIMITER_hash:
//		err = p.eatFreelyAndAdvance(tokenizerSpec.DELIMITER_hash)
//		if err != nil {
//			return instructionMode, operandList, instructionXYindex, err // ❌ Fails
//		}
//		instructionMode = instructionData.IMM

// --------------------------------------------------------
// < for enforced ZP (if auto ZP is off this will still force ZP if desired)

//	case tokenizerSpec.OPERATOR_relational:
//		if p.lookaheadValue == "<" {
//			err = p.eatFreelyAndAdvance(tokenizerSpec.OPERATOR_relational)
//			if err != nil {
//				return instructionMode, operandList, instructionXYindex, err // ❌ Fails
//			}
//			instructionMode = instructionData.ZP
//		}
// ❌ Will fail later - main parser will catch the bad token

// --------------------------------------------------------
// Anything else is absolute or relative (branches)

/* default:
	if isBranch {
		instructionMode = instructionData.REL
	} else {
		instructionMode = instructionData.ABS
	}
} */

// ````````````````````````````````````````````````````````
// Finally, parse the operand for non-indirect modes

// --------------------------------------------------------
// Check for followups (X,Y indexes)

/*	switch p.lookaheadType {
	case tokenizerSpec.None:
		return instructionMode, operandList, instructionXYindex, nil // 🟢 No index Succeeds
	case tokenizerSpec.DELIMITER_comma:
		err = p.eatAndAdvance(tokenizerSpec.DELIMITER_comma)
		if err != nil {
			return instructionMode, operandList, instructionXYindex, err // ❌ Fails
		}
		err = p.checkIfXYIndes()
		if err != nil {
			return instructionMode, operandList, instructionXYindex, err // ❌ Fails
		}
		instructionXYindex = p.lookaheadType
		err = p.eatAndAdvance(instructionXYindex)
		if err != nil {
			return instructionMode, operandList, instructionXYindex, err // ❌ Fails
		}
		return instructionMode, operandList, instructionXYindex, nil // 🟢 Index Succeeds
	} */

//return instructionMode, operandList, instructionXYindex, nil // 🟢 Succeeds for now

//}

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
/* func (p *InstructionOperandParser) checkValidXYIndexes(targetIndex tokenizerSpec.TokenType) error {

	err := p.eatAndAdvance(tokenizerSpec.DELIMITER_comma)
	if err != nil {
		return err
	}
	if p.lookaheadType == tokenizerSpec.None {
		return errors.New("Index is MISSING!")
	}
	err = p.checkIfXYIndes()
	if err != nil {
		return err
	}
	if p.lookaheadType != targetIndex {
		return errors.New("Wrong index used for operation")
	}
	err = p.eatAndAdvance(targetIndex)
	if err != nil {
		return err
	}

	return nil
} */

// +++++++++++++++++

// Operand wants to use an X/Y index. See if index is used properly and eat it
/* func (p *InstructionOperandParser) checkIfXYIndes() error {
	if p.lookaheadType != tokenizerSpec.REGISTER_X && p.lookaheadType != tokenizerSpec.REGISTER_Y {
		return errors.New("Invalid index value")
	}
	return nil
}*/

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
