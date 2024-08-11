package funcTable

import (
	"misc/nintasm/assemble/errorHandler"
	enumErrorCodes "misc/nintasm/constants/enums/errorCodes"
	"misc/nintasm/interpreter/operandFactory"
)

type Node = operandFactory.Node

var functionTable = map[string]*Node{}

func AddIdentifierKeyToFunctionTable(functionName string) {
	functionTable[functionName] = nil
}

func AddNodeToFunction(functionName string, node *Node) {
	functionTable[functionName] = node
}

func lookupFunctionInEnvironment(functionName string) (*Node, bool) {
	functionData, ok := functionTable[functionName]
	return functionData, ok
}

func LookupAndGetFunctionInEnvironment(functionName string) (*Node, error) {
	functionData, ok := lookupFunctionInEnvironment(functionName)
	if ok {
		return functionData, nil
	} else {
		return nil, errorHandler.AddNew(enumErrorCodes.InterpreterFuncUndefined, functionName)
	}
}
