package operandFactory

import (
	"fmt"
	enumNodeTypes "misc/nintasm/constants/enums/nodeTypes"
	enumTokenTypes "misc/nintasm/constants/enums/tokenTypes"
	"strings"
)

type Node = NodeStruct
type nodeEnum = enumNodeTypes.Def
type tokenEnum = enumTokenTypes.Def

type NodeStruct struct {
	NodeType      nodeEnum
	Resolved      bool
	NodeTokenType tokenEnum
	NodeValue     string
	AsBool        bool
	AsNumber      int
	Left          *Node
	Right         *Node
	ArgumentList  *[]Node
	//Consequent    *Node
	//Alternate     *Node
}

func newNode(nodeType tokenEnum, nodeValue string, nType nodeEnum) NodeStruct {
	return NodeStruct{
		NodeTokenType: nodeType,
		NodeValue:     nodeValue,
		NodeType:      nType,
	}
}

func EmptyNode() Node {
	node := newNode(enumTokenTypes.None, "", enumNodeTypes.Empty)
	return node
}

// Used for errors during parsing
func ErrorNode(nodeValue string) Node {
	node := newNode(enumTokenTypes.None, nodeValue, enumNodeTypes.Error)
	return node
}

/*
// """Operand for .func expressions"""
func FunctionOperand(numArgs int, operands []string) (int, []string) {
	return numArgs, operands
}

// """Any instruction that takes an operand"""
func InterpretedInstructionBytes(opcode []string, value []string) []string {
	//return opcode, ...value
	return opcode
} */

// """Calling a function"""
func CreateCallExpressionNode(callee string, arguments []Node) Node {
	node := newNode(enumTokenTypes.IDENTIFIER, callee, enumNodeTypes.CallExpression)
	node.ArgumentList = &arguments
	return node
}

// Called by interpreter at time of creating global env.
func CreateAssemblerReservedWordNode(nodeValue string) Node {
	node := newNode(enumTokenTypes.IDENTIFIER, nodeValue, enumNodeTypes.AssemblerReservedWord)
	node.Resolved = true
	return node
}

// """Label.label expressions"""
func CreateMemberExpressionNode(parent string, key string, computed bool) Node {
	var parentKey string
	if !computed {
		parentKey = fmt.Sprintf("%v.%v", parent, key)
	} else {
		parentKey = fmt.Sprintf("%v[%v]", parent, key)
	}

	node := newNode(enumTokenTypes.IDENTIFIER, parentKey, enumNodeTypes.MemberExpression)
	return node
}

// """Standard binary expression"""
// nodeValue is the operator
func CreateBinaryExpressionNode(nodeType tokenEnum, nodeValue string, left Node, right Node) Node {
	node := newNode(nodeType, nodeValue, enumNodeTypes.BinaryExpression)
	node.Left = &left
	node.Right = &right
	return node
}

// """Standard unary expression"""
func CreateUnaryExpressionNode(nodeType tokenEnum, nodeValue string, right Node) Node {
	node := newNode(nodeType, nodeValue, enumNodeTypes.UnaryExpression)
	node.Right = &right
	return node
}

// """Assign a symbol"""
func CreateAssignmentNode(left Node, right Node) Node {
	node := newNode(enumTokenTypes.ASSIGN_simple, "=", enumNodeTypes.AssignmentExpression)
	node.Left = &left
	node.Right = &right
	return node
}

//===================================================
//Types of identifiers

func CreateIdentifierNode(nodeType tokenEnum, nodeValue string) Node {
	node := newNode(nodeType, nodeValue, enumNodeTypes.Identifier)
	return node
}

// """Substitutions for arguments passed to macros or functions"""
func CreateSubstitutionIdNode(nodeType tokenEnum, nodeValue string) Node {
	capturedString := nodeValue[1:]
	adjustedString := fmt.Sprintf("\\%v", capturedString)
	node := newNode(nodeType, adjustedString, enumNodeTypes.SubstitutionID)
	return node
}

//===================================================
//Literals

func CreateBooleanLiteralNode(nodeType tokenEnum, nodeValue string, asBool bool) Node {
	node := newNode(nodeType, nodeValue, enumNodeTypes.BooleanLiteral)
	node.AsBool = asBool
	node.Resolved = true
	return node
}

// Numbers
func CreateNumericLiteralNode(nodeType tokenEnum, nodeValue string, asNumber int) Node {
	node := newNode(nodeType, nodeValue, enumNodeTypes.NumericLiteral)
	node.AsNumber = asNumber
	node.Resolved = true
	return node
}

// Any string in 'single' or "double" quotes
func CreateStringLiteralNode(nodeType tokenEnum, nodeValue string) Node {
	capturedString := nodeValue[1 : len(nodeValue)-1]
	node := newNode(nodeType, capturedString, enumNodeTypes.StringLiteral)
	node.Resolved = true
	return node
}

// Any string in `backticks` - These will still need to be resolved via interpreter
func CreateBacktickStringLiteralNode(nodeType tokenEnum, nodeValue string) Node {
	capturedString := nodeValue[1 : len(nodeValue)-1]
	capturedString = strings.TrimSpace(capturedString)
	node := newNode(nodeType, capturedString, enumNodeTypes.BacktickStringLiteral)
	return node
}

//->->->->->->->->->->->->->->->->->->->->->->->->->->
// Conversions

func ConvertNodeToNumericLiteral(node *Node) {
	node.NodeType = enumNodeTypes.NumericLiteral
	node.NodeTokenType = enumTokenTypes.NUMBER_decimal
	node.NodeValue = fmt.Sprintf("%d", node.AsNumber)
	node.Resolved = true
	return
}

func ConvertNodeToBooleanLiteral(node *Node) {
	node.NodeType = enumNodeTypes.BooleanLiteral
	node.NodeTokenType = enumTokenTypes.None
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
	node.NodeTokenType = enumTokenTypes.STRING
	node.Resolved = true
	return
}

func ValidateNodeIsBoolean(node *Node) bool {
	return node.NodeType == enumNodeTypes.BooleanLiteral
}
func ValidateNodeIsString(node *Node) bool {
	return node.NodeType == enumNodeTypes.StringLiteral
}
func ValidateNodeIsSubstitutionID(node *Node) bool {
	return node.NodeType == enumNodeTypes.SubstitutionID
}
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
