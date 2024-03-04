package interpreter

import (
	"fmt"
	"log"
	"misc/nintasm/assemble/errorHandler"
	enumErrorCodes "misc/nintasm/constants/enums/errorCodes"
	enumNodeTypes "misc/nintasm/constants/enums/nodeTypes"
	"misc/nintasm/interpreter/environment"
	"misc/nintasm/interpreter/environment/charmapTable"
	"misc/nintasm/interpreter/environment/exprmapTable"
	"misc/nintasm/interpreter/environment/funcTable"
	"misc/nintasm/interpreter/environment/symbolAsNodeTable"
	"misc/nintasm/interpreter/operandFactory"
	"strings"
)

type Node = operandFactory.Node

type assemblerFunction struct {
	minArgs          int
	maxArgs          int
	argMustResolveTo []enumNodeTypes.Def
}

var assemblerBuiltInFunctions = map[string]assemblerFunction{
	"high":      {1, 1, []enumNodeTypes.Def{enumNodeTypes.NumericLiteral}},
	"low":       {1, 1, []enumNodeTypes.Def{enumNodeTypes.NumericLiteral}},
	"bank":      {1, 1, []enumNodeTypes.Def{enumNodeTypes.Identifier}},
	"toCharmap": {1, 1, []enumNodeTypes.Def{enumNodeTypes.StringLiteral}},
}

func EvaluateNode(node Node) (Node, error) {
	switch node.NodeType {
	case enumNodeTypes.Empty,
		enumNodeTypes.BooleanLiteral,
		enumNodeTypes.NumericLiteral,
		enumNodeTypes.StringLiteral:
		return node, nil

	case enumNodeTypes.Error:
		return node, nil

	case enumNodeTypes.Identifier,
		enumNodeTypes.MemberExpression:
		resolvedNode, resolved, err := environment.LookupIdentifierInSymbolAsNodeTable(node.NodeValue)
		if resolved {
			return resolvedNode, err
		}
		return node, err

	case enumNodeTypes.BacktickStringLiteral:
		_, err := exprmapTable.GetCurrentExprmap()
		if err != nil {
			return node, err
		}

	case enumNodeTypes.SubstitutionID:
		substitutionNode, err := environment.LookupSubstitutionIDfromStack(node.NodeValue)
		if err != nil {
			return substitutionNode, err
		}
		return EvaluateNode(substitutionNode)

	case enumNodeTypes.AssignLabelExpression,
		enumNodeTypes.AssignmentExpression:
		var err error
		left := *node.Left
		right, err := EvaluateNode(*node.Right)
		if err != nil {
			return node, err
		}
		symbolName := left.NodeValue
		node.Left = nil
		node.Right = nil
		isLocal := strings.HasPrefix(symbolName, ".")
		if isLocal {
			parentLabel, err := GetParentLabel()
			if err != nil {
				return node, err
			}
			symbolName = parentLabel + symbolName
		}
		node.Resolved = true

		isLabel := node.NodeType == enumNodeTypes.AssignLabelExpression
		err = environment.AddIdentifierToSymbolAsNodeTable(symbolName, right)
		if err != nil {
			return node, err
		}
		if isLabel {
			environment.AddToLabelAsBankTable(symbolName)
		}

		return node, nil

	case enumNodeTypes.BinaryExpression:
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
			return node, errorHandler.AddNew(enumErrorCodes.InterpreterBinaryMismatchedTypes, badLeftValue, operation, badRightValue)
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
			node.AsNumber = left.AsNumber % right.AsNumber
		case "|":
			node.AsNumber = left.AsNumber | right.AsNumber
		case "&":
			node.AsNumber = left.AsNumber & right.AsNumber
		case "^":
			node.AsNumber = left.AsNumber ^ right.AsNumber
		case "<<":
			node.AsNumber = left.AsNumber << right.AsNumber
		case ">>":
			node.AsNumber = left.AsNumber >> right.AsNumber
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

	case enumNodeTypes.UnaryExpression:
		operation := node.NodeValue
		right, err := EvaluateNode(*node.Right)
		if err != nil {
			return node, err
		}

		if right.NodeType != enumNodeTypes.NumericLiteral {
			return node, errorHandler.AddNew(enumErrorCodes.InterpreterUnaryNotNumeric, operation, right.NodeValue)
		}

		node.Right = nil
		switch operation {
		case "+":
			node.AsNumber = right.AsNumber
		case "-":
			node.AsNumber = -right.AsNumber
		case "~":
			node.AsNumber = ^right.AsNumber
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

	case enumNodeTypes.CallExpression:
		wasAsmFunc, err := ProcessAssemblerFunction(&node)
		if err != nil {
			return node, err
		}
		if wasAsmFunc {
			return node, nil
		}

		// Not ASM function, look for user-def function
		functionNode, err := funcTable.LookupAndGetFunctionInEnvironment(node.NodeValue)
		if err != nil {
			return node, err
		}
		symbolAsNodeTable.PushToSymbolTableStack()
		for i, n := range *node.ArgumentList {
			symbolAsNodeTable.AddSymbolToTopTableStack(fmt.Sprintf("\\%d", i+1), n)
		}
		evaluatedFuncNode, err := EvaluateNode(*functionNode)
		symbolAsNodeTable.PopFromSymbolTableStack()
		if err != nil {
			return node, err
		}
		return evaluatedFuncNode, err

	default:
		errorHandler.AddNew(enumErrorCodes.Other, "???")
		errMsg := fmt.Sprintf("🛑 UNKNOWN NODE BEING INTERPRETED!!! %v", node.NodeValue)
		panic(errMsg)
	}

	return node, nil
}

func ProcessAssemblerFunction(node *Node) (bool, error) {
	funcName := node.NodeValue
	functionData, isAsmFunc := assemblerBuiltInFunctions[funcName]
	if isAsmFunc {
		numArgs := len(*node.ArgumentList)
		if numArgs < functionData.minArgs {
			return isAsmFunc, errorHandler.AddNew(enumErrorCodes.InterpreterFuncTooFewArgs)
		}
		if numArgs > functionData.maxArgs {
			return isAsmFunc, errorHandler.AddNew(enumErrorCodes.InterpreterFuncTooManyArgs)
		}

		for i, a := range *node.ArgumentList {
			if a.NodeType != functionData.argMustResolveTo[i] {
				return isAsmFunc, errorHandler.AddNew(enumErrorCodes.InterpreterFuncArgWrongType)
			}
		}

		switch funcName {
		case "high":
			node.AsNumber = ((*node.ArgumentList)[0].AsNumber & 0x0ff00) >> 8
		case "low":
			node.AsNumber = ((*node.ArgumentList)[0].AsNumber & 0x000ff)
		case "toCharmap":
			nodeString := ((*node.ArgumentList)[0].NodeValue)
			replacedString, err := charmapTable.MapStringToCharmap(nodeString)
			if err != nil {
				return isAsmFunc, err
			}
			node.NodeValue = replacedString
		case "bank":
			log.Println((*node.ArgumentList)[0])
			//node.AsNumber = ((*node.ArgumentList)[0].AsNumber & 0x000ff)
		}

		switch funcName {
		case "high", "low":
			operandFactory.ConvertNodeToNumericLiteral(node)
		case "toCharmap":
			operandFactory.ConvertNodeToStringLiteral(node)
		}
	}
	return isAsmFunc, nil
}
