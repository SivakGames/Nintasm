package enumParserTypes

// The possible values for simple operations
type Def int

const (
	None Def = iota + 1
	CaptureBlock
	Directive
	Instruction
	Label
	Macro
)
