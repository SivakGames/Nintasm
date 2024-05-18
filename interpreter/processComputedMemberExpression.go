package interpreter

import (
	"misc/nintasm/assemble/errorHandler"
	enumErrorCodes "misc/nintasm/constants/enums/errorCodes"
	enumNodeTypes "misc/nintasm/constants/enums/nodeTypes"
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

	var result Node
	for _, n := range evaluatedArgList {
		index := int(n.AsNumber)
		if parentNode.NodeType != enumNodeTypes.MultiByte {
			return node, errorHandler.AddNew(enumErrorCodes.Other, "Not an array at all!")
		}
		if index < 0 {
			return node, errorHandler.AddNew(enumErrorCodes.Other, "Index must be positive")
		}
		if index > len((*parentNode.ArgumentList))-1 {
			return node, errorHandler.AddNew(enumErrorCodes.Other, "NO! Array len")
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
