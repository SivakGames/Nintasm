package operandFactory

import (
	"fmt"
	"misc/nintasm/tokenizer/tokenizerSpec"
	"strings"
)

type Node = NodeStruct

type nodeType int

const (
	NodeTypeEmpty nodeType = iota + 0
	NodeTypeError
	NodeTypeBinaryExpression
	NodeTypeCallExpression
	NodeTypeLogicalExpression
	NodeTypeMemberExpression
	NodeTypeUnaryExpression

	NodeTypeBacktickStringLiteral
	NodeTypeBooleanLiteral
	NodeTypeIdentifier
	NodeTypeNumericLiteral
	NodeTypeStringLiteral
	NodeTypeSubstitutionID
)

type NodeStruct struct {
	NodeType      nodeType
	Resolved      bool
	NodeTokenType tokenizerSpec.TokenType
	NodeValue     string
	AsBool        bool
	AsNumber      int
	Left          *Node
	Right         *Node
	ArgumentList  *[]Node
	//Consequent    *Node
	//Alternate     *Node
}

func newNode(nodeType tokenizerSpec.TokenType, nodeValue string, nType nodeType) NodeStruct {
	return NodeStruct{
		NodeTokenType: nodeType,
		NodeValue:     nodeValue,
		NodeType:      nType,
	}
}

func EmptyNode() Node {
	node := newNode(tokenizerSpec.None, "", NodeTypeEmpty)
	return node
}

// Used for errors during parsing
func ErrorNode(nodeValue string) Node {
	node := newNode(tokenizerSpec.None, nodeValue, NodeTypeError)
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
func CallExpression(callee string, arguments []Node) Node {
	node := newNode(tokenizerSpec.IDENTIFIER, callee, NodeTypeCallExpression)
	node.ArgumentList = &arguments
	return node
}

// """Label.label expressions"""
func MemberExpression(parent string, key string, computed bool) Node {
	var parentKey string
	if !computed {
		parentKey = fmt.Sprintf("%v.%v", parent, key)
	} else {
		parentKey = fmt.Sprintf("%v[%v]", parent, key)
	}

	node := newNode(tokenizerSpec.IDENTIFIER, parentKey, NodeTypeMemberExpression)
	return node
}

// """Standard binary expression"""
// nodeValue is the operator
func BinaryExpression(nodeType tokenizerSpec.TokenType, nodeValue string, left Node, right Node) Node {
	node := newNode(nodeType, nodeValue, NodeTypeBinaryExpression)
	node.Left = &left
	node.Right = &right
	return node
}

// """Standard unary expression"""
func UnaryExpression(nodeType tokenizerSpec.TokenType, nodeValue string, argument Node) Node {
	node := newNode(nodeType, nodeValue, NodeTypeUnaryExpression)
	node.Right = &argument
	return node
}

//===================================================
//Types of identifiers

func Identifier(nodeType tokenizerSpec.TokenType, nodeValue string) Node {
	node := newNode(nodeType, nodeValue, NodeTypeIdentifier)
	return node
}

// """Substitutions for arguments passed to macros or functions"""
func SubstitutionId(nodeType tokenizerSpec.TokenType, nodeValue string) Node {
	capturedString := nodeValue[1:]
	adjustedString := fmt.Sprintf("\\%v", capturedString)
	node := newNode(nodeType, adjustedString, NodeTypeSubstitutionID)
	return node
}

//===================================================
//Literals

// Numbers
func NumericLiteral(nodeType tokenizerSpec.TokenType, nodeValue string, asNumber int) Node {
	node := newNode(nodeType, nodeValue, NodeTypeNumericLiteral)
	node.AsNumber = asNumber
	node.Resolved = true
	return node
}

func ConvertNodeToNumericLiteral(node *Node) {
	node.NodeType = NodeTypeNumericLiteral
	node.NodeTokenType = tokenizerSpec.NUMBER_decimal
	node.NodeValue = fmt.Sprintf("%d", node.AsNumber)
	node.Resolved = true
	return
}

func ConvertNodeToBooleanLiteral(node *Node) {
	node.NodeType = NodeTypeBooleanLiteral
	node.NodeTokenType = tokenizerSpec.None
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

// Any string in 'single' or "double" quotes
func StringLiteral(nodeType tokenizerSpec.TokenType, nodeValue string) Node {
	capturedString := nodeValue[1 : len(nodeValue)-1]
	node := newNode(nodeType, capturedString, NodeTypeStringLiteral)
	node.Resolved = true
	return node
}

// Any string in `backticks` - These will still need to be resolved via interpreter
func BacktickStringLiteral(nodeType tokenizerSpec.TokenType, nodeValue string) Node {
	capturedString := nodeValue[1 : len(nodeValue)-1]
	capturedString = strings.TrimSpace(capturedString)
	node := newNode(nodeType, capturedString, NodeTypeBacktickStringLiteral)
	return node
}
