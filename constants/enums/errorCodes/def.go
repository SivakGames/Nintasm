package enumErrorCodes

type Def string

const (
	None             Def = "None"
	UnresolvedSymbol     = "UnresolvedSymbol"

	IncludeFileNotExist = "IncludeFileNotExist"
	FailOpenFile        = "FailOpenFile"
	FailScanFile        = "FailScanFile"

	INESPRGSet            = "INESPRGSet"
	INESPRGBadValue       = "INESPRGBadValue"
	INESPRGUnacceptable   = "INESPRGUnacceptable"
	INESCHRSet            = "INESCHRSet"
	INESCHRBadValue       = "INESCHRBadValue"
	INESCHRUnacceptable   = "INESCHRUnacceptable"
	INESMapperSet         = "INESMapperSet"
	INESMapperBadValue    = "INESMapperBadValue"
	INESMirroringSet      = "INESMirroringSet"
	INESMirroringBadValue = "INESMirroringBadValue"
)

type Severity int

const (
	// Will just be a warning - stops nothing
	Warning Severity = iota
	// Stop operand processing
	Error
	// Stop building entirely
	Fatal
)
