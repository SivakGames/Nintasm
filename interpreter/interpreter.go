package interpreter

import (
	"errors"
	"log"
	enumNodeTypes "misc/nintasm/enums/nodeTypes"
	"misc/nintasm/interpreter/environment"
	"misc/nintasm/parser/operandFactory"
)

type Node = operandFactory.Node

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
		log.Println("Cally", *node.ArgumentList)

	default:
		return node, errors.New("UNKNOWN NODE!!!")
	}

	return node, nil
}
