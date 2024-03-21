package interpreter

import (
	"misc/nintasm/assemble/errorHandler"
	enumErrorCodes "misc/nintasm/constants/enums/errorCodes"
	"misc/nintasm/interpreter/environment/exprmapTable"
	"misc/nintasm/interpreter/operandFactory"
)

func processBacktickLiteral(node Node) (Node, error) {
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
}