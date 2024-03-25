package operandFactory

import (
	"fmt"
	enumNodeTypes "misc/nintasm/constants/enums/nodeTypes"
)

//->->->->->->->->->->->->->->->->->->->->->->->->->->
// Conversions

func ConvertNodeToNumericLiteral(node *Node) {
	node.NodeType = enumNodeTypes.NumericLiteral
	node.NodeValue = fmt.Sprintf("%d", int(node.AsNumber))
	node.Resolved = true
	node.ArgumentList = nil
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

func ConvertNodeToMultiBytes(node *Node, mutliBytes []Node) {
	node.NodeType = enumNodeTypes.MultiByte
	node.Resolved = true
	node.ArgumentList = &mutliBytes
	return
}

//-----------------------------------
//Special node for branch instructions

func ConvertToBranchBinaryExpressionNode(originalNode Node, orgToSubtract int) Node {
	orgToSubtractNode := CreateNumericLiteralNode(float64(orgToSubtract))
	innerNode := CreateBinaryExpressionNode("-", originalNode, orgToSubtractNode)
	branchNode := CreateBinaryExpressionNode("-", innerNode, CreateNumericLiteralNode(2))
	return branchNode
}
