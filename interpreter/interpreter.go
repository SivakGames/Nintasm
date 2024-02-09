package interpreter

import (
	"errors"
	enumNodeTypes "misc/nintasm/enums/nodeTypes"
	"misc/nintasm/interpreter/environment"
	"misc/nintasm/parser/operandFactory"
	"strings"
)

type Node = operandFactory.Node

type assemblerFunction struct {
	minArgs          int
	maxArgs          int
	argMustResolveTo []enumNodeTypes.Def
}

var assemblerBuiltInFunctions = map[string]assemblerFunction{
	"high": {1, 1, []enumNodeTypes.Def{enumNodeTypes.NumericLiteral}},
}

func EvaluateNode(node Node) (Node, error) {
	switch node.NodeType {
	case enumNodeTypes.Empty,
		enumNodeTypes.StringLiteral,
		enumNodeTypes.NumericLiteral:
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
			return node, errors.New("SOMETHING IS VERY WRONG")

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
			return node, errors.New("SOMETHING IS VERY WRONG")
		}

		switch operation {
		case "+", "-", "~":
			operandFactory.ConvertNodeToNumericLiteral(&node)
		case "!":
			operandFactory.ConvertNodeToBooleanLiteral(&node)
		}

	case enumNodeTypes.CallExpression:
		potentialAssemblerFunctionName := strings.ToLower(node.NodeValue)
		wasAsmFunc, err := ProcessAssemblerFunction(&node, potentialAssemblerFunctionName)
		if err != nil {
			return node, err
		}
		if wasAsmFunc {
			switch potentialAssemblerFunctionName {
			case "high":
				node.AsNumber = (node.AsNumber & 0x0ff00) >> 8
			case "low":
				node.AsNumber = (node.AsNumber & 0x000ff)
			}

			switch potentialAssemblerFunctionName {
			case "high", "low":
				operandFactory.ConvertNodeToNumericLiteral(&node)
			}

		}
		// Look up user def functions

	default:
		return node, errors.New("UNKNOWN NODE!!!")
	}

	return node, nil
}

func ProcessAssemblerFunction(node *Node, funcName string) (bool, error) {
	functionData, isAsmFunc := assemblerBuiltInFunctions[funcName]
	if isAsmFunc {
		numArgs := len(*node.ArgumentList)
		if numArgs < functionData.minArgs {
			return isAsmFunc, errors.New("Too few arguments for function!")
		}
		if numArgs > functionData.maxArgs {
			return isAsmFunc, errors.New("Too many arguments for function!")
		}
		for i, a := range *node.ArgumentList {
			if a.NodeType != functionData.argMustResolveTo[i] {
				return isAsmFunc, errors.New("Argument node is wrong type...")
			}
		}
	}
	return isAsmFunc, nil
}
