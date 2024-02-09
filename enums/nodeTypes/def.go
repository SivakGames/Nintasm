package enumNodeTypes

type Def int

const (
	Empty Def = iota + 0
	Error
	AssemblerBuiltInFunction

	AssignmentExpression
	BinaryExpression
	CallExpression
	LogicalExpression
	MemberExpression
	UnaryExpression

	BacktickStringLiteral
	BooleanLiteral
	Identifier
	NumericLiteral
	StringLiteral
	SubstitutionID
)