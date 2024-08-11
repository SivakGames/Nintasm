package interpreter

import (
	"fmt"
	"misc/nintasm/assemble/errorHandler"
	enumErrorCodes "misc/nintasm/constants/enums/errorCodes"
	enumNodeTypes "misc/nintasm/constants/enums/nodeTypes"
	"misc/nintasm/interpreter/environment"
	"misc/nintasm/interpreter/operandFactory"
)

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++

type Node = operandFactory.Node

// ======================================================

func EvaluateNode(node Node) (Node, error) {
	switch node.NodeType {
	case enumNodeTypes.Empty,
		enumNodeTypes.BooleanLiteral,
		enumNodeTypes.NumericLiteral,
		enumNodeTypes.ImplicitUndefined,
		enumNodeTypes.StringLiteral:
		return node, nil

	case enumNodeTypes.Array:
		return processArray(node)

	case enumNodeTypes.Identifier:
		resolvedNode, resolved, err := environment.LookupIdentifierInSymbolAsNodeTable(node.NodeValue)
		if resolved {
			return resolvedNode, err
		}
		return node, err

	case enumNodeTypes.MemberExpression:
		// node[index] style
		if node.ArgumentList != nil && len(*node.ArgumentList) > 0 {
			return processComputedMemberExpression(node)
		}
		// node.key style
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
		return processBacktickLiteral(node)

	case enumNodeTypes.AssignLabelExpression,
		enumNodeTypes.AssignmentExpression:
		return processAssignmentExpression(node)

	case enumNodeTypes.BinaryExpression:
		return processBinaryExpression(node)

	case enumNodeTypes.TernaryExpression:
		return processTernaryExpression(node)

	case enumNodeTypes.UnaryExpression:
		return processUnaryExpression(node)

	case enumNodeTypes.CallExpression:
		if isAssemblerFunction(&node) {
			err := processAssemblerFunction(&node)
			return node, err
		}
		return processCallExpression(node)

	default:
		errorHandler.AddNew(enumErrorCodes.OtherFatal, "???")
		errMsg := fmt.Sprintf("🛑 UNKNOWN NODE BEING INTERPRETED!!! %v", node.NodeValue)
		panic(errMsg)
	}

}
