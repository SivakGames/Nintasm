package enumParserTypes

// The possible values for simple operations
type Def int

const (
	None Def = iota + 1
	Directive
	Instruction
	Label
	Macro
)
