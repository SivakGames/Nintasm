package interpreter

import (
	"fmt"
	"misc/nintasm/interpreter/environment"
	"misc/nintasm/parser/operandFactory"
)

type Node = operandFactory.Node

type EvaluateType int

func InterpretOperands(nodes []Node) []Node {
	for i, n := range nodes {
		newNode := evaluate(n)
		nodes[i] = newNode
	}

	return nodes
}

func evaluate(node Node) Node {
	switch node.NodeType {
	case operandFactory.NodeTypeEmpty,
		operandFactory.NodeTypeStringLiteral,
		operandFactory.NodeTypeNumericLiteral:
		return node

	case operandFactory.NodeTypeIdentifier,
		operandFactory.NodeTypeMemberExpression:
		return environment.LookupInEnvironment(node.NodeValue)

	case operandFactory.NodeTypeBinaryExpression:
		left := evaluate(*node.Left)
		right := evaluate(*node.Right)
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
		}

	case operandFactory.NodeTypeUnaryExpression:
		right := evaluate(*node.Right)
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
		}
	default:
		fmt.Println("UNKNOWN NODE!!!")
	}

	return node
}
