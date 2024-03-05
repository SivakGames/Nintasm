package operandFactory

import (
	"fmt"
	enumNodeTypes "misc/nintasm/constants/enums/nodeTypes"
	enumTokenTypes "misc/nintasm/constants/enums/tokenTypes"
	"strconv"
	"strings"
)

type Node = NodeStruct
type nodeEnum = enumNodeTypes.Def
type tokenEnum = enumTokenTypes.Def

type NodeStruct struct {
	NodeType      nodeEnum
	NodeTokenEnum tokenEnum
	Resolved      bool
	NodeValue     string
	AsBool        bool
	AsNumber      int
	Left          *Node
	Right         *Node
	ArgumentList  *[]Node
}

func newNode(nodeTokenEnum tokenEnum, nodeValue string, nType nodeEnum) NodeStruct {
	return NodeStruct{
		NodeTokenEnum: nodeTokenEnum,
		NodeValue:     nodeValue,
		NodeType:      nType,
	}
}

//----------------------------------------

func EmptyNode() Node {
	node := newNode(enumTokenTypes.None, "", enumNodeTypes.Empty)
	return node
}

// Used for errors during parsing
func ErrorNode(nodeValue string) Node {
	node := newNode(enumTokenTypes.None, nodeValue, enumNodeTypes.Error)
	return node
}

//----------------------------------------

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

//====================================================

// """Assign a symbol"""
func CreateAssignmentNode(left Node, right Node) Node {
	node := newNode(enumTokenTypes.ASSIGN_simple, "=", enumNodeTypes.AssignmentExpression)
	node.Left = &left
	node.Right = &right
	return node
}

// """Assign a label symbol"""
func CreateAssignLabelNode(labelName string, org int) Node {
	node := newNode(enumTokenTypes.ASSIGN_simple, labelName, enumNodeTypes.AssignLabelExpression)
	left := CreateIdentifierNode(labelName)
	right := CreateNumericLiteralNode(org)
	node.Left = &left
	node.Right = &right
	return node
}

//===================================================
//Types of identifiers

func CreateIdentifierNode(nodeValue string) Node {
	node := newNode(enumTokenTypes.IDENTIFIER, nodeValue, enumNodeTypes.Identifier)
	return node
}

// """Substitutions for arguments passed to macros or functions"""
func CreateSubstitutionIdNode(tokenEnumType tokenEnum, nodeValue string) Node {
	capturedString := nodeValue[1:]
	adjustedString := fmt.Sprintf("\\%v", capturedString)
	node := newNode(tokenEnumType, adjustedString, enumNodeTypes.SubstitutionID)
	return node
}

// Special node recognized by macros
func CreateMacroReplacementNode(nodeValue string) Node {
	node := newNode(enumTokenTypes.None, nodeValue, enumNodeTypes.MacroReplacementString)
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
	node := newNode(enumTokenTypes.NUMBER_decimal, val, enumNodeTypes.BooleanLiteral)
	node.AsBool = asBool
	node.Resolved = true
	return node
}

// Numbers
func CreateNumericLiteralNode(asNumber int) Node {
	node := newNode(enumTokenTypes.NUMBER_decimal, strconv.Itoa(asNumber), enumNodeTypes.NumericLiteral)
	node.AsNumber = asNumber
	node.Resolved = true
	return node
}

// Any string in 'single' or "double" quotes
func CreateStringLiteralNode(nodeValue string) Node {
	capturedString := nodeValue[1 : len(nodeValue)-1]
	node := newNode(enumTokenTypes.STRING, capturedString, enumNodeTypes.StringLiteral)
	node.Resolved = true
	return node
}

// Any string in `backticks` - These will still need to be resolved via interpreter
func CreateBacktickStringLiteralNode(nodeValue string) Node {
	capturedString := nodeValue[1 : len(nodeValue)-1]
	capturedString = strings.TrimSpace(capturedString)
	node := newNode(enumTokenTypes.BACKTICK_STRING, capturedString, enumNodeTypes.BacktickStringLiteral)
	return node
}

// A _ or can also be made when invoking a macro with arguments missing
func CreateUndefinedNode(nodeValue string) Node {
	node := newNode(enumTokenTypes.None, nodeValue, enumNodeTypes.Undefined)
	return node
}
