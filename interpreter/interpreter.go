package interpreter

import (
	"fmt"
	"misc/nintasm/assemble/errorHandler"
	enumErrorCodes "misc/nintasm/constants/enums/errorCodes"
	enumNodeTypes "misc/nintasm/constants/enums/nodeTypes"
	"misc/nintasm/interpreter/environment"
	"misc/nintasm/interpreter/environment/exprmapTable"
	"misc/nintasm/interpreter/environment/funcTable"
	"misc/nintasm/interpreter/environment/namespaceTable"
	"misc/nintasm/interpreter/environment/symbolAsNodeTable"
	"misc/nintasm/interpreter/operandFactory"
	"strings"
)

type Node = operandFactory.Node

func EvaluateNode(node Node) (Node, error) {
	switch node.NodeType {
	case enumNodeTypes.Empty,
		enumNodeTypes.BooleanLiteral,
		enumNodeTypes.NumericLiteral,
		enumNodeTypes.MultiByte,
		enumNodeTypes.StringLiteral:
		return node, nil

	case enumNodeTypes.Identifier,
		enumNodeTypes.MemberExpression:
		resolvedNode, resolved, err := environment.LookupIdentifierInSymbolAsNodeTable(node.NodeValue)
		if resolved {
			return resolvedNode, err
		}
		return node, err

	case enumNodeTypes.SubstitutionID:
		substitutionNode, err := environment.LookupSubstitutionID(node.NodeValue)
		if err != nil {
			return node, err
		}
		return substitutionNode, nil

	case enumNodeTypes.BacktickStringLiteral:
		_, err := exprmapTable.GetCurrentExprmap()
		if err != nil {
			return node, err
		}
		exprAsNum, exprExists := exprmapTable.CheckIfDefinedInExprmap(node.NodeValue)
		if !exprExists {
			return node, errorHandler.AddNew(enumErrorCodes.ExprMapUndefExpr, node.NodeValue)
		}
		node.AsNumber = exprAsNum
		operandFactory.ConvertNodeToNumericLiteral(&node)
		return node, nil

	case enumNodeTypes.AssignLabelExpression,
		enumNodeTypes.AssignmentExpression:

		assignmentTypeIsLabel := node.NodeType == enumNodeTypes.AssignLabelExpression
		nodeHasResolved := false

		//Left node is the label itself
		symbolName := (*node.Left).NodeValue
		originalSymbolName := symbolName
		isLocal := strings.HasPrefix(symbolName, ".")
		if isLocal {
			parentLabel, err := GetParentLabel()
			if err != nil {
				return node, err
			}
			symbolName = parentLabel + symbolName
			if namespaceTable.IsDefiningNamespace {
				defer func() {
					namespaceTable.AddKeyToCurrentNamespace(parentLabel, originalSymbolName, nodeHasResolved)
				}()
			}
		}

		// Right node is the expression to set the label to
		evaluatedLabelNode, err := EvaluateNode(*node.Right)
		if err != nil {
			return node, err
		}

		err = environment.AddIdentifierToSymbolAsNodeTable(symbolName, evaluatedLabelNode)
		if err != nil {
			return node, err
		}
		if assignmentTypeIsLabel {
			environment.AddToLabelAsBankTable(symbolName)
		}

		nodeHasResolved = true
		node.Left = nil
		node.Right = nil
		return node, nil

	case enumNodeTypes.BinaryExpression:
		var err error
		node, err = processBinaryExpression(node)
		if err != nil {
			return node, err
		}

	case enumNodeTypes.UnaryExpression:
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
		wasAsmFunc, err := processAssemblerFunction(&node)
		if err != nil {
			return node, err
		}
		if wasAsmFunc {
			return node, nil
		}

		// ---------------------------------------------
		// Not ASM function, look for user-def function
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

	default:
		errorHandler.AddNew(enumErrorCodes.Other, "???")
		errMsg := fmt.Sprintf("🛑 UNKNOWN NODE BEING INTERPRETED!!! %v", node.NodeValue)
		panic(errMsg)
	}

	return node, nil
}
