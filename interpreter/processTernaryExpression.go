package interpreter

import (
	"misc/nintasm/assemble/errorHandler"
	enumErrorCodes "misc/nintasm/constants/enums/errorCodes"
	"misc/nintasm/interpreter/operandFactory"
)

func processTernaryExpression(node Node) (Node, error) {
	conditionNode := (*node.ArgumentList)[0]

	conditionResult, err := EvaluateNode(conditionNode)
	if err != nil {
		return node, err
	}
	if !operandFactory.ValidateNodeIsBoolean(&conditionResult) {
		return node, errorHandler.AddNew(enumErrorCodes.NodeTypeNotBool)
	}

	var resultNode *Node
	if conditionResult.AsBool {
		resultNode = node.Left
	} else {
		resultNode = node.Right
	}
	result, err := EvaluateNode(*resultNode)
	if err != nil {
		return node, err
	}
	return result, nil
}
