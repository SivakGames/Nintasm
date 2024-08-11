package enumSymbolTableTypes

type Def int

const (
	SymbolAsNode Def = iota + 0
	CharMap
	ExprMap
	Function
	Label
	Macro
	Namespace
)
