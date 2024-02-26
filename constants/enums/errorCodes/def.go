package enumErrorCodes

type Def int

const (
	None Def = iota + 0
	UnresolvedSymbol
	SyntaxError
)
