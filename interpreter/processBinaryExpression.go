package interpreter

import (
	"misc/nintasm/assemble/errorHandler"
	enumErrorCodes "misc/nintasm/constants/enums/errorCodes"
	"misc/nintasm/interpreter/operandFactory"
)

func processBinaryExpression(node Node) (Node, error) {
	var err error
	operation := node.NodeValue
	left, err := EvaluateNode(*node.Left)
	if err != nil {
		return node, err
	}

	right, err := EvaluateNode(*node.Right)
	if err != nil {
		return node, err
	}

	if left.NodeType != right.NodeType {
		badLeftValue := node.Left.NodeValue
		badRightValue := node.Right.NodeValue
		operandFactory.ConvertNodeToNumericLiteral(&node)
		return node, errorHandler.AddNew(
			enumErrorCodes.InterpreterBinaryMismatchedTypes,
			badLeftValue, operation, badRightValue,
		)
	}

	node.Left = nil
	node.Right = nil

	switch operation {
	case "+":
		node.AsNumber = left.AsNumber + right.AsNumber
	case "-":
		node.AsNumber = left.AsNumber - right.AsNumber
	case "*":
		node.AsNumber = left.AsNumber * right.AsNumber
	case "/":
		node.AsNumber = left.AsNumber / right.AsNumber
	case "%":
		node.AsNumber = float64(int(left.AsNumber) % int(right.AsNumber))
	case "|":
		node.AsNumber = float64(int(left.AsNumber) | int(right.AsNumber))
	case "&":
		node.AsNumber = float64(int(left.AsNumber) & int(right.AsNumber))
	case "^":
		node.AsNumber = float64(int(left.AsNumber) ^ int(right.AsNumber))
	case "<<":
		node.AsNumber = float64(int(left.AsNumber) << int(right.AsNumber))
	case ">>":
		node.AsNumber = float64(int(left.AsNumber) >> int(right.AsNumber))
	case "<":
		node.AsBool = left.AsNumber < right.AsNumber
	case "<=":
		node.AsBool = left.AsNumber <= right.AsNumber
	case ">":
		node.AsBool = left.AsNumber > right.AsNumber
	case ">=":
		node.AsBool = left.AsNumber > right.AsNumber
	case "==":
		node.AsBool = left.AsNumber == right.AsNumber
	case "!=":
		node.AsBool = left.AsNumber != right.AsNumber
	case "&&":
		node.AsBool = left.AsBool && right.AsBool
	case "||":
		node.AsBool = left.AsBool || right.AsBool
	default:
		panic("🛑 SOMETHING IS VERY WRONG WITH BINARY EXPRESSION INTERPRETING")
	}
	switch operation {
	case "+", "-", "*", "/", "%", "|", "&", "^", "<<", ">>":
		operandFactory.ConvertNodeToNumericLiteral(&node)
	case "<", "<=", ">", ">=", "==", "!=", "&&", "||":
		operandFactory.ConvertNodeToBooleanLiteral(&node)
	}
	return node, nil
}
