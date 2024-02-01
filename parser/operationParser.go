package parser

import (
	"errors"
	"fmt"
	"misc/nintasm/parser/parserTypes"
	"misc/nintasm/tokenizer/tokenizerSpec"
)

const LINE_OPERATION_TARGET_TOKENIZER = "startLine"

// Used to determine the operation at the start of the line
type OperationParser struct {
	Parser
	operationLabel       string
	operationType        tokenizerSpec.TokenType
	operationValue       string
	operationSimpleType  parserTypes.SimpleOperation
	operandStartPosition int
}

// Create helper
func NewOperationParser() OperationParser {
	return OperationParser{
		operationLabel:       "",
		operationType:        tokenizerSpec.None,
		operationValue:       "",
		operationSimpleType:  parserTypes.None,
		operandStartPosition: 0,
	}
}

// ++++++++++++++++++++++
// 🛠️ Get info about the successfully-parsed operation
func (p *OperationParser) GetOperationDetails() (tokenizerSpec.TokenType, parserTypes.SimpleOperation, string, string, int) {
	return p.operationType, p.operationSimpleType, p.operationValue, p.operationLabel, p.operandStartPosition
}

// ====================================================================
func (p *OperationParser) Process(line string) (err error) {
	p.operationLabel = ""
	p.operationType = tokenizerSpec.None
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
	if p.lookaheadType != tokenizerSpec.WHITESPACE {
		return p.getLabelOperation()
	}
	return p.getRegularOperation()
}

//====================================================================
// REGULAR OPERATION
// There IS whitespace at the start

func (p *OperationParser) getRegularOperation() error {
	p.eat(tokenizerSpec.WHITESPACE)
	p.advanceToNext()

	var operationSimpleType parserTypes.SimpleOperation

	switch p.lookaheadType {

	case tokenizerSpec.INSTRUCTION:
		operationSimpleType = parserTypes.Instruction
		break

	case tokenizerSpec.DELIMITER_period:
		p.eat(tokenizerSpec.DELIMITER_period)
		p.advanceToNext()
		if p.lookaheadType > tokenizerSpec.DIRECTIVE_RANGE_START && p.lookaheadType < tokenizerSpec.DIRECTIVE_RANGE_END {
			operationSimpleType = parserTypes.Directive
			break
		}
		return errors.New("UNKNOWN DIRECTIVE")
	case tokenizerSpec.IDENTIFIER:
		operationSimpleType = parserTypes.Macro
		break
	case tokenizerSpec.None:
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
	if p.lookaheadType == tokenizerSpec.WHITESPACE || p.lookaheadType == tokenizerSpec.None {
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
	isLocal := p.lookaheadType == tokenizerSpec.DELIMITER_period
	if isLocal {
		p.eat(tokenizerSpec.DELIMITER_period)
		p.advanceToNext()
	}

	// Label itself has been determined
	operationLabel := p.lookaheadValue
	if isLocal {
		operationLabel = "." + operationLabel
	}

	//Will expect an identifier to signify a label
	err := p.eat(tokenizerSpec.IDENTIFIER)
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
	case tokenizerSpec.None:
		// ❌ Fails if NO colon
		return errors.New("LABEL IS INCOMPLETE - WHERE'S YOUR COLON!?")

	case tokenizerSpec.DELIMITER_colon:
		p.eat(tokenizerSpec.DELIMITER_colon)
		p.advanceToNext()
		if p.lookaheadType == tokenizerSpec.None {
			// 🟢 Label parsing succeeds
			p.operationLabel = operationLabel
			p.operationType = tokenizerSpec.IDENTIFIER
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
	if p.lookaheadType == tokenizerSpec.WHITESPACE {
		p.eat(tokenizerSpec.WHITESPACE)
		p.advanceToNext()
		if p.lookaheadType != tokenizerSpec.None {
			return p.getLabelFollowup(operationLabel, true)
		}
	}

	var operationValue string
	var operationType tokenizerSpec.TokenType

	switch p.lookaheadType {

	//Equals sign
	case tokenizerSpec.ASSIGN_simple:
		operationType = p.lookaheadType
		operationValue = p.lookaheadValue
		p.eat(tokenizerSpec.ASSIGN_simple)
		p.advanceToNext()

	//EQU
	case tokenizerSpec.ASSIGN_EQU:
		operationType = p.lookaheadType
		operationValue = p.lookaheadValue
		p.eat(tokenizerSpec.ASSIGN_EQU)
		p.advanceToNext()
		err := p.eat(tokenizerSpec.WHITESPACE)
		if err != nil {
			return err
		}

	//Actual directive
	case tokenizerSpec.DELIMITER_period:
		if !hadWhitespace {
			return errors.New("need spacing for labeled directive")
		}
		p.eat(tokenizerSpec.DELIMITER_period)
		p.advanceToNext()

		switch p.lookaheadType {
		case tokenizerSpec.DIRECTIVE_labeled,
			tokenizerSpec.DIRECTIVE_labeledBlockStart,
			tokenizerSpec.DIRECTIVE_labeledBlockEnd:
			operationType = p.lookaheadType
			operationValue = p.lookaheadValue
			p.eat(p.lookaheadType)
			err := p.advanceToNext()
			if err != nil {
				return err
			}

			if p.lookaheadType != tokenizerSpec.None && p.lookaheadType != tokenizerSpec.WHITESPACE {
				// ❌ Fails
				err := p.eat(tokenizerSpec.WHITESPACE)
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
