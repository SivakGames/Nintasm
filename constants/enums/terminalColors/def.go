package enumTerminalColors

// The possible values for simple operations
type Def int

const (
	None Def = iota + 0
	Red
	Green
	Yellow
	Blue
	Magenta
	Cyan
	LightRed
	LightGreen
	LightYellow
	LightBlue
	LightMagenta
	LightCyan
	ANSI_START
	AnsiPurple
	AnsiTeal
	AnsiGreen
	AnsiDarkSeaGreen
	AnsiSeaGreen
	AnsiOrange
	AnsiRed
	AnsiGray1
	AnsiGray2
	AnsiGray3
	AnsiGray4
	AnsiGray5
)
