package operandFactory

import enumNodeTypes "misc/nintasm/constants/enums/nodeTypes"

//-----------------------------------Primitives

func ValidateNodeIsEmpty(node *Node) bool {
	return node.NodeType == enumNodeTypes.Empty
}

//-----------------------------------Primitives

func ValidateNodeIsBoolean(node *Node) bool {
	return node.NodeType == enumNodeTypes.BooleanLiteral
}
func ValidateNodeIsString(node *Node) bool {
	return node.NodeType == enumNodeTypes.StringLiteral
}

//-----------------------------------Identifiers/IDs

func ValidateNodeIsIdentifier(node *Node) bool {
	return node.NodeType == enumNodeTypes.Identifier
}
func ValidateNodeIsSubstitutionID(node *Node) bool {
	return node.NodeType == enumNodeTypes.SubstitutionID
}

//-----------------------------------1234567890

func ValidateNodeIsNumeric(node *Node) bool {
	return node.NodeType == enumNodeTypes.NumericLiteral
}
func ValidateNumericNodeIsGTValue(node *Node, minValue int) bool {
	return node.AsNumber > minValue
}
func ValidateNumericNodeIsGTEValue(node *Node, minValue int) bool {
	return node.AsNumber >= minValue
}
func ValidateNumericNodeIsLTValue(node *Node, maxValue int) bool {
	return node.AsNumber < maxValue
}
func ValidateNumericNodeIsLTEValue(node *Node, maxValue int) bool {
	return node.AsNumber <= maxValue
}
func ValidateNumericNodeIsGTEandLTEValues(node *Node, minValue int, maxValue int) bool {
	return node.AsNumber >= minValue && node.AsNumber <= maxValue
}
func ValidateNumericNodeIsPositive(node *Node) bool {
	return node.AsNumber >= 0
}
func ValidateNumericNodeIsGTZero(node *Node) bool {
	return node.AsNumber > 0
}
func ValidateNumericNodeIs8BitValue(node *Node) bool {
	return node.AsNumber > -0x00100 && node.AsNumber < 0x00100
}
func ValidateNumericNodeIs16BitValue(node *Node) bool {
	return node.AsNumber > -0x10000 && node.AsNumber < 0x10000
}
