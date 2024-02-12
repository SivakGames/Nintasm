package parser

import (
	"errors"
	"fmt"
	enumParserTypes "misc/nintasm/enums/parserTypes"
	enumTokenTypes "misc/nintasm/enums/tokenTypes"
	"misc/nintasm/util"
)

const LINE_OPERATION_TARGET_TOKENIZER = "startLine"

// Used to determine the operation at the start of the line
type OperationParser struct {
	Parser
	operationLabel       string
	operationTokenEnum   tokenEnum
	operationTokenValue  string
	parentParserEnum     enumParserTypes.Def
	operandStartPosition int
}

// Create helper
func NewOperationParser() OperationParser {
	return OperationParser{
		operationLabel:       "",
		operationTokenEnum:   enumTokenTypes.None,
		operationTokenValue:  "",
		parentParserEnum:     enumParserTypes.None,
		operandStartPosition: 0,
	}
}

// ====================================================================
func (p *OperationParser) Process(line string) (err error) {
	p.operationLabel = ""
	p.operationTokenEnum = enumTokenTypes.None
	p.operationTokenValue = ""
	p.parentParserEnum = enumParserTypes.None
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

//+++++++++++++++++++++++++++++

func (p *OperationParser) GetLineOperationValues() util.LineOperationParsedValues {
	return util.NewLineOperationParsedValues(
		p.operandStartPosition,
		p.operationLabel,
		p.operationTokenEnum,
		p.operationTokenValue,
		p.parentParserEnum,
	)
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

	var operationSimpleType enumParserTypes.Def

	switch p.lookaheadType {

	case enumTokenTypes.INSTRUCTION:
		operationSimpleType = enumParserTypes.Instruction
		break

	case enumTokenTypes.DELIMITER_period:
		err := p.eatFreelyAndAdvance(enumTokenTypes.DELIMITER_period)
		if err != nil {
			return err
		}

		if p.lookaheadType > enumTokenTypes.DIRECTIVE_RANGE_START && p.lookaheadType < enumTokenTypes.DIRECTIVE_RANGE_END {
			operationSimpleType = enumParserTypes.Directive
			break
		}

		return errors.New("UNKNOWN DIRECTIVE")
	case enumTokenTypes.IDENTIFIER:
		operationSimpleType = enumParserTypes.Macro
		break
	case enumTokenTypes.None:
		return errors.New("UNEXPECTED EMPTY OPERATION???")
	default:
		return errors.New("UNKNOWN OPERATION")
	}

	operationTokenEnum := p.lookaheadType
	operationTokenValue := p.lookaheadValue

	// ❔ See what's next. Should be either whitespace or nothing
	p.eat(p.lookaheadType)
	err := p.advanceToNext()
	if err != nil {
		return err
	}

	// 🟢 Operation parsing succeeds
	if p.lookaheadType == enumTokenTypes.WHITESPACE || p.lookaheadType == enumTokenTypes.None {
		p.operationTokenValue = operationTokenValue
		p.operationTokenEnum = operationTokenEnum
		p.parentParserEnum = operationSimpleType
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
			p.operationTokenEnum = enumTokenTypes.IDENTIFIER
			p.operationTokenValue = ""
			p.parentParserEnum = enumParserTypes.Label
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

	var operationTokenValue string
	var operationTokenEnum tokenEnum

	switch p.lookaheadType {

	//Equals sign
	case enumTokenTypes.ASSIGN_simple:
		operationTokenEnum = p.lookaheadType
		operationTokenValue = p.lookaheadValue
		p.eat(enumTokenTypes.ASSIGN_simple)
		p.advanceToNext()

	//EQU
	case enumTokenTypes.ASSIGN_EQU:
		operationTokenEnum = p.lookaheadType
		operationTokenValue = p.lookaheadValue
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
			operationTokenEnum = p.lookaheadType
			operationTokenValue = p.lookaheadValue
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
	p.operationTokenEnum = operationTokenEnum
	p.operationTokenValue = operationTokenValue
	p.parentParserEnum = enumParserTypes.Label
	return nil
}
