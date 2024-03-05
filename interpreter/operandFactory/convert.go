package operandFactory

import (
	"fmt"
	enumNodeTypes "misc/nintasm/constants/enums/nodeTypes"
	enumTokenTypes "misc/nintasm/constants/enums/tokenTypes"
)

//->->->->->->->->->->->->->->->->->->->->->->->->->->
// Conversions

func ConvertNodeToNumericLiteral(node *Node) {
	node.NodeType = enumNodeTypes.NumericLiteral
	node.NodeTokenEnum = enumTokenTypes.NUMBER_decimal
	node.NodeValue = fmt.Sprintf("%d", node.AsNumber)
	node.Resolved = true
	return
}

func ConvertNodeToBooleanLiteral(node *Node) {
	node.NodeType = enumNodeTypes.BooleanLiteral
	node.NodeTokenEnum = enumTokenTypes.None
	if node.AsBool {
		node.NodeValue = "1"
		node.AsNumber = 1
	} else {
		node.NodeValue = "0"
		node.AsNumber = 0
	}
	node.Resolved = true
	return
}

func ConvertNodeToStringLiteral(node *Node) {
	node.NodeType = enumNodeTypes.StringLiteral
	node.NodeTokenEnum = enumTokenTypes.STRING
	node.Resolved = true
	return
}

//-----------------------------------
//Special node for branch instructions

func ConvertToBranchBinaryExpressionNode(originalNode Node, orgToSubtract int) Node {
	orgToSubtractNode := CreateNumericLiteralNode(orgToSubtract)
	branchNode := CreateBinaryExpressionNode(enumTokenTypes.OPERATOR_additive, "-", originalNode, orgToSubtractNode)
	return branchNode
}
