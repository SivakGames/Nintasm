package enumNodeTypes

type Def int

const (
	Empty Def = iota + 0
	Error
	AssemblerReservedWord

	AssignmentExpression
	AssignLabelExpression
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

	MacroReplacementString
)
