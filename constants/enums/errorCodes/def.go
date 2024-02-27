package enumErrorCodes

type Def string

const (
	None             Def = "None"
	UnresolvedSymbol     = "UnresolvedSymbol"

	IncludeFileNotExist = "IncludeFileNotExist"
	FailOpenFile        = "FailOpenFile"
	FailScanFile        = "FailScanFile"

	NodeTypeNotNumeric   = "NodeTypeNotNumeric"
	NodeTypeNotString    = "NodeTypeNotString"
	NodeValueNotPositive = "NodeValueNotPositive"
	NodeValueNot8Bit     = "NodeValueNot8Bit"
	NodeValueNotPowerOf2 = "NodeValueNotPowerOf2"

	NodeValueNotGT        = "NodeValueNotGT"
	NodeValueNotLT        = "NodeValueNotLT"
	NodeValueNotGTE       = "NodeValueNotGTE"
	NodeValueNotLTE       = "NodeValueNotLTE"
	NodeValueNotGTEandLTE = "NodeValueNotGTEandLTE"

	InvalidValueAlias = "InvalidValueAlias"
	UnacceptableAlias = "UnacceptableAlias"

	INESValueAlreadySet   = "INESValueAlreadySet"
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
