package enumSymbolTableTypes

type Def int

const (
	SymbolAsNode Def = iota + 0
	CharMap
	ExprMap
	Function
	KVMacro
	Label
	Macro
	Namespace
)
