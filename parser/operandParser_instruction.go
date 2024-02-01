package parser

import (
	"errors"
	"fmt"
	"misc/nintasm/instructionData"
	"misc/nintasm/parser/operandFactory"
	"misc/nintasm/tokenizer"
	"misc/nintasm/tokenizer/tokenizerSpec"
	"strings"
)

type InstructionOperandParser struct {
	OperandParser
	altTokenizer tokenizer.Tokenizer
}

func NewInstructionOperandParser() InstructionOperandParser {
	return InstructionOperandParser{
		altTokenizer: tokenizer.New(),
	}
}

func (p *InstructionOperandParser) Process(operationValue string) {
	operand := operandFactory.EmptyNode()
	instructionIndex := tokenizerSpec.None
	// The instruction itself (in upper case)
	instruction := strings.ToUpper(operationValue)
	// What modes this instruction can use
	allowedModesForInstruction := instructionData.AllowedModes[instruction]

	isBranch := false
	var err error = nil

	// What mode will be used
	var instructionMode instructionData.InstructionModes

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
		instructionMode, operand, instructionIndex, err = p.getOperandModeByLeadToken(isBranch)
	}

	if p.lookaheadType != tokenizerSpec.None {
		fmt.Println("NO MORE TOKENS allowed", p.lookaheadValue) // ❌ Fails
		return
	}

	//If index is present, see if it's usable with desired mode

	if instructionIndex == tokenizerSpec.REGISTER_X || instructionIndex == tokenizerSpec.REGISTER_Y {
		p.checkModeSupportsXY(instructionMode, instructionIndex)
	} else if instructionIndex != tokenizerSpec.None {
		fmt.Println("MAJOR ERROR!!!")
	}

	// Convert to ZP?

	// Check if instruction itself supports mode

	fmt.Println("SUCCESS")
	fmt.Println(instructionMode, operand, instructionIndex)

	return
}

// See what token (if any) precedes the operand. This will determine general mode...
func (p *InstructionOperandParser) getOperandModeByLeadToken(isBranch bool) (instructionData.InstructionModes, operandFactory.Node, tokenizerSpec.TokenType, error) {
	instructionMode := instructionData.None
	instructionXYindex := tokenizerSpec.None
	var err error = nil
	var operand operandFactory.Node = operandFactory.EmptyNode()

	switch p.lookaheadType {

	// --------------------------------------------------------
	// [ for indirect addressing.
	// *This is more complex in assuming an operand wrapped in braces
	//  so most handling is done here.  This will also directly return the operand

	case tokenizerSpec.DELIMITER_leftSquareBracket:
		instructionMode = instructionData.IND

		err = p.eatFreelyAndAdvance(tokenizerSpec.DELIMITER_leftSquareBracket)
		if err != nil {
			return instructionMode, operand, instructionXYindex, err // ❌ Fails
		}

		//Get the operand
		operand, err = p.Statement()
		if err != nil {
			return instructionMode, operand, instructionXYindex, err // ❌ Fails
		}

		switch p.lookaheadType {
		case tokenizerSpec.None:
			return instructionMode, operand, instructionXYindex, errors.New("Indirect End of input") // ❌ Fails

		// xxxxxxxxxxxxxxx
		// Indirect X
		case tokenizerSpec.DELIMITER_comma:
			err = p.checkValidXYIndexes(tokenizerSpec.REGISTER_X)
			if err != nil {
				return instructionMode, operand, instructionXYindex, err // ❌ Fails
			}
			err = p.eatFreelyAndAdvance(tokenizerSpec.DELIMITER_rightSquareBracket)
			if err != nil {
				return instructionMode, operand, instructionXYindex, err // ❌ Fails
			}
			instructionXYindex = tokenizerSpec.REGISTER_X
			//instructionMode = instructionData.IND_X
			return instructionMode, operand, instructionXYindex, nil // 🟢 Indirect X Succeeds

		// yyyyyyyyyyyyyy
		// Indirect only or Indirect Y
		case tokenizerSpec.DELIMITER_rightSquareBracket:
			err = p.eatFreelyAndAdvance(tokenizerSpec.DELIMITER_rightSquareBracket)
			if err != nil {
				return instructionMode, operand, instructionXYindex, err // ❌ Fails
			}
			if p.lookaheadType != tokenizerSpec.None {
				err = p.checkValidXYIndexes(tokenizerSpec.REGISTER_Y)
				if err != nil {
					return instructionMode, operand, instructionXYindex, err // ❌ Fails
				}
				instructionXYindex = tokenizerSpec.REGISTER_Y
			}
			return instructionMode, operand, instructionXYindex, nil // 🟢 Indirect or Indirect Y Succeeds

		default:
			return instructionMode, operand, instructionXYindex, errors.New("Unknown token for indirect op") // ❌ Fails
		}

	// --------------------------------------------------------
	//# for immediate mode

	case tokenizerSpec.DELIMITER_hash:
		err = p.eatFreelyAndAdvance(tokenizerSpec.DELIMITER_hash)
		if err != nil {
			return instructionMode, operand, instructionXYindex, err // ❌ Fails
		}
		instructionMode = instructionData.IMM

	// --------------------------------------------------------
	// < for enforced ZP (if auto ZP is off this will still force ZP if desired)

	case tokenizerSpec.OPERATOR_relational:
		if p.lookaheadValue == "<" {
			err = p.eatFreelyAndAdvance(tokenizerSpec.OPERATOR_relational)
			if err != nil {
				return instructionMode, operand, instructionXYindex, err // ❌ Fails
			}
			instructionMode = instructionData.ZP
		}
		// ❌ Will fail later - main parser will catch the bad token

	// --------------------------------------------------------
	// Anything else is absolute or relative (branches)

	default:
		if isBranch {
			instructionMode = instructionData.REL
		} else {
			instructionMode = instructionData.ABS
		}
	}

	// ````````````````````````````````````````````````````````
	// Finally, parse the operand for non-indirect modes

	operand, err = p.GetFirstOperandOnly()
	if err != nil {
		return instructionMode, operand, instructionXYindex, err // ❌ Fails
	}

	// --------------------------------------------------------
	// Check for followups (X,Y indexes)

	switch p.lookaheadType {
	case tokenizerSpec.None:
		return instructionMode, operand, instructionXYindex, nil // 🟢 No index Succeeds
	case tokenizerSpec.DELIMITER_comma:
		err = p.eatAndAdvance(tokenizerSpec.DELIMITER_comma)
		if err != nil {
			return instructionMode, operand, instructionXYindex, err // ❌ Fails
		}
		err = p.checkIfXYIndes()
		if err != nil {
			return instructionMode, operand, instructionXYindex, err // ❌ Fails
		}
		instructionXYindex = p.lookaheadType
		err = p.eatAndAdvance(instructionXYindex)
		if err != nil {
			return instructionMode, operand, instructionXYindex, err // ❌ Fails
		}
		return instructionMode, operand, instructionXYindex, nil // 🟢 Index Succeeds

	}

	return instructionMode, operand, instructionXYindex, nil // 🟢 Succeeds for now

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
func (p *InstructionOperandParser) checkValidXYIndexes(targetIndex tokenizerSpec.TokenType) error {

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
}

// +++++++++++++++++

// Operand wants to use an X/Y index. See if index is used properly and eat it
func (p *InstructionOperandParser) checkIfXYIndes() error {
	if p.lookaheadType != tokenizerSpec.REGISTER_X && p.lookaheadType != tokenizerSpec.REGISTER_Y {
		return errors.New("Invalid index value")
	}
	return nil
}

// +++++++++++++++++
// Operand wants to use an X/Y index. See if index is used properly and eat it
func (p *InstructionOperandParser) checkModeSupportsXY(instructionMode instructionData.InstructionModes, instructionIndex tokenizerSpec.TokenType) (instructionData.InstructionModes, error) {
	switch instructionMode {
	case instructionData.ABS:
		if instructionIndex == tokenizerSpec.REGISTER_X {
			return instructionData.ABS_X, nil
		}
		return instructionData.ABS_Y, nil
	case instructionData.IND:
		if instructionIndex == tokenizerSpec.REGISTER_X {
			return instructionData.IND_X, nil
		}
		return instructionData.IND_Y, nil
	case instructionData.ZP:
		if instructionIndex == tokenizerSpec.REGISTER_X {
			return instructionData.ZP_X, nil
		}
		return instructionData.ZP_Y, nil

	default:
		return instructionMode, errors.New("Index cannot be used with target mode")
	}
}
