package interpreter

import "misc/nintasm/interpreter/operandFactory"

func processMultiByte(node Node) (Node, error) {
	evaluatedArgList := make([]Node, len(*node.ArgumentList))
	for i, n := range *node.ArgumentList {
		evalN, err := EvaluateNode(n)
		if err != nil {
			return node, err
		}
		evaluatedArgList[i] = evalN
	}
	evaluatedNode := operandFactory.CreateMultiByteNode(evaluatedArgList)

	return evaluatedNode, nil
}
