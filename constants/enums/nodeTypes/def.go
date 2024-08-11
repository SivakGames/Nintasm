package enumNodeTypes

type Def int

const (
	Empty Def = iota + 0
	Error
	AssemblerReservedWord

	Undefined
	ImplicitUndefined

	AssignmentExpression
	AssignLabelExpression
	BinaryExpression
	CallExpression
	LogicalExpression
	MemberExpression
	TernaryExpression
	UnaryExpression

	Array
	BacktickStringLiteral
	BooleanLiteral
	Identifier
	NumericLiteral
	StringLiteral
	SubstitutionID

	MacroReplacementString
)
