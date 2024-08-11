package operandFactory

import (
	"math"
	enumNodeTypes "misc/nintasm/constants/enums/nodeTypes"
)

//-----------------------------------Primitives

func ValidateNodeIsEmpty(node *Node) bool {
	return node.NodeType == enumNodeTypes.Empty
}
func ValidateNodeIsError(node *Node) bool {
	return node.NodeType == enumNodeTypes.Error
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

// -----------------------------------[1,2,3,4]
func ValidateNodeIsArray(node *Node) bool {
	return node.NodeType == enumNodeTypes.Array
}

//-----------------------------------1234567890

func ValidateNodeIsNumeric(node *Node) bool {
	return node.NodeType == enumNodeTypes.NumericLiteral
}
func ValidateNodeIsInt(node *Node) bool {
	_, d := math.Modf(node.AsNumber)
	return d == 0
}
func ValidateNumericNodeIsGTValue(node *Node, minValue int) bool {
	return node.AsNumber > float64(minValue)
}
func ValidateNumericNodeIsGTEValue(node *Node, minValue int) bool {
	return node.AsNumber >= float64(minValue)
}
func ValidateNumericNodeIsLTValue(node *Node, maxValue int) bool {
	return node.AsNumber < float64(maxValue)
}
func ValidateNumericNodeIsLTEValue(node *Node, maxValue int) bool {
	return node.AsNumber <= float64(maxValue)
}
func ValidateNumericNodeIsGTEandLTEValues(node *Node, minValue int, maxValue int) bool {
	return node.AsNumber >= float64(minValue) && node.AsNumber <= float64(maxValue)
}
func ValidateNumericNodeIsPositive(node *Node) bool {
	return node.AsNumber >= 0
}
func ValidateNumericNodeIsGTZero(node *Node) bool {
	return node.AsNumber > 0
}
func ValidateNumericNodeIsSigned8BitValue(node *Node) bool {
	return node.AsNumber > -0x00081 && node.AsNumber < 0x00080
}
func ValidateNumericNodeIs8BitValue(node *Node) bool {
	return node.AsNumber > -0x00100 && node.AsNumber < 0x00100
}
func ValidateNumericNodeIs16BitValue(node *Node) bool {
	return node.AsNumber > -0x10000 && node.AsNumber < 0x10000
}
