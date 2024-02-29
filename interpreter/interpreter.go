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
	"misc/nintasm/interpreter/operandFactory"
	"strings"
)

type Node = operandFactory.Node

type assemblerFunction struct {
	minArgs          int
	maxArgs          int
	argMustResolveTo []enumNodeTypes.Def
}

var PopParentLabelWhenBlockOpDone bool = false
var parentLabelStack []string

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

	case enumNodeTypes.BacktickStringLiteral:
		_, err := exprmapTable.GetCurrentExprmap()
		if err != nil {
			return node, err
		}

		exprValue, exists := exprmapTable.CheckIfDefinedInExprmap(node.NodeValue)
		if !exists {
			return node, errorHandler.AddNew(enumErrorCodes.ToExprMapUndefExpr, node.NodeValue)
		}

		node.AsNumber = exprValue
		operandFactory.ConvertNodeToNumericLiteral(&node)
		return node, nil

	case enumNodeTypes.Identifier,
		enumNodeTypes.MemberExpression:
		return environment.LookupInEnvironment(node.NodeValue)

	case enumNodeTypes.AssignmentExpression:
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

		_, err = environment.AddToEnvironment(symbolName, right)
		if err != nil {
			return node, err
		}
		return node, nil

	case enumNodeTypes.BinaryExpression:
		var err error
		left, err := EvaluateNode(*node.Left)
		if err != nil {
			return node, err
		}
		right, err := EvaluateNode(*node.Right)
		if err != nil {
			return node, err
		}

		operation := node.NodeValue
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
			panic("SOMETHING IS VERY WRONG WITH BINARY EXPRESSION INTERPRETING")
		}
		switch operation {
		case "+", "-", "*", "/", "%", "|", "&", "^", "<<", ">>":
			operandFactory.ConvertNodeToNumericLiteral(&node)
		case "<", "<=", ">", ">=", "==", "!=", "&&", "||":
			operandFactory.ConvertNodeToBooleanLiteral(&node)
		}

	case enumNodeTypes.UnaryExpression:
		right, err := EvaluateNode(*node.Right)
		if err != nil {
			return node, err
		}
		operation := node.NodeValue
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
			panic("SOMETHING IS VERY WRONG WITH UNARY EXPRESSION INTERPRETING")
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

		// Look up user def functions
		fmt.Println(node)

	default:
		errMsg := fmt.Sprintf("UNKNOWN NODE BEING INTERPRETED!!! %v", node.NodeValue)
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

//++++++++++++++++++++++++++++++++++++++++++++++++++++++++

func GetParentLabel() (string, error) {
	if len(parentLabelStack) == 0 {
		return "", errorHandler.AddNew(enumErrorCodes.InterpreterNoParentLabel)
	}
	parentLabel := parentLabelStack[len(parentLabelStack)-1]
	return parentLabel, nil
}

func AppendParentLabel(newLabel string) {
	parentLabelStack = append(parentLabelStack, newLabel)
	return
}

func PopParentLabel() {
	parentLabelStack = parentLabelStack[:len(parentLabelStack)-1]
	return
}

// Will overwrite at current position or add if none
func OverwriteParentLabel(newLabel string) {
	if len(parentLabelStack) == 0 {
		parentLabelStack = append(parentLabelStack, newLabel)
		return
	}
	parentLabelStack[len(parentLabelStack)-1] = newLabel
	return
}
