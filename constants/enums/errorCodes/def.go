package enumErrorCodes

type Def string

const (
	None Def = "None"

	IncludeFileNotExist = "IncludeFileNotExist"
	FailOpenFile        = "FailOpenFile"
	FailScanFile        = "FailScanFile"

	BinFileSeekAtEnd          = "BinFileSeekAtEnd"
	BinFileSeekAfterEnd       = "BinFileSeekAfterEnd"
	BinFileReadBeyondFileSize = "BinFileReadBeyondFileSize"

	ParserEndOfInput      = "ParserEndOfInput"
	ParserUnexpectedToken = "ParserUnexpectedToken"

	OperationBadTokenAfter           = "OperationBadTokenAfter"
	OperationDirectiveUnknown        = "OperationDirectiveUnknown"
	OperationEmpty                   = "OperationEmpty"
	OperationLabelBadTokenAfter      = "OperationLabelBadTokenAfter"
	OperationLabelMissingColon       = "OperationLabelMissingColon"
	OperationLabeledDirectiveNoSpace = "OperationLabeledDirectiveNoSpace"
	OperationLabeledDirectiveUnknown = "OperationLabeledDirectiveUnknown"
	OperationUNKNOWN                 = "OperationUNKNOWN"

	OperandListStartingComma = "OperandListStartingComma"
	OperandListTooMany       = "OperandListTooMany"
	OperandListTooFew        = "OperandListTooFew"
	OperandBadCalleeName     = "OperandBadCalleeName"

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

	InstUnsupportedMode      = "InstUnsupportedMode"
	InstTokenAfterOperand    = "InstTokenAfterOperand"
	InstBadAccumMode         = "InstBadAccumMode"
	InstBadIndexValue        = "InstBadIndexValue"
	InstXYUnusableMode       = "InstXYUnusableMode"
	InstIndirectIndexMustBeX = "InstIndirectIndexMustBeX"
	InstIndirectIndexMustBeY = "InstIndirectIndexMustBeY"

	MacroInvokeDoubleCurlyBrace   = "MacroInvokeDoubleCurlyBrace"
	MacroInvokeUnclosedCurlyBrace = "MacroInvokeUnclosedCurlyBrace"
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
