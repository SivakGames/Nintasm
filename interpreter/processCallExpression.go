package interpreter

import (
	"fmt"
	"misc/nintasm/interpreter/environment/funcTable"
	"misc/nintasm/interpreter/environment/symbolAsNodeTable"
)

func processCallExpression(node Node) (Node, error) {
	functionPtr, err := funcTable.LookupAndGetFunctionInEnvironment(node.NodeValue)
	if err != nil {
		return node, err
	}

	functionNode := *functionPtr
	evaluatedArgList := make([]Node, len(*node.ArgumentList))

	for i, n := range *node.ArgumentList {
		evalN, err := EvaluateNode(n)
		if err != nil {
			return node, err
		}
		evaluatedArgList[i] = evalN
	}

	//Add newly evaluated arguments to stack
	symbolAsNodeTable.PushToSymbolTableStack()
	defer symbolAsNodeTable.PopFromSymbolTableStack()
	for i, n := range evaluatedArgList {
		symbolAsNodeTable.AddSymbolToTopTableStack(fmt.Sprintf("\\%d", i+1), n)
	}

	evaluatedFuncNode, err := EvaluateNode(functionNode)
	if err != nil {
		return node, err
	}
	return evaluatedFuncNode, err
}
