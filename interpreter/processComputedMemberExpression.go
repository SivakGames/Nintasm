package interpreter

import (
	"misc/nintasm/assemble/errorHandler"
	enumErrorCodes "misc/nintasm/constants/enums/errorCodes"
	"misc/nintasm/interpreter/operandFactory"
)

func processComputedMemberExpression(node Node) (Node, error) {
	parentNode, err := EvaluateNode(*node.Left)
	if err != nil {
		return node, err
	}

	evaluatedArgList := make([]Node, len(*node.ArgumentList))
	for i, n := range *node.ArgumentList {
		evalN, err := EvaluateNode(n)
		if err != nil {
			return node, err
		}
		evaluatedArgList[i] = evalN
	}

	if !operandFactory.ValidateNodeIsArray(&parentNode) {
		return node, errorHandler.AddNew(enumErrorCodes.NodeTypeNotArray)
	}

	maxRange := len((*parentNode.ArgumentList)) - 1
	var result Node

	for _, n := range evaluatedArgList {
		index := int(n.AsNumber)
		if index < 0 {
			return node, errorHandler.AddNew(enumErrorCodes.InterpreterComputedMemberNegativeIndex, index)
		}
		if index > maxRange {
			return node, errorHandler.AddNew(enumErrorCodes.InterpreterComputedMemberIndexTooBig, index, maxRange)
		}
		result = (*parentNode.ArgumentList)[int(n.AsNumber)]
		parentNode = result
	}

	evaluatedComputedMemberNode, err := EvaluateNode(result)
	if err != nil {
		return node, err
	}
	return evaluatedComputedMemberNode, err
}
