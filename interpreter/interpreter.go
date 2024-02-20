package interpreter

import (
	"errors"
	"fmt"
	"log"
	enumNodeTypes "misc/nintasm/constants/enums/nodeTypes"
	"misc/nintasm/interpreter/environment"
	"misc/nintasm/interpreter/environment/charmapTable"
	"misc/nintasm/interpreter/environment/exprmapTable"
	"misc/nintasm/interpreter/operandFactory"
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

	case enumNodeTypes.BacktickStringLiteral:
		_, err := exprmapTable.GetCurrentExprmap()
		if err != nil {
			return node, err
		}

		exprValue, exists := exprmapTable.CheckIfDefinedInExprmap(node.NodeValue)
		if !exists {
			return node, errors.New("Bad expr char")
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

		fmt.Println(node)
		return node, errors.New("UNKNOWN NODE BEING INTERPRETED!!!")
	}

	return node, nil
}

func ProcessAssemblerFunction(node *Node) (bool, error) {
	funcName := node.NodeValue
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
				return isAsmFunc, errors.New("Argument for node is wrong type...")
			}
		}

		switch funcName {
		case "high":
			node.AsNumber = ((*node.ArgumentList)[0].AsNumber & 0x0ff00) >> 8
		case "low":
			node.AsNumber = ((*node.ArgumentList)[0].AsNumber & 0x000ff)
		case "toCharmap":
			nodeString := ((*node.ArgumentList)[0].NodeValue)
			runeArray := []rune(nodeString)
			currCharmap, err := charmapTable.GetCurrentCharmap()
			if err != nil {
				return isAsmFunc, err
			}

			replacedString := ""
			for _, r := range runeArray {
				d, exists := currCharmap[r]
				if !exists {
					return isAsmFunc, errors.New("Char doesn't exit")
				}
				for _, v := range d {
					replacedString += string(rune(v.AsNumber))
				}
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
