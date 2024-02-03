package parser

import (
	"errors"
	"fmt"
	enumTokenTypes "misc/nintasm/enums/tokenTypes"
	"misc/nintasm/parser/parserTypes"
)

const LINE_OPERATION_TARGET_TOKENIZER = "startLine"

// Used to determine the operation at the start of the line
type OperationParser struct {
	Parser
	operationLabel       string
	operationType        tokenEnum
	operationValue       string
	operationSimpleType  parserTypes.SimpleOperation
	operandStartPosition int
}

// Create helper
func NewOperationParser() OperationParser {
	return OperationParser{
		operationLabel:       "",
		operationType:        enumTokenTypes.None,
		operationValue:       "",
		operationSimpleType:  parserTypes.None,
		operandStartPosition: 0,
	}
}

// ++++++++++++++++++++++
// 🛠️ Get info about the successfully-parsed operation
func (p *OperationParser) GetOperationDetails() (tokenEnum, parserTypes.SimpleOperation, string, string, int) {
	return p.operationType, p.operationSimpleType, p.operationValue, p.operationLabel, p.operandStartPosition
}

// ====================================================================
func (p *OperationParser) Process(line string) (err error) {
	p.operationLabel = ""
	p.operationType = enumTokenTypes.None
	p.operationValue = ""
	p.operationSimpleType = parserTypes.None
	p.operandStartPosition = 0

	// Get tokenizer to start
	err = p.startAndAdvanceToNext(line, LINE_OPERATION_TARGET_TOKENIZER)
	if err != nil {
		return err
	}
	err = p.determineLabelOrOperation()
	if err != nil {
		return err
	}

	// 🟢 Parsing has succeeded, so get the operand start position too
	p.operandStartPosition = p.tokenizer.GetCursor()

	return nil
}

// ====================================================================
func (p *OperationParser) determineLabelOrOperation() error {
	if p.lookaheadType != enumTokenTypes.WHITESPACE {
		return p.getLabelOperation()
	}
	return p.getRegularOperation()
}

//====================================================================
// REGULAR OPERATION
// There IS whitespace at the start

func (p *OperationParser) getRegularOperation() error {
	p.eat(enumTokenTypes.WHITESPACE)
	p.advanceToNext()

	var operationSimpleType parserTypes.SimpleOperation

	switch p.lookaheadType {

	case enumTokenTypes.INSTRUCTION:
		operationSimpleType = parserTypes.Instruction
		break

	case enumTokenTypes.DELIMITER_period:
		p.eat(enumTokenTypes.DELIMITER_period)
		p.advanceToNext()
		if p.lookaheadType > enumTokenTypes.DIRECTIVE_RANGE_START && p.lookaheadType < enumTokenTypes.DIRECTIVE_RANGE_END {
			operationSimpleType = parserTypes.Directive
			break
		}
		return errors.New("UNKNOWN DIRECTIVE")
	case enumTokenTypes.IDENTIFIER:
		operationSimpleType = parserTypes.Macro
		break
	case enumTokenTypes.None:
		return errors.New("UNEXPECTED EMPTY OPERATION???")
	default:
		return errors.New("UNKNOWN OPERATION")
	}

	operationType := p.lookaheadType
	operationValue := p.lookaheadValue

	// ❔ See what's next. Should be either whitespace or nothing
	p.eat(p.lookaheadType)
	err := p.advanceToNext()
	if err != nil {
		return err
	}

	// 🟢 Operation parsing succeeds
	if p.lookaheadType == enumTokenTypes.WHITESPACE || p.lookaheadType == enumTokenTypes.None {
		p.operationValue = operationValue
		p.operationType = operationType
		p.operationSimpleType = operationSimpleType
		return nil
	}

	// ❌ Fails
	badTokenAfterOperationMessage := fmt.Sprintf("ILLEGAL token after operation: %v", p.lookaheadValue)
	return errors.New(badTokenAfterOperationMessage)

}

// ====================================================================

// LABEL OPERATION
// Line has no whitespace at the start
func (p *OperationParser) getLabelOperation() error {

	// Check for local label
	isLocal := p.lookaheadType == enumTokenTypes.DELIMITER_period
	if isLocal {
		p.eat(enumTokenTypes.DELIMITER_period)
		p.advanceToNext()
	}

	// Label itself has been determined
	operationLabel := p.lookaheadValue
	if isLocal {
		operationLabel = "." + operationLabel
	}

	//Will expect an identifier to signify a label
	err := p.eat(enumTokenTypes.IDENTIFIER)
	if err != nil {
		//⚠️ In the case of a LOCAL label, other label-likes are allowed
		if !isLocal || !p.tokenizer.IsTokenIdentifierLike(operationLabel) {
			return err
		}
	}

	err = p.advanceToNext()
	if err != nil {
		return err
	}

	//See what follows the label
	switch p.lookaheadType {
	case enumTokenTypes.None:
		// ❌ Fails if NO colon
		return errors.New("LABEL IS INCOMPLETE - WHERE'S YOUR COLON!?")

	case enumTokenTypes.DELIMITER_colon:
		p.eat(enumTokenTypes.DELIMITER_colon)
		p.advanceToNext()
		if p.lookaheadType == enumTokenTypes.None {
			// 🟢 Label parsing succeeds
			p.operationLabel = operationLabel
			p.operationType = enumTokenTypes.IDENTIFIER
			p.operationValue = ""
			p.operationSimpleType = parserTypes.Label
			return nil
		}
		// ❌ Fails if tokens follow colon
		return errors.New("STUFF FOLLOWS THE COLON!!!")
	}

	//Potential labeled directive or assignment statement
	return p.getLabelFollowup(operationLabel, false)
}

// ====================================================================

// Labeled directive
func (p *OperationParser) getLabelFollowup(operationLabel string, hadWhitespace bool) error {
	if p.lookaheadType == enumTokenTypes.WHITESPACE {
		p.eat(enumTokenTypes.WHITESPACE)
		p.advanceToNext()
		if p.lookaheadType != enumTokenTypes.None {
			return p.getLabelFollowup(operationLabel, true)
		}
	}

	var operationValue string
	var operationType tokenEnum

	switch p.lookaheadType {

	//Equals sign
	case enumTokenTypes.ASSIGN_simple:
		operationType = p.lookaheadType
		operationValue = p.lookaheadValue
		p.eat(enumTokenTypes.ASSIGN_simple)
		p.advanceToNext()

	//EQU
	case enumTokenTypes.ASSIGN_EQU:
		operationType = p.lookaheadType
		operationValue = p.lookaheadValue
		p.eat(enumTokenTypes.ASSIGN_EQU)
		p.advanceToNext()
		err := p.eat(enumTokenTypes.WHITESPACE)
		if err != nil {
			return err
		}

	//Actual directive
	case enumTokenTypes.DELIMITER_period:
		if !hadWhitespace {
			return errors.New("need spacing for labeled directive")
		}
		p.eat(enumTokenTypes.DELIMITER_period)
		p.advanceToNext()

		switch p.lookaheadType {
		case enumTokenTypes.DIRECTIVE_labeled,
			enumTokenTypes.DIRECTIVE_labeledBlockStart,
			enumTokenTypes.DIRECTIVE_labeledBlockEnd:
			operationType = p.lookaheadType
			operationValue = p.lookaheadValue
			p.eat(p.lookaheadType)
			err := p.advanceToNext()
			if err != nil {
				return err
			}

			if p.lookaheadType != enumTokenTypes.None && p.lookaheadType != enumTokenTypes.WHITESPACE {
				// ❌ Fails
				err := p.eat(enumTokenTypes.WHITESPACE)
				return err
			}

		default:
			return errors.New("Unknown labeled directive")

		}
	default:
		// ❌ Fails
		return errors.New("Illegal token for labeled operation")
	}

	// 🟢 Labeled directive parsing succeeds
	p.operationLabel = operationLabel
	p.operationType = operationType
	p.operationValue = operationValue
	p.operationSimpleType = parserTypes.Label
	return nil
}
