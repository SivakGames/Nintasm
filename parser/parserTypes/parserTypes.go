package parserTypes

// The possible values for simple operations
type SimpleOperation int

const (
	None SimpleOperation = iota + 1
	Directive
	Instruction
	Label
	Macro
)
