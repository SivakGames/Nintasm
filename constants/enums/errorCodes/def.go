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

	TokenizerUnknownIllegalToken = "TokenizerUnknownIllegalToken"

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

	OperandListStartingComma       = "OperandListStartingComma"
	OperandListTooMany             = "OperandListTooMany"
	OperandListTooFew              = "OperandListTooFew"
	OperandBadCalleeName           = "OperandBadCalleeName"
	OperandMisplacedLiteral        = "OperandMisplacedLiteral"
	OperandMisplacedIdentifier     = "OperandMisplacedIdentifier"
	OperandMissingPrimaryExpr      = "OperandMissingPrimaryExpr"
	OperandPeriodMissingIdentifier = "OperandPeriodMissingIdentifier"
	OperandBadPrimaryExpr          = "OperandBadPrimaryExpr"

	NodeTypeNotBool           = "NodeTypeNotBool"
	NodeTypeNotString         = "NodeTypeNotString"
	NodeTypeNotIdentifier     = "NodeTypeNotIdentifier"
	NodeTypeNotSubstitutionID = "NodeTypeNotSubstitutionID"

	NodeTypeNotNumeric   = "NodeTypeNotNumeric"
	NodeValueNotPositive = "NodeValueNotPositive"
	NodeValueNot8Bit     = "NodeValueNot8Bit"
	NodeValueNot16Bit    = "NodeValueNot16Bit"
	NodeValueNotPowerOf2 = "NodeValueNotPowerOf2"

	NodeValueNotGT        = "NodeValueNotGT"
	NodeValueNotLT        = "NodeValueNotLT"
	NodeValueNotGTE       = "NodeValueNotGTE"
	NodeValueNotLTE       = "NodeValueNotLTE"
	NodeValueNotGTEandLTE = "NodeValueNotGTEandLTE"

	InvalidValueAlias = "InvalidValueAlias"
	UnacceptableAlias = "UnacceptableAlias"

	INESValueAlreadySet = "INESValueAlreadySet"

	BankSizeUneven       = "BankSizeUneven"
	BankNotSequential    = "BankNotSequential"
	BankNumberTooHigh    = "BankNumberTooHigh"
	BankOverflow         = "BankOverflow"
	OrgTooSmall          = "OrgTooSmall"
	OrgTooBig            = "OrgTooBig"
	OrgLTEProgramCounter = "OrgLTEProgramCounter"

	InstUnsupportedMode      = "InstUnsupportedMode"
	InstTokenAfterOperand    = "InstTokenAfterOperand"
	InstBadAccumMode         = "InstBadAccumMode"
	InstBadIndexValue        = "InstBadIndexValue"
	InstXYUnusableMode       = "InstXYUnusableMode"
	InstIndirectIndexMustBeX = "InstIndirectIndexMustBeX"
	InstIndirectIndexMustBeY = "InstIndirectIndexMustBeY"

	DirectiveUnopenedEndBlock  = "DirectiveUnopenedEndBlock"
	DirectiveUnmatchedEndBlock = "DirectiveUnmatchedEndBlock"
	DirectiveNestedLabelBlock  = "DirectiveNestedLabelBlock"

	MacroInvokeDoubleCurlyBrace   = "MacroInvokeDoubleCurlyBrace"
	MacroInvokeUnclosedCurlyBrace = "MacroInvokeUnclosedCurlyBrace"

	AssignmentMissingOperand = "AssignmentMissingOperand"

	IfStatementElseIfAfterElse = "IfStatementElseIfAfterElse"
	IfStatementDuplicateElse   = "IfStatementDuplicateElse"

	CharMapNoneDefined     = "CharMapNoneDefined"
	CharMapNotExist        = "CharMapNotExist"
	CharMapDuplicateKey    = "CharMapDuplicateKey"
	ToCharMapUndefChar     = "ToCharMapUndefChar"
	DefCharTooLong         = "DefCharTooLong"
	DefCharRangeEndSmaller = "DefCharRangeEndSmaller"

	ExprMapNoneDefined  = "ExprMapNoneDefined"
	ExprMapNotExist     = "ExprMapNotExist"
	ExprMapDuplicateKey = "ExprMapDuplicateKey"
	ToExprMapUndefExpr  = "ToExprMapUndefExpr"

	InterpreterNoParentLabel    = "InterpreterNoParentLabel"
	InterpreterFuncTooFewArgs   = "InterpreterFuncTooFewArgs"
	InterpreterFuncTooManyArgs  = "InterpreterFuncTooManyArgs"
	InterpreterFuncArgWrongType = "InterpreterFuncArgWrongType"
	InterpreterAlreadyDefined   = "InterpreterAlreadyDefined"
	InterpreterSymbolNotFound   = "InterpreterSymbolNotFound"

	BlockOpUncapturableByParent = "BlockOpUncapturableByParent"
	ResolvedValueNot8Bit        = "ResolvedValueNot8Bit"
	ResolvedValueNot16Bit       = "ResolvedValueNot16Bit"
	ResolvedValueIsBool         = "ResolvedValueIsBool"
	ResolvedValue16BitBool      = "ResolvedValue16BitBool"
	ResolvedValueMultiByteChar  = "ResolvedValueMultiByteChar"
	ResolvedValue16BitString    = "ResolvedValue16BitString"

	RsNotSet = "RsNotSet"
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
