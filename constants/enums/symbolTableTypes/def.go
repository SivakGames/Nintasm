package enumSymbolTableTypes

type Def int

const (
	SymbolAsNode Def = iota + 0
	Label
	Macro
	KVMacro
	CharMap
	ExprMap
)
