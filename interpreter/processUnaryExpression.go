package interpreter

import (
	"misc/nintasm/assemble/errorHandler"
	enumErrorCodes "misc/nintasm/constants/enums/errorCodes"
	enumNodeTypes "misc/nintasm/constants/enums/nodeTypes"
	"misc/nintasm/interpreter/operandFactory"
)

func processUnaryExpression(node Node) (Node, error) {
	var err error

	operation := node.NodeValue
	right, err := EvaluateNode(*node.Right)
	if err != nil {
		return node, err
	}
	switch operation {
	case "+", "-", "~":
		if right.NodeType != enumNodeTypes.NumericLiteral {
			return node, errorHandler.AddNew(enumErrorCodes.InterpreterUnaryNotNumeric, operation, right.NodeValue)
		}
	case "!":
		if right.NodeType != enumNodeTypes.BooleanLiteral {
			return node, errorHandler.AddNew(enumErrorCodes.InterpreterUnaryNotBoolean, operation, right.NodeValue)
		}
	}

	node.Right = nil
	switch operation {
	case "+":
		node.AsNumber = right.AsNumber
	case "-":
		node.AsNumber = -right.AsNumber
	case "~":
		node.AsNumber = float64(^int(right.AsNumber))
	case "!":
		node.AsBool = !right.AsBool
	default:
		panic("🛑 SOMETHING IS VERY WRONG WITH UNARY EXPRESSION INTERPRETING")
	}

	switch operation {
	case "+", "-", "~":
		operandFactory.ConvertNodeToNumericLiteral(&node)
	case "!":
		operandFactory.ConvertNodeToBooleanLiteral(&node)
	}
	return node, nil
}
