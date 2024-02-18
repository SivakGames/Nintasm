package enumInstructionModes

type Def int

const (
	None Def = iota + 0
	A
	ABS
	ABS_X
	ABS_Y
	IMM
	IMPL
	IND
	IND_X
	IND_Y
	REL
	ZP
	ZP_X
	ZP_Y
)
