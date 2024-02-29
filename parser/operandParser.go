package parser

import (
	"misc/nintasm/assemble/errorHandler"
	enumErrorCodes "misc/nintasm/constants/enums/errorCodes"
	enumInstructionModes "misc/nintasm/constants/enums/instructionModes"
	enumTokenTypes "misc/nintasm/constants/enums/tokenTypes"
	"misc/nintasm/interpreter"
	"misc/nintasm/interpreter/operandFactory"
	"strconv"
)

type instModeEnum = enumInstructionModes.Def
type Node = operandFactory.Node

// General type for other operand parsers to borrow from
type OperandParser struct {
	operandListStringStartPosition int
	operandStartPosition           int
	instructionMode                instModeEnum
	instructionXYIndex             tokenEnum
	manuallyEvalOperands           bool
	Parser
}

//=============================================
//=============================================
// Main parser for operands starts here
//=============================================
//=============================================

// Used by most operations.  Will go through and separate operands by commas returning each one as an AST
func (p *OperandParser) GetOperandList(
	minOperands int, maxOperands int,
	manuallyEvalOperands bool, captureMasks []string) ([]Node, error) {

	operandList := []Node{}
	operandCount := 0
	p.manuallyEvalOperands = manuallyEvalOperands

	//See if there are no operands at all
	if p.lookaheadType != enumTokenTypes.None {

		//No commas allowed at the beginning...
		if p.lookaheadType == enumTokenTypes.DELIMITER_comma {
			return operandList, errorHandler.AddNew(enumErrorCodes.OperandListStartingComma) // ❌ Fails
		}

		// Get first operand
		err := p.getOperandAndAppend(&operandList, &captureMasks)
		if err != nil {
			return operandList, err // ❌ Fails
		}

		//From here get subsequent operands, if any. Operands are comma-separated
		for p.lookaheadType != enumTokenTypes.None && p.lookaheadType == enumTokenTypes.DELIMITER_comma {
			err := p.eatFreelyAndAdvance(enumTokenTypes.DELIMITER_comma)
			if err != nil {
				return operandList, err // ❌ Fails
			}
			err = p.getOperandAndAppend(&operandList, &captureMasks)
			if err != nil {
				return operandList, err // ❌ Fails
			}
			operandCount++
		}
	}

	//Check if too many or too few operands...

	if len(operandList) > maxOperands {
		return operandList, errorHandler.AddNew(enumErrorCodes.OperandListTooMany, maxOperands) // ❌ Fails
	} else if len(operandList) < minOperands {
		return operandList, errorHandler.AddNew(enumErrorCodes.OperandListTooFew, minOperands) // ❌ Fails
	}

	return operandList, nil // 🟢 Succeeds
}

// =============================================
// =============================================

// Get the actual operand. What type of operand being captured dictates the followup function
func (p *OperandParser) getOperandAndAppend(operandList *[]Node, captureMasks *[]string) error {
	captureStatementFunction := p.statement

	if len(*captureMasks) > len(*operandList) {
		switch (*captureMasks)[len(*operandList)] {
		case "instruction":
			captureStatementFunction = p.instructionPrefix
		case "macro":
			captureStatementFunction = p.macroReplaceStatement
		default:
			captureStatementFunction = p.statement
		}
	}

	operand, err := captureStatementFunction()
	if err != nil {
		err := errorHandler.CheckErrorContinuesUpwardPropagation(err, enumErrorCodes.Error)
		if err != nil {
			return err // ❌❌ CONTINUES Failing!
		}
	}
	*operandList = append(*operandList, operand)
	return nil
}

// =============================================
// Instructions only
// =============================================

func (p *OperandParser) instructionPrefix() (Node, error) {
	p.instructionMode = enumInstructionModes.ABS
	p.instructionXYIndex = enumTokenTypes.None
	xyIndex := enumTokenTypes.None
	checkXYfollowup := false
	var statement Node
	var err error

	switch p.lookaheadType {

	//[][][][][][][][][][][][][][][][][][][][][]
	//Indirect

	case enumTokenTypes.DELIMITER_leftSquareBracket:
		p.instructionMode = enumInstructionModes.IND
		err = p.eatFreelyAndAdvance(enumTokenTypes.DELIMITER_leftSquareBracket)
		if err != nil {
			return operandFactory.ErrorNode(p.lookaheadValue), err // ❌ Fails
		}
		statement, err = p.statement()
		if err != nil {
			return statement, err // ❌ Fails
		}
		// For indirect X
		if p.lookaheadType == enumTokenTypes.DELIMITER_comma {
			xyIndex, err = p.checkInstructionXYIndex()
			if err != nil {
				return statement, err // ❌ Fails
			}
			if xyIndex != enumTokenTypes.REGISTER_X {
				return statement, errorHandler.AddNew(enumErrorCodes.InstIndirectIndexMustBeX) // ❌ Fails
			}
			err = p.eatFreelyAndAdvance(enumTokenTypes.DELIMITER_rightSquareBracket)
			if err != nil {
				return statement, err // ❌ Fails
			}
		} else {
			err = p.eatFreelyAndAdvance(enumTokenTypes.DELIMITER_rightSquareBracket)
			if err != nil {
				return statement, err // ❌ Fails
			}
			if p.lookaheadType == enumTokenTypes.DELIMITER_comma {
				xyIndex, err = p.checkInstructionXYIndex()
				if err != nil {
					return statement, err // ❌ Fails
				}
				if xyIndex != enumTokenTypes.REGISTER_Y {
					return statement, errorHandler.AddNew(enumErrorCodes.InstIndirectIndexMustBeY) // ❌ Fails
				}
			}
		}

	//######################################
	//Immediate mode

	case enumTokenTypes.DELIMITER_hash:
		p.instructionMode = enumInstructionModes.IMM
		err = p.eatFreelyAndAdvance(enumTokenTypes.DELIMITER_hash)
		if err != nil {
			return operandFactory.ErrorNode(p.lookaheadValue), err // ❌ Fails
		}
		statement, err = p.statement()
		if err != nil {
			return operandFactory.ErrorNode(p.lookaheadValue), err // ❌ Fails
		}
		checkXYfollowup = true

	//<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<
	//Explicit ZP mode

	case enumTokenTypes.OPERATOR_relational:
		if p.lookaheadValue == "<" {
			p.instructionMode = enumInstructionModes.ZP
			err = p.eatFreelyAndAdvance(enumTokenTypes.OPERATOR_relational)
			if err != nil {
				return operandFactory.ErrorNode(p.lookaheadValue), err // ❌ Fails
			}
		}
		statement, err = p.statement()
		if err != nil {
			return operandFactory.ErrorNode(p.lookaheadValue), err // ❌ Fails
		}
		checkXYfollowup = true

	//-------------------------------------
	//Absolute mode

	default:
		statement, err = p.statement()
		if err != nil {
			return operandFactory.ErrorNode(p.lookaheadValue), err // ❌ Fails
		}
		checkXYfollowup = true
	}

	//-------------------------------------
	//-------------------------------------

	if checkXYfollowup && p.lookaheadType == enumTokenTypes.DELIMITER_comma {
		xyIndex, err = p.checkInstructionXYIndex()
		if err != nil {
			return operandFactory.ErrorNode(p.lookaheadValue), err // ❌ Fails
		}
	}
	p.instructionXYIndex = xyIndex

	if p.lookaheadType != enumTokenTypes.None {
		return statement, errorHandler.AddNew(enumErrorCodes.InstTokenAfterOperand, p.lookaheadValue)
	}

	return statement, nil // 🟢 Succeeds
}

// ++++++++++++++++++++++++++++++++

func (p *OperandParser) checkInstructionXYIndex() (tokenEnum, error) {
	err := p.eatFreelyAndAdvance(enumTokenTypes.DELIMITER_comma)
	targetIndex := enumTokenTypes.None
	if err != nil {
		return targetIndex, err // ❌ Fails
	}

	if p.lookaheadType != enumTokenTypes.REGISTER_X && p.lookaheadType != enumTokenTypes.REGISTER_Y {
		return targetIndex, errorHandler.AddNew(enumErrorCodes.InstBadIndexValue, p.lookaheadValue)
	}

	targetIndex = p.lookaheadType
	err = p.eatFreelyAndAdvance(p.lookaheadType)
	if err != nil {
		return targetIndex, err // ❌ Fails
	}

	return targetIndex, nil
}

// =============================================
/* func (p *OperandParser) operandStatementList() (Node, error) {
	return p.statementList(enumTokenTypes.DELIMITER_comma)
}
*/

// =============================================
// Get statements - used for things like function arguments
/*
func (p *OperandParser) statementList(stopTokenType tokenEnum) (Node, error) {
	statementList, err := p.statement()
	if err != nil {
		return statementList, err
	}

	statementList, err = interpreter.EvaluateNode(statementList)
	if err != nil {
		return statementList, err
	}
	//If somehow, after evaluation, the next token is not the stop token...
	for p.lookaheadType != enumTokenTypes.None && p.lookaheadType != stopTokenType {
		log.Println("\x1b[38;5;202mEvaluating next statement...\x1b[0m")
		return p.statement()
	}

	return statementList, nil
} */

// =============================================

// Will return a special node recognized by macros that stops at commas.
// Statements within curly braces will ignore commas
func (p *OperandParser) macroReplaceStatement() (Node, error) {
	var closingTokenEnum []enumTokenTypes.Def

	replacement := ""
	closingTokenEnum = append(closingTokenEnum, enumTokenTypes.DELIMITER_comma)

	for len(closingTokenEnum) > 0 && p.lookaheadType != enumTokenTypes.None {
		topOfStackEnum := closingTokenEnum[len(closingTokenEnum)-1]
		switch p.lookaheadType {
		case topOfStackEnum:
			closingTokenEnum = closingTokenEnum[:len(closingTokenEnum)-1]
			if len(closingTokenEnum) > 0 {
				err := p.eatFreelyAndAdvance(topOfStackEnum)
				if err != nil {
					return operandFactory.ErrorNode(p.lookaheadValue), err // ❌ Fails
				}
			}

		case enumTokenTypes.DELIMITER_leftCurlyBrace:
			if topOfStackEnum == enumTokenTypes.DELIMITER_rightCurlyBrace {
				return operandFactory.ErrorNode(p.lookaheadValue),
					errorHandler.AddNew(enumErrorCodes.MacroInvokeDoubleCurlyBrace) // ❌ Fails
			}
			err := p.eatFreelyAndAdvance(enumTokenTypes.DELIMITER_leftCurlyBrace)
			if err != nil {
				return operandFactory.ErrorNode(p.lookaheadValue), err // ❌ Fails
			}
			closingTokenEnum = append(closingTokenEnum, enumTokenTypes.DELIMITER_rightCurlyBrace)

		default:
			replacement += p.lookaheadValue
			err := p.eatFreelyAndAdvance(p.lookaheadType)
			if err != nil {
				return operandFactory.ErrorNode(p.lookaheadValue), err // ❌ Fails
			}
		}
	}
	if len(closingTokenEnum) > 1 {
		return operandFactory.ErrorNode(p.lookaheadValue),
			errorHandler.AddNew(enumErrorCodes.MacroInvokeUnclosedCurlyBrace) // ❌ Fails
	}

	return operandFactory.CreateMacroReplacementNode(replacement), nil
}

// =============================================

func (p *OperandParser) statement() (Node, error) {
	if p.lookaheadType == enumTokenTypes.None {
		return operandFactory.EmptyNode(), nil
	}
	statement, err := p.bitwiseOrExpression()
	if err != nil {
		return statement, err
	}
	if !p.manuallyEvalOperands {
		statement, err = interpreter.EvaluateNode(statement)
	}
	return statement, err
}

/*
------------------------------------------
❕❕❕❕❕❕❕❕❕❕❕❕❕❕❕❕❕❕❕❕❕❕❕❕❕❕❕❕❕❕❕❕❕❕❕❕❕❕❕❕❕❕❕❕❕❕
	From this point, things finally start getting
	evaluated. This point will check
	if there is a left hand side operator and trickle
	through the various operations in order of precedence
------------------------------------------
------------------------------------------
*/

/*
	Here we have the food chain of expressions in order of
	lowest predence to highest
*/

// || expression
func (p *OperandParser) logicalOrExpression() (Node, error) {
	return p._logicalExpression(p.logicalAndExpression, enumTokenTypes.OPERATOR_logicalOr)
}

// && expression
func (p *OperandParser) logicalAndExpression() (Node, error) {
	return p._logicalExpression(p.bitwiseOrExpression, enumTokenTypes.OPERATOR_logicalAnd)
}

// | expression
func (p *OperandParser) bitwiseOrExpression() (Node, error) {
	return p._logicalExpression(p.bitwiseXorExpression, enumTokenTypes.OPERATOR_bitwiseOr)
}

// ^ expression
func (p *OperandParser) bitwiseXorExpression() (Node, error) {
	return p._logicalExpression(p.bitwiseAndExpression, enumTokenTypes.OPERATOR_bitwiseXor)
}

// & expression
func (p *OperandParser) bitwiseAndExpression() (Node, error) {
	return p._logicalExpression(p.equalityExpression, enumTokenTypes.OPERATOR_bitwiseAnd)
}

// ==, != expression
func (p *OperandParser) equalityExpression() (Node, error) {
	return p._logicalExpression(p.shiftExpression, enumTokenTypes.OPERATOR_equality)
}

// <<, >> expression
func (p *OperandParser) shiftExpression() (Node, error) {
	return p._logicalExpression(p.relationalExpression, enumTokenTypes.OPERATOR_shift)
}

// <,<=,>=,> expression
func (p *OperandParser) relationalExpression() (Node, error) {
	return p._logicalExpression(p.additiveExpression, enumTokenTypes.OPERATOR_relational)
}

// +,- expression
func (p *OperandParser) additiveExpression() (Node, error) {
	return p._logicalExpression(p.multiplicativeExpression, enumTokenTypes.OPERATOR_additive)
}

// *,/,% expression
func (p *OperandParser) multiplicativeExpression() (Node, error) {
	return p._logicalExpression(p.unaryExpression, enumTokenTypes.OPERATOR_multiplicative)
}

//---------------------

// Preceding -,~,! expression
func (p *OperandParser) unaryExpression() (Node, error) {
	if p.lookaheadType != enumTokenTypes.None {
		switch p.lookaheadType {
		case enumTokenTypes.OPERATOR_additive,
			enumTokenTypes.OPERATOR_logicalNot,
			enumTokenTypes.OPERATOR_negate:
			unaryType := p.lookaheadType
			unaryValue := p.lookaheadValue
			err := p.eatFreelyAndAdvance(p.lookaheadType)
			if err != nil {
				return operandFactory.ErrorNode(p.lookaheadValue), err // ❌ Fails
			}
			argument, err := p.unaryExpression()
			if err != nil {
				return operandFactory.ErrorNode(p.lookaheadValue), err // ❌ Fails
			}
			return operandFactory.CreateUnaryExpressionNode(unaryType, unaryValue, argument), nil
		}
	}
	return p.callMemberExpression()
}

//---------------------

func (p *OperandParser) callMemberExpression() (Node, error) {
	callMemberValue := p.lookaheadValue
	callMemberType := p.lookaheadType

	member, err := p.memberExpression()
	if err != nil {
		return operandFactory.ErrorNode(p.lookaheadValue), err // ❌ Fails
	}

	if p.lookaheadType == enumTokenTypes.DELIMITER_leftParenthesis {
		if !_checkValidAssignmentTarget(callMemberType) {
			return operandFactory.ErrorNode(p.lookaheadValue),
				errorHandler.AddNew(enumErrorCodes.OperandBadCalleeName, callMemberValue) // ❌ Fails
		}
		return p._callExpression(callMemberValue)
	}
	return member, nil
}

// +++++++++++++++++++++++

// Call expressions MUST begin with an identifier
func _checkValidAssignmentTarget(assignmentType tokenEnum) bool {
	return (assignmentType == enumTokenTypes.IDENTIFIER)
}

// ---------------------
// Functions
func (p *OperandParser) _callExpression(callee string) (Node, error) {
	var err error = nil

	if callee == "bank" {
		err = p.eatFreelyAndAdvance(enumTokenTypes.DELIMITER_leftParenthesis)
		if err != nil {
			return operandFactory.ErrorNode(p.lookaheadValue), err // ❌ Fails
		}
		bankArgument := operandFactory.CreateIdentifierNode(p.lookaheadType, p.lookaheadValue)
		err = p.eatAndAdvance(enumTokenTypes.IDENTIFIER)
		if err != nil {
			return operandFactory.ErrorNode(p.lookaheadValue), err // ❌ Fails
		}
		err = p.eatAndAdvance(enumTokenTypes.DELIMITER_rightParenthesis)
		if err != nil {
			return operandFactory.ErrorNode(p.lookaheadValue), err // ❌ Fails
		}

		return operandFactory.CreateCallExpressionNode(callee, []Node{bankArgument}), nil
	}

	arguments, err := p.arguments()
	if err != nil {
		return operandFactory.ErrorNode(p.lookaheadValue), err // ❌ Fails
	}

	callExpr := operandFactory.CreateCallExpressionNode(callee, arguments)

	// ⚠️ TODO: I forgot why this was in the original assembler...
	/* if p.lookaheadType == tokenizerSpec.DELIMITER_leftParenthesis {
		newCallee := callExpr.NodeValue
		newExpr, err := p._callExpression(newCallee)
		if err != nil {
			return operandFactory.EmptyNode(), err
		}
		callExpr = newExpr
	} */

	return callExpr, nil
}

//--------------------

func (p *OperandParser) arguments() ([]Node, error) {
	var argumentList []Node
	var err error = nil

	err = p.eatAndAdvance(enumTokenTypes.DELIMITER_leftParenthesis)
	if err != nil {
		return argumentList, err // ❌ Fails
	}

	if p.lookaheadType != enumTokenTypes.DELIMITER_rightParenthesis {
		argumentList, err = p.argumentList()
		if err != nil {
			return argumentList, err // ❌ Fails
		}
	}

	err = p.eatAndAdvance(enumTokenTypes.DELIMITER_rightParenthesis)
	if err != nil {
		return argumentList, err // ❌ Fails
	}
	return argumentList, nil
}

//--------------------

func (p *OperandParser) argumentList() ([]Node, error) {
	argumentList := []Node{}
	firstArgument, err := p.statement()
	if err != nil {
		return argumentList, err // ❌ Fails
	}
	argumentList = append(argumentList, firstArgument)

	for p.lookaheadType == enumTokenTypes.DELIMITER_comma {
		err = p.eatAndAdvance(enumTokenTypes.DELIMITER_comma)
		if err != nil {
			return argumentList, err // ❌ Fails
		}
		nextArgument, err := p.statement()
		if err != nil {
			return argumentList, err // ❌ Fails
		}
		argumentList = append(argumentList, nextArgument)
	}

	return argumentList, nil
}

//--------------------------
//Things such as dots before labels

func (p *OperandParser) memberExpression() (Node, error) {
	//The parent label
	parent := p.lookaheadValue

	result, err := p.primaryExpression()
	if err != nil {
		return operandFactory.ErrorNode(p.lookaheadValue), err // ❌ Fails
	}

	//If nothing else, just exit
	if p.lookaheadType == enumTokenTypes.None {
		return result, nil // 🟢 Succeeds
	}

	if p._isLiteral(p.lookaheadType) {
		return operandFactory.ErrorNode(p.lookaheadValue),
			errorHandler.AddNew(enumErrorCodes.OperandMisplacedLiteral, p.lookaheadValue) // ❌ Fails
	}

	if p.lookaheadType == enumTokenTypes.IDENTIFIER {
		return operandFactory.ErrorNode(p.lookaheadValue),
			errorHandler.AddNew(enumErrorCodes.OperandMisplacedIdentifier, p.lookaheadValue) // ❌ Fails
	}

	//A dot indicates member
	if p.lookaheadType == enumTokenTypes.DELIMITER_period {
		err = p.eatFreelyAndAdvance(enumTokenTypes.DELIMITER_period)
		if err != nil {
			return operandFactory.ErrorNode(p.lookaheadValue), err // ❌ Fails
		}
		if p.lookaheadType != enumTokenTypes.IDENTIFIER {
			if p.tokenizer.IsTokenIdentifierLike(p.lookaheadValue) {
				p.lookaheadType = enumTokenTypes.IDENTIFIER
			}
		}

		key := p.lookaheadValue
		_, err := p.identifier()
		if err != nil {
			return operandFactory.ErrorNode(p.lookaheadValue), err // ❌ Fails
		}
		return operandFactory.CreateMemberExpressionNode(parent, key, false), nil
	}

	return result, nil

}

// !!!!!!!!!!!!!!!!!!!!!!!
// Top of the food chain - highest precedence
func (p *OperandParser) primaryExpression() (Node, error) {
	if p.lookaheadType == enumTokenTypes.None {
		return operandFactory.ErrorNode(p.lookaheadValue),
			errorHandler.AddNew(enumErrorCodes.OperandMissingPrimaryExpr) // ❌ Fails
	}
	if p._isLiteral(p.lookaheadType) {
		return p.literal()
	}

	switch p.lookaheadType {
	case enumTokenTypes.DELIMITER_leftParenthesis:
		return p.parenthesizedExpression()

	//Dots will prepend the parent label
	case enumTokenTypes.DELIMITER_period:
		parentLabel, err := interpreter.GetParentLabel()
		if err != nil {
			return operandFactory.ErrorNode(p.lookaheadValue), err // ❌ Fails
		}
		err = p.eatFreelyAndAdvance(enumTokenTypes.DELIMITER_period)
		if err != nil {
			return operandFactory.ErrorNode(p.lookaheadValue), err // ❌ Fails
		}
		if p.lookaheadType == enumTokenTypes.IDENTIFIER {
			p.lookaheadValue = parentLabel + "." + p.lookaheadValue
			return p.identifier()
		}
		return operandFactory.ErrorNode(p.lookaheadValue),
			errorHandler.AddNew(enumErrorCodes.OperandPeriodMissingIdentifier) // ❌ Fails

	case enumTokenTypes.IDENTIFIER:
		return p.identifier()
	}

	// ❌ Fails
	return operandFactory.ErrorNode(p.lookaheadValue),
		errorHandler.AddNew(enumErrorCodes.OperandBadPrimaryExpr, p.lookaheadValue)
}

// ++++++++++++++++++++++++++++
// Helper for logical expressions
func (p *OperandParser) _logicalExpression(builderName func() (Node, error), operatorToken tokenEnum) (Node, error) {
	var left Node
	var right Node
	var err error = nil
	left, err = builderName()
	if err != nil {
		return operandFactory.ErrorNode(p.lookaheadValue), err // ❌ Fails
	}

	for p.lookaheadType != enumTokenTypes.None && p.lookaheadType == operatorToken {
		logicalExpressionType := p.lookaheadType
		logicalExpressionValue := p.lookaheadValue
		err = p.eatFreelyAndAdvance(operatorToken)
		if err != nil {
			return operandFactory.ErrorNode(p.lookaheadValue), err // ❌ Fails
		}
		right, err = builderName()
		if err != nil {
			return operandFactory.ErrorNode(p.lookaheadValue), err // ❌ Fails
		}
		left = operandFactory.CreateBinaryExpressionNode(logicalExpressionType, logicalExpressionValue, left, right)
	}

	return left, nil
}

// ++++++++++++++++++++++++++++++
// Helper to see if value is one of the literal types
func (p *OperandParser) _isLiteral(tokenType tokenEnum) bool {
	switch tokenType {
	case enumTokenTypes.NUMBER_binary,
		enumTokenTypes.NUMBER_decimal,
		enumTokenTypes.NUMBER_hex,
		enumTokenTypes.STRING,
		enumTokenTypes.BACKTICK_STRING,
		enumTokenTypes.SUBSTITUTION_numMacroArgs,
		enumTokenTypes.SUBSTITUTION_stringID,
		enumTokenTypes.SUBSTITUTION_numericID:
		return true
	}
	return false
}

// ((((((((((((((((
func (p *OperandParser) parenthesizedExpression() (Node, error) {
	err := p.eatFreelyAndAdvance(enumTokenTypes.DELIMITER_leftParenthesis)
	if err != nil {
		return operandFactory.ErrorNode(p.lookaheadValue), err // ❌ Fails
	}

	expression, err := p.statement()
	if err != nil {
		return operandFactory.ErrorNode(p.lookaheadValue), err // ❌ Fails
	}

	err = p.eatAndAdvance(enumTokenTypes.DELIMITER_rightParenthesis)
	if err != nil {
		return operandFactory.ErrorNode(p.lookaheadValue), err // ❌ Fails
	}

	return expression, nil
}

// -----------------
func (p *OperandParser) literal() (Node, error) {
	literalType := p.lookaheadType
	literalValue := p.lookaheadValue

	err := p.eatFreelyAndAdvance(p.lookaheadType)
	if err != nil {
		return operandFactory.EmptyNode(), err
	}

	switch literalType {
	case enumTokenTypes.NUMBER_hex:
		asNumber, _ := strconv.ParseInt(literalValue[1:], 16, 64)
		return operandFactory.CreateNumericLiteralNode(literalType, literalValue, int(asNumber)), nil
	case enumTokenTypes.NUMBER_binary:
		asNumber, _ := strconv.ParseInt(literalValue[1:], 2, 64)
		return operandFactory.CreateNumericLiteralNode(literalType, literalValue, int(asNumber)), nil
	case enumTokenTypes.NUMBER_decimal:
		asNumber, _ := strconv.ParseInt(literalValue, 10, 64)
		return operandFactory.CreateNumericLiteralNode(literalType, literalValue, int(asNumber)), nil
	case enumTokenTypes.STRING:
		return operandFactory.CreateStringLiteralNode(literalType, literalValue), nil
	case enumTokenTypes.BACKTICK_STRING:
		return operandFactory.CreateBacktickStringLiteralNode(literalType, literalValue), nil
	case enumTokenTypes.SUBSTITUTION_numericID:
		return operandFactory.CreateSubstitutionIdNode(literalType, literalValue), nil
	case enumTokenTypes.SUBSTITUTION_stringID:
		return operandFactory.CreateSubstitutionIdNode(literalType, literalValue), nil
	case enumTokenTypes.SUBSTITUTION_numMacroArgs:
		return operandFactory.CreateSubstitutionIdNode(literalType, literalValue), nil
	}
	// ❌ Fails
	panic("Something is greatly wrong with literal type")
}

// -----------------
func (p *OperandParser) identifier() (Node, error) {
	literalType := p.lookaheadType
	literalValue := p.lookaheadValue
	err := p.eatFreelyAndAdvance(enumTokenTypes.IDENTIFIER)
	if err != nil {
		return operandFactory.ErrorNode(p.lookaheadValue), err // ❌ Fails
	}

	return operandFactory.CreateIdentifierNode(literalType, literalValue), nil
}
