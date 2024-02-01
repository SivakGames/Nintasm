package parser

import (
	"errors"
	"fmt"
	"misc/nintasm/parser/operandFactory"
	"misc/nintasm/tokenizer/tokenizerSpec"
	"strconv"
)

// General type for other operand parsers to borrow from
type OperandParser struct {
	operandLine     string
	operandPosition int
	Parser
}

// Used by instructions.  Will just get the first operand found (if any)
// Doing it this way will make parsing potential indexes earier
func (p *OperandParser) GetFirstOperandOnly() (operandFactory.Node, error) {
	var firstOperand operandFactory.Node = operandFactory.EmptyNode()

	//No operands at all
	if p.lookaheadType == tokenizerSpec.None {
		return firstOperand, nil // 🟢 Succeeds
	}
	//No commas at the beginning...
	if p.lookaheadType == tokenizerSpec.DELIMITER_comma {
		return firstOperand, errors.New("First operand cannot be a comma!") // ❌ Fails
	}

	firstOperand, err := p.operandStatementList()
	if err != nil {
		return firstOperand, err // ❌ Fails
	}

	return firstOperand, nil // 🟢 Succeeds
}

//=============================================
//=============================================
// Main parser for operands starts here
//=============================================
//=============================================

// Used by most operations.  Will go through and separate operands by commas returning each one as an AST
func (p *OperandParser) GetOperandList() ([]operandFactory.Node, error) {
	operandList := []operandFactory.Node{}
	operandCount := 0

	//No operands at all
	if p.lookaheadType == tokenizerSpec.None {
		return operandList, nil // 🟢 Succeeds
	}

	//No commas at the beginning...
	if p.lookaheadType == tokenizerSpec.DELIMITER_comma {
		return operandList, errors.New("Operand list cannot start with a comma!") // ❌ Fails
	}

	//There is at least one operand
	firstOperand, err := p.operandStatementList()
	if err != nil {
		return operandList, err // ❌ Fails
	}

	operandList = append(operandList, firstOperand)

	//From here, operands are comma separated
	for p.lookaheadType != tokenizerSpec.None && p.lookaheadType == tokenizerSpec.DELIMITER_comma {
		err = p.eatFreelyAndAdvance(tokenizerSpec.DELIMITER_comma)
		if err != nil {
			return operandList, err // ❌ Fails
		}
		data, err := p.operandStatementList()
		if err != nil {
			return operandList, err // ❌ Fails
		}
		operandList = append(operandList, data)
		operandCount++
	}

	return operandList, nil // 🟢 Succeeds
}

// =============================================
func (p *OperandParser) operandStatementList() (operandFactory.Node, error) {
	return p.statementList(tokenizerSpec.DELIMITER_comma)
}

// Get statements
func (p *OperandParser) statementList(stopTokenType tokenizerSpec.TokenType) (operandFactory.Node, error) {
	statementList, err := p.Statement()
	if err != nil {
		return statementList, err
	}

	//Subsequent operands
	for p.lookaheadType != tokenizerSpec.None && p.lookaheadType != stopTokenType {
		return p.Statement()
	}

	return statementList, nil
}

func (p *OperandParser) Statement() (operandFactory.Node, error) {
	if p.lookaheadType == tokenizerSpec.None {
		return operandFactory.EmptyNode(), nil
	}
	return p.logicalOrExpression()
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
func (p *OperandParser) logicalOrExpression() (operandFactory.Node, error) {
	return p._logicalExpression(p.logicalAndExpression, tokenizerSpec.OPERATOR_logicalOr)
}

// && expression
func (p *OperandParser) logicalAndExpression() (operandFactory.Node, error) {
	return p._logicalExpression(p.bitwiseOrExpression, tokenizerSpec.OPERATOR_logicalAnd)
}

// | expression
func (p *OperandParser) bitwiseOrExpression() (operandFactory.Node, error) {
	return p._logicalExpression(p.bitwiseXorExpression, tokenizerSpec.OPERATOR_bitwiseOr)
}

// ^ expression
func (p *OperandParser) bitwiseXorExpression() (operandFactory.Node, error) {
	return p._logicalExpression(p.bitwiseAndExpression, tokenizerSpec.OPERATOR_bitwiseXor)
}

// & expression
func (p *OperandParser) bitwiseAndExpression() (operandFactory.Node, error) {
	return p._logicalExpression(p.equalityExpression, tokenizerSpec.OPERATOR_bitwiseAnd)
}

// ==, != expression
func (p *OperandParser) equalityExpression() (operandFactory.Node, error) {
	return p._logicalExpression(p.shiftExpression, tokenizerSpec.OPERATOR_equality)
}

// <<, >> expression
func (p *OperandParser) shiftExpression() (operandFactory.Node, error) {
	return p._logicalExpression(p.relationalExpression, tokenizerSpec.OPERATOR_shift)
}

// <,<=,>=,> expression
func (p *OperandParser) relationalExpression() (operandFactory.Node, error) {
	return p._logicalExpression(p.additiveExpression, tokenizerSpec.OPERATOR_relational)
}

// +,- expression
func (p *OperandParser) additiveExpression() (operandFactory.Node, error) {
	return p._logicalExpression(p.multiplicativeExpression, tokenizerSpec.OPERATOR_additive)
}

// *,/,% expression
func (p *OperandParser) multiplicativeExpression() (operandFactory.Node, error) {
	return p._logicalExpression(p.unaryExpression, tokenizerSpec.OPERATOR_multiplicative)
}

//---------------------

// Preceding -,~,! expression
func (p *OperandParser) unaryExpression() (operandFactory.Node, error) {
	if p.lookaheadType != tokenizerSpec.None {
		switch p.lookaheadType {
		case tokenizerSpec.OPERATOR_additive,
			tokenizerSpec.OPERATOR_logicalNot,
			tokenizerSpec.OPERATOR_negate:
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
			return operandFactory.UnaryExpression(unaryType, unaryValue, argument), nil
		}
	}
	return p.callMemberExpression()
}

//---------------------

func (p *OperandParser) callMemberExpression() (operandFactory.Node, error) {
	callMemberValue := p.lookaheadValue
	callMemberType := p.lookaheadType

	member, err := p.memberExpression()
	if err != nil {
		return operandFactory.ErrorNode(p.lookaheadValue), err // ❌ Fails
	}

	if p.lookaheadType == tokenizerSpec.DELIMITER_leftParenthesis {
		if !_checkValidAssignmentTarget(callMemberType) {
			badCallee := fmt.Sprintf("Illegal functional callee name: %v", callMemberValue)
			return operandFactory.ErrorNode(p.lookaheadValue), errors.New(badCallee) // ❌ Fails
		}
		return p._callExpression(callMemberValue)
	}
	return member, nil
}

// Call expressions MUST begin with an identifier
func _checkValidAssignmentTarget(assignmentType tokenizerSpec.TokenType) bool {
	return (assignmentType == tokenizerSpec.IDENTIFIER)
}

// ---------------------
// Functions
func (p *OperandParser) _callExpression(callee string) (operandFactory.Node, error) {
	arguments, err := p.arguments()
	if err != nil {
		return operandFactory.ErrorNode(p.lookaheadValue), err // ❌ Fails
	}

	callExpr := operandFactory.CallExpression(callee, arguments)

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

func (p *OperandParser) arguments() ([]operandFactory.Node, error) {
	var argumentList []operandFactory.Node
	var err error = nil

	p.eatAndAdvance(tokenizerSpec.DELIMITER_leftParenthesis)
	if err != nil {
		return argumentList, err // ❌ Fails
	}

	if p.lookaheadType != tokenizerSpec.DELIMITER_rightParenthesis {
		argumentList, err = p.argumentList()
		if err != nil {
			return argumentList, err // ❌ Fails
		}
	}

	err = p.eatAndAdvance(tokenizerSpec.DELIMITER_rightParenthesis)
	if err != nil {
		return argumentList, err // ❌ Fails
	}
	return argumentList, nil
}

//--------------------

func (p *OperandParser) argumentList() ([]operandFactory.Node, error) {
	argumentList := []operandFactory.Node{}
	firstArgument, err := p.Statement()
	if err != nil {
		return argumentList, err // ❌ Fails
	}
	argumentList = append(argumentList, firstArgument)

	for p.lookaheadType == tokenizerSpec.DELIMITER_comma {
		err = p.eatAndAdvance(tokenizerSpec.DELIMITER_comma)
		if err != nil {
			return argumentList, err // ❌ Fails
		}
		nextArgument, err := p.Statement()
		if err != nil {
			return argumentList, err // ❌ Fails
		}
		argumentList = append(argumentList, nextArgument)
	}

	return argumentList, nil
}

//--------------------------
//Things such as dots before labels

func (p *OperandParser) memberExpression() (operandFactory.Node, error) {
	//The parent label
	parent := p.lookaheadValue

	result, err := p.primaryExpression()
	if err != nil {
		return operandFactory.ErrorNode(p.lookaheadValue), err // ❌ Fails
	}

	//If nothing else, just exit
	if p.lookaheadType == tokenizerSpec.None {
		return result, nil // 🟢 Succeeds
	}

	if p._isLiteral(p.lookaheadType) {
		return operandFactory.ErrorNode(p.lookaheadValue), errors.New("Misplaced literal") // ❌ Fails
	}

	if p.lookaheadType == tokenizerSpec.IDENTIFIER {
		return operandFactory.ErrorNode(p.lookaheadValue), errors.New("Misplaced identifier") // ❌ Fails
	}

	//A dot indicates member
	if p.lookaheadType == tokenizerSpec.DELIMITER_period {
		err = p.eatFreelyAndAdvance(tokenizerSpec.DELIMITER_period)
		if err != nil {
			return operandFactory.ErrorNode(p.lookaheadValue), err // ❌ Fails
		}
		if p.lookaheadType != tokenizerSpec.IDENTIFIER {
			if p.tokenizer.IsTokenIdentifierLike(p.lookaheadValue) {
				p.lookaheadType = tokenizerSpec.IDENTIFIER
			}
		}

		key := p.lookaheadValue
		_, err := p.identifier()
		if err != nil {
			return operandFactory.ErrorNode(p.lookaheadValue), err // ❌ Fails
		}
		return operandFactory.MemberExpression(parent, key, false), nil
	}

	return result, nil

}

// !!!!!!!!!!!!!!!!!!!!!!!
// Top of the food chain - highest precedence
func (p *OperandParser) primaryExpression() (operandFactory.Node, error) {
	if p.lookaheadType == tokenizerSpec.None {
		return operandFactory.ErrorNode(p.lookaheadValue), errors.New("THERE'S NO PRIMARY EXPR!!!") // ❌ Fails
	}
	if p._isLiteral(p.lookaheadType) {
		return p.literal()
	}

	switch p.lookaheadType {
	case tokenizerSpec.DELIMITER_leftParenthesis:
		return p.parenthesizedExpression()

	case tokenizerSpec.DELIMITER_period:
		// ⚠️ TODO: Add period checker for local labels
		return operandFactory.ErrorNode(p.lookaheadValue), errors.New("Period doesn't exist yet") // ❌ Fails

	case tokenizerSpec.IDENTIFIER:
		return p.identifier()
	}

	// ❌ Fails
	badPrimary := fmt.Sprintf("\x1b[31mBAD primary expr!!!\x1b[0m - \x1b[33m%v\x1b[0m", p.lookaheadValue)
	return operandFactory.ErrorNode(p.lookaheadValue), errors.New(badPrimary)
}

// ++++++++++++++++++++++++++++
// Helper for logical expressions
func (p *OperandParser) _logicalExpression(builderName func() (operandFactory.Node, error), operatorToken tokenizerSpec.TokenType) (operandFactory.Node, error) {
	var left operandFactory.Node
	var right operandFactory.Node
	var err error = nil
	left, err = builderName()
	if err != nil {
		return operandFactory.ErrorNode(p.lookaheadValue), err // ❌ Fails
	}

	for p.lookaheadType != tokenizerSpec.None && p.lookaheadType == operatorToken {
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
		left = operandFactory.LogicalExpression(logicalExpressionType, logicalExpressionValue, left, right)
	}

	return left, nil
}

// ++++++++++++++++++++++++++++++
// Helper to see if value is one of the literal types
func (p *OperandParser) _isLiteral(tokenType tokenizerSpec.TokenType) bool {
	switch tokenType {
	case tokenizerSpec.NUMBER_binary,
		tokenizerSpec.NUMBER_decimal,
		tokenizerSpec.NUMBER_hex,
		tokenizerSpec.STRING,
		tokenizerSpec.BACKTICK_STRING,
		tokenizerSpec.SUBSTITUTION_numMacroArgs,
		tokenizerSpec.SUBSTITUTION_stringID,
		tokenizerSpec.SUBSTITUTION_numericID:
		return true
	}
	return false
}

// ((((((((((((((((
func (p *OperandParser) parenthesizedExpression() (operandFactory.Node, error) {
	err := p.eatFreelyAndAdvance(tokenizerSpec.DELIMITER_leftParenthesis)
	if err != nil {
		return operandFactory.ErrorNode(p.lookaheadValue), err // ❌ Fails
	}

	expression, err := p.Statement()
	if err != nil {
		return operandFactory.ErrorNode(p.lookaheadValue), err // ❌ Fails
	}

	err = p.eatAndAdvance(tokenizerSpec.DELIMITER_rightParenthesis)
	if err != nil {
		return operandFactory.ErrorNode(p.lookaheadValue), err // ❌ Fails
	}

	return expression, nil
}

// -----------------
func (p *OperandParser) literal() (operandFactory.Node, error) {
	literalType := p.lookaheadType
	literalValue := p.lookaheadValue
	err := p.eatFreelyAndAdvance(p.lookaheadType)
	if err != nil {
		return operandFactory.EmptyNode(), err
	}

	switch literalType {
	case tokenizerSpec.NUMBER_hex:
		asInt, _ := strconv.ParseInt(literalValue[1:], 16, 64)
		literalValue := fmt.Sprintf("%d", asInt)
		return operandFactory.NumericLiteral(literalType, literalValue), nil
	case tokenizerSpec.NUMBER_binary:
		asInt, _ := strconv.ParseInt(literalValue[1:], 2, 64)
		literalValue := fmt.Sprintf("%d", asInt)
		return operandFactory.NumericLiteral(literalType, literalValue), nil
	case tokenizerSpec.NUMBER_decimal:
		return operandFactory.NumericLiteral(literalType, literalValue), nil
	case tokenizerSpec.STRING:
		return operandFactory.StringLiteral(literalType, literalValue), nil
	case tokenizerSpec.BACKTICK_STRING:
		return operandFactory.BacktickStringLiteral(literalType, literalValue), nil
	case tokenizerSpec.SUBSTITUTION_numericID:
		return operandFactory.SubstitutionId(literalType, literalValue), nil
	case tokenizerSpec.SUBSTITUTION_stringID:
		return operandFactory.SubstitutionId(literalType, literalValue), nil
	case tokenizerSpec.SUBSTITUTION_numMacroArgs:
		return operandFactory.SubstitutionId(literalType, literalValue), nil
	}
	// ❌ Fails
	return operandFactory.ErrorNode(p.lookaheadValue), errors.New("\x1b[31mERROR!!!!!!!!\x1b[0m")
}

// -----------------
func (p *OperandParser) identifier() (operandFactory.Node, error) {
	literalType := p.lookaheadType
	literalValue := p.lookaheadValue
	err := p.eatFreelyAndAdvance(tokenizerSpec.IDENTIFIER)
	if err != nil {
		return operandFactory.ErrorNode(p.lookaheadValue), err // ❌ Fails
	}

	return operandFactory.Identifier(literalType, literalValue), nil
}
