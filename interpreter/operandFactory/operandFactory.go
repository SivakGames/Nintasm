package operandFactory

import (
	"fmt"
	enumNodeTypes "misc/nintasm/constants/enums/nodeTypes"
	"strings"
)

type Node = NodeStruct
type nodeEnum = enumNodeTypes.Def

type NodeStruct struct {
	NodeType     nodeEnum
	Resolved     bool
	NodeValue    string
	AsBool       bool
	AsNumber     float64
	Left         *Node
	Right        *Node
	ArgumentList *[]Node
}

func newNode(nodeValue string, nType nodeEnum) NodeStruct {
	return NodeStruct{
		NodeValue: nodeValue,
		NodeType:  nType,
	}
}

//----------------------------------------

func EmptyNode() Node {
	node := newNode("", enumNodeTypes.Empty)
	return node
}

// Used for errors during parsing
func ErrorNode(nodeValue string) Node {
	node := newNode(nodeValue, enumNodeTypes.Error)
	return node
}

//----------------------------------------

// """Calling a function"""
func CreateCallExpressionNode(callee string, arguments []Node) Node {
	node := newNode(callee, enumNodeTypes.CallExpression)
	node.ArgumentList = &arguments
	return node
}

// Called by interpreter at time of creating global env.
func CreateAssemblerReservedWordNode(nodeValue string) Node {
	node := newNode(nodeValue, enumNodeTypes.AssemblerReservedWord)
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

	node := newNode(parentKey, enumNodeTypes.MemberExpression)
	return node
}

// """Standard binary expression"""
// nodeValue is the operator
func CreateBinaryExpressionNode(nodeValue string, left Node, right Node) Node {
	node := newNode(nodeValue, enumNodeTypes.BinaryExpression)
	node.Left = &left
	node.Right = &right
	return node
}

// """Standard unary expression"""
func CreateUnaryExpressionNode(nodeValue string, right Node) Node {
	node := newNode(nodeValue, enumNodeTypes.UnaryExpression)
	node.Right = &right
	return node
}

//====================================================

// """Assign a symbol"""
func CreateAssignmentNode(left Node, right Node) Node {
	node := newNode("=", enumNodeTypes.AssignmentExpression)
	node.Left = &left
	node.Right = &right
	return node
}

// """Assign a label symbol"""
func CreateAssignLabelNode(labelName string, org int) Node {
	node := newNode(labelName, enumNodeTypes.AssignLabelExpression)
	left := CreateIdentifierNode(labelName)
	right := CreateNumericLiteralNode(float64(org))
	node.Left = &left
	node.Right = &right
	return node
}

// ++++++++++++++++++++++++++++++++++++++++++++++++++++

// Helper function to take a symbol and number and make an assignment node
func CreateAssignmentNodeForNumber(symbolName string, number float64) Node {
	return CreateAssignmentNode(CreateIdentifierNode(symbolName), CreateNumericLiteralNode(number))
}

//====================================================
//Types of identifiers

func CreateIdentifierNode(nodeValue string) Node {
	node := newNode(nodeValue, enumNodeTypes.Identifier)
	return node
}

// """Substitutions for arguments passed to macros or functions"""
func CreateSubstitutionIdNode(nodeValue string) Node {
	capturedString := nodeValue[1:]
	adjustedString := fmt.Sprintf("\\%v", capturedString)
	node := newNode(adjustedString, enumNodeTypes.SubstitutionID)
	return node
}

// Special node recognized by macros
func CreateMacroReplacementNode(nodeValue string) Node {
	node := newNode(nodeValue, enumNodeTypes.MacroReplacementString)
	return node
}

// ===================================================
// Multi node
func CreateMultiByteNode(nodes []Node) Node {
	node := newNode("multiNode", enumNodeTypes.MultiByte)
	node.ArgumentList = &nodes
	node.Resolved = true
	return node
}

//===================================================
//Literals

func CreateBooleanLiteralNode(asBool bool) Node {
	var val string
	if asBool {
		val = "1"
	} else {
		val = "0"
	}
	node := newNode(val, enumNodeTypes.BooleanLiteral)
	node.AsBool = asBool
	node.Resolved = true
	return node
}

// Numbers
func CreateNumericLiteralNode(asNumber float64) Node {
	node := newNode(fmt.Sprintf("%d", int(asNumber)), enumNodeTypes.NumericLiteral)
	node.AsNumber = asNumber
	node.Resolved = true
	return node
}

// Any string in 'single' or "double" quotes
func CreateStringLiteralNode(nodeValue string) Node {
	capturedString := nodeValue[1 : len(nodeValue)-1]
	node := newNode(capturedString, enumNodeTypes.StringLiteral)
	node.Resolved = true
	return node
}

// Any string in `backticks` - These will still need to be resolved via interpreter
func CreateBacktickStringLiteralNode(nodeValue string) Node {
	capturedString := nodeValue[1 : len(nodeValue)-1]
	capturedString = strings.TrimSpace(capturedString)
	node := newNode(capturedString, enumNodeTypes.BacktickStringLiteral)
	return node
}

// A _ or can also be made when invoking a macro with arguments missing
func CreateUndefinedNode(nodeValue string) Node {
	node := newNode(nodeValue, enumNodeTypes.Undefined)
	return node
}
