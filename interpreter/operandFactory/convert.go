package operandFactory

import (
	"fmt"
	enumNodeTypes "misc/nintasm/constants/enums/nodeTypes"
)

//->->->->->->->->->->->->->->->->->->->->->->->->->->
// Conversions

func ConvertNodeToNumericLiteral(node *Node) {
	node.NodeType = enumNodeTypes.NumericLiteral
	node.NodeValue = fmt.Sprintf("%d", node.AsNumber)
	node.Resolved = true
	return
}

func ConvertNodeToBooleanLiteral(node *Node) {
	node.NodeType = enumNodeTypes.BooleanLiteral
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
	node.Resolved = true
	return
}

//-----------------------------------
//Special node for branch instructions

func ConvertToBranchBinaryExpressionNode(originalNode Node, orgToSubtract int) Node {
	orgToSubtractNode := CreateNumericLiteralNode(orgToSubtract)
	branchNode := CreateBinaryExpressionNode("-", originalNode, orgToSubtractNode)
	return branchNode
}
