package enumNodeTypes

type Def int

const (
	Empty Def = iota + 0
	Error
	AssemblerReservedWord
	Undefined

	AssignmentExpression
	AssignLabelExpression
	BinaryExpression
	CallExpression
	LogicalExpression
	MemberExpression
	TernaryExpression
	UnaryExpression

	BacktickStringLiteral
	BooleanLiteral
	Identifier
	NumericLiteral
	StringLiteral
	SubstitutionID

	MultiByte
	MacroReplacementString
)
