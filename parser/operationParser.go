package parser

import (
	"misc/nintasm/assemble/errorHandler"
	enumErrorCodes "misc/nintasm/constants/enums/errorCodes"
	enumParserTypes "misc/nintasm/constants/enums/parserTypes"
	enumTokenTypes "misc/nintasm/constants/enums/tokenTypes"
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
	err := p.eatFreelyAndAdvance(enumTokenTypes.WHITESPACE)
	if err != nil {
		return err
	}

	var parentParserEnum enumParserTypes.Def

	switch p.lookaheadType {

	case enumTokenTypes.INSTRUCTION:
		parentParserEnum = enumParserTypes.Instruction
		break

	case enumTokenTypes.DELIMITER_period:
		err := p.eatFreelyAndAdvance(enumTokenTypes.DELIMITER_period)
		if err != nil {
			return err
		}

		if p.lookaheadType > enumTokenTypes.DIRECTIVE_RANGE_START && p.lookaheadType < enumTokenTypes.DIRECTIVE_RANGE_END {
			parentParserEnum = enumParserTypes.Directive
			break
		}
		return errorHandler.AddNew(enumErrorCodes.OperationDirectiveUnknown) // ❌ Fails

	case enumTokenTypes.IDENTIFIER:
		parentParserEnum = enumParserTypes.Macro
		break
	case enumTokenTypes.None:
		return errorHandler.AddNew(enumErrorCodes.OperationEmpty) // ❌ Fails
	default:
		return errorHandler.AddNew(enumErrorCodes.OperationUNKNOWN) // ❌ Fails
	}

	operationTokenEnum := p.lookaheadType
	operationTokenValue := p.lookaheadValue

	// ❔ See what's next. Should be either whitespace or nothing
	err = p.eatAndAdvance(p.lookaheadType)
	if err != nil {
		return err
	}

	// 🟢 Operation parsing succeeds
	if p.lookaheadType == enumTokenTypes.WHITESPACE || p.lookaheadType == enumTokenTypes.None {
		p.operationTokenValue = operationTokenValue
		p.operationTokenEnum = operationTokenEnum
		p.parentParserEnum = parentParserEnum
		return nil
	}

	return errorHandler.AddNew(enumErrorCodes.OperationBadTokenAfter, p.lookaheadValue) // ❌ Fails

}

// ====================================================================

// LABEL OPERATION
// Line has no whitespace at the start
func (p *OperationParser) getLabelOperation() error {

	// Check for local label
	isLocal := p.lookaheadType == enumTokenTypes.DELIMITER_period
	if isLocal {
		err := p.eatFreelyAndAdvance(enumTokenTypes.DELIMITER_period)
		if err != nil {
			return err
		}

	}

	// Label itself has been determined
	operationLabel := p.lookaheadValue
	if isLocal {
		operationLabel = "." + operationLabel
	}

	//⭐⭐ We JUST try to eat first ⭐⭐
	//Will expect an identifier to signify a label...
	err := p.eat(enumTokenTypes.IDENTIFIER)
	if err != nil {
		//⚠️ ... but in the case of a LOCAL label, label-likes ARE allowed
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
		return errorHandler.AddNew(enumErrorCodes.OperationLabelMissingColon) // ❌ Fails if NO colon

	case enumTokenTypes.DELIMITER_colon:
		err := p.eatFreelyAndAdvance(enumTokenTypes.DELIMITER_colon)
		if err != nil {
			return err
		}

		if p.lookaheadType == enumTokenTypes.None {
			// 🟢 Label parsing succeeds
			p.operationLabel = operationLabel
			p.operationTokenEnum = enumTokenTypes.IDENTIFIER
			p.operationTokenValue = ""
			p.parentParserEnum = enumParserTypes.Label
			return nil
		}
		return errorHandler.AddNew(enumErrorCodes.OperationLabelBadTokenAfter, p.lookaheadValue) // ❌ Fails if tokens follow colon
	}

	//Potential labeled directive or assignment statement
	return p.getLabelFollowup(operationLabel, false)
}

// ====================================================================

// Labeled directive
func (p *OperationParser) getLabelFollowup(operationLabel string, hadWhitespace bool) error {
	if p.lookaheadType == enumTokenTypes.WHITESPACE {
		err := p.eatFreelyAndAdvance(enumTokenTypes.WHITESPACE)
		if err != nil {
			return err
		}

		if p.lookaheadType != enumTokenTypes.None {
			return p.getLabelFollowup(operationLabel, true)
		}
	}

	var operationTokenValue string
	var operationTokenEnum tokenEnum
	var parentParserEnum enumParserTypes.Def

	switch p.lookaheadType {

	//Equals sign
	case enumTokenTypes.ASSIGN_simple:
		operationTokenEnum = p.lookaheadType
		operationTokenValue = p.lookaheadValue
		parentParserEnum = enumParserTypes.Label
		err := p.eatFreelyAndAdvance(enumTokenTypes.ASSIGN_simple)
		if err != nil {
			return err
		}

	//EQU
	case enumTokenTypes.ASSIGN_EQU:
		operationTokenEnum = p.lookaheadType
		operationTokenValue = p.lookaheadValue
		parentParserEnum = enumParserTypes.Label
		err := p.eatFreelyAndAdvance(enumTokenTypes.ASSIGN_EQU)
		if err != nil {
			return err
		}

	//Actual directive
	case enumTokenTypes.DELIMITER_period:
		if !hadWhitespace {
			return errorHandler.AddNew(enumErrorCodes.OperationLabeledDirectiveNoSpace) // ❌ Fails
		}
		err := p.eatFreelyAndAdvance(enumTokenTypes.DELIMITER_period)
		if err != nil {
			return err
		}

		switch p.lookaheadType {
		case enumTokenTypes.DIRECTIVE_labeled,
			enumTokenTypes.DIRECTIVE_labeledBlockStart,
			enumTokenTypes.DIRECTIVE_labeledBlockEnd:
			operationTokenEnum = p.lookaheadType
			operationTokenValue = p.lookaheadValue
			parentParserEnum = enumParserTypes.Directive
			err := p.eatAndAdvance(p.lookaheadType)
			if err != nil {
				return err
			}

			if p.lookaheadType != enumTokenTypes.None && p.lookaheadType != enumTokenTypes.WHITESPACE {
				// ❌ Fails as unexpected token
				err := p.eat(enumTokenTypes.WHITESPACE)
				return err
			}

		default:
			return errorHandler.AddNew(enumErrorCodes.OperationLabeledDirectiveUnknown) // ❌ Fails

		}
	default:
		return errorHandler.AddNew(enumErrorCodes.OperationLabelBadTokenAfter, p.lookaheadValue) // ❌ Fails
	}

	// 🟢 Labeled directive parsing succeeds
	p.operationLabel = operationLabel
	p.operationTokenEnum = operationTokenEnum
	p.operationTokenValue = operationTokenValue
	p.parentParserEnum = parentParserEnum
	return nil
}
