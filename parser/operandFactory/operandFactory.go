package operandFactory

import (
	"fmt"
	"misc/nintasm/tokenizer/tokenizerSpec"
	"strings"
)

type Node struct {
	NodeType     tokenizerSpec.TokenType
	NodeValue    string
	Left         *Node
	Right        *Node
	Argument     *Node
	ArgumentList *[]Node
	Consequent   *Node
	Alternate    *Node
}

func newNode(nodeType tokenizerSpec.TokenType, nodeValue string) Node {
	return Node{
		NodeType:  nodeType,
		NodeValue: nodeValue,
	}
}

func EmptyNode() Node {
	return newNode(tokenizerSpec.None, "")
}

// Used for errors during parsing
func ErrorNode(nodeValue string) Node {
	return newNode(tokenizerSpec.None, nodeValue)
}

// """Operand for .func expressions"""
func FunctionOperand(numArgs int, operands []string) (int, []string) {
	return numArgs, operands
}

// """Any instruction that takes an operand"""
func InterpretedInstructionBytes(opcode []string, value []string) []string {
	//return opcode, ...value
	return opcode
}

// """Calling a function"""
func CallExpression(callee string, arguments []Node) Node {
	node := newNode(tokenizerSpec.IDENTIFIER, callee)
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

	node := newNode(tokenizerSpec.IDENTIFIER, parentKey)
	return node
}

// """Standard binary expression"""
// nodeValue is the operator
func BinaryExpression(nodeType tokenizerSpec.TokenType, nodeValue string, left Node, right Node) Node {
	node := newNode(nodeType, nodeValue)
	node.Left = &left
	node.Right = &right
	return node
}

// """Basically same as binary expression"""
// nodeValue is the operator
func LogicalExpression(nodeType tokenizerSpec.TokenType, nodeValue string, left Node, right Node) Node {
	node := newNode(nodeType, nodeValue)
	node.Left = &left
	node.Right = &right
	return node
}

// """Standard unary expression"""
func UnaryExpression(nodeType tokenizerSpec.TokenType, nodeValue string, argument Node) Node {
	node := newNode(nodeType, nodeValue)
	node.Argument = &argument
	return node
}

//===================================================
//Literals

// """Any string in single or double quotes"""
func StringLiteral(nodeType tokenizerSpec.TokenType, nodeValue string) Node {
	capturedString := nodeValue[1 : len(nodeValue)-1]
	adjustedString := fmt.Sprintf("\"%v\"", capturedString)
	node := newNode(nodeType, adjustedString)
	return node
}

// """Any string in `backticks`"""
func BacktickStringLiteral(nodeType tokenizerSpec.TokenType, nodeValue string) Node {
	capturedString := nodeValue[1 : len(nodeValue)-1]
	capturedString = strings.TrimSpace(capturedString)
	adjustedString := fmt.Sprintf("`%v`", capturedString)
	node := newNode(nodeType, adjustedString)
	return node
}

func NumericLiteral(nodeType tokenizerSpec.TokenType, nodeValue string) Node {
	node := newNode(nodeType, nodeValue)
	return node
}

func Identifier(nodeType tokenizerSpec.TokenType, nodeValue string) Node {
	node := newNode(nodeType, nodeValue)
	return node
}

// """Substitutions for arguments passed to macros or functions"""
func SubstitutionId(nodeType tokenizerSpec.TokenType, nodeValue string) Node {
	capturedString := nodeValue[1:]
	adjustedString := fmt.Sprintf("\\%v", capturedString)
	node := newNode(nodeType, adjustedString)
	return node
}
