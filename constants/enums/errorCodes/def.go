package enumErrorCodes

type Def string

const (
	None Def = "None"

	Other      = "Other"
	OtherFatal = "OtherFatal"

	IncludeFileNotExist = "IncludeFileNotExist"
	FailOpenFile        = "FailOpenFile"
	FailScanFile        = "FailScanFile"

	BinFileSeekAtEnd          = "BinFileSeekAtEnd"
	BinFileSeekAfterEnd       = "BinFileSeekAfterEnd"
	BinFileReadBeyondFileSize = "BinFileReadBeyondFileSize"

	TokenizerUnknownIllegalToken = "TokenizerUnknownIllegalToken"

	ParserEndOfInput      = "ParserEndOfInput"
	OperandStatementEmpty = "OperandStatementEmpty"
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

	BankSizeUneven      = "BankSizeUneven"
	BankNotSequential   = "BankNotSequential"
	BankNumberTooHigh   = "BankNumberTooHigh"
	BankOverflow        = "BankOverflow"
	OrgTooSmall         = "OrgTooSmall"
	OrgTooBig           = "OrgTooBig"
	OrgLTProgramCounter = "OrgLTProgramCounter"

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

	MacroNotExist                 = "MacroNotExist"
	MacroInvokeDoubleCurlyBrace   = "MacroInvokeDoubleCurlyBrace"
	MacroInvokeUnclosedCurlyBrace = "MacroInvokeUnclosedCurlyBrace"

	AssignmentMissingOperand      = "AssignmentMissingOperand"
	AssignmentLocalNotInNamespace = "AssignmentLocalNotInNamespace"
	AssignmentNamespaceNotLocal   = "AssignmentNamespaceNotLocal"

	IfStatementElseIfAfterElse        = "IfStatementElseIfAfterElse"
	IfStatementDuplicateElse          = "IfStatementDuplicateElse"
	SwitchStatementBadOperand         = "SwitchStatementBadOperand"
	SwitchStatementCaseAfterDefault   = "SwitchStatementCaseAfterDefault"
	SwitchStatementMismatchedCaseType = "SwitchStatementMismatchedCaseType"
	SwitchStatementDuplicateDefault   = "SwitchStatementDuplicateDefault"
	CaseNoSwitch                      = "CaseNoSwitch"
	DefaultNoSwitch                   = "DefaultNoSwitch"

	NamespaceNotExist            = "NamespaceNotExist"
	NamespaceToValuesNotResolved = "NamespaceToValuesNotResolved"

	CharMapNoneDefined     = "CharMapNoneDefined"
	CharMapNotExist        = "CharMapNotExist"
	CharMapDuplicateKey    = "CharMapDuplicateKey"
	ToCharMapUndefChar     = "ToCharMapUndefChar"
	DefCharTooLong         = "DefCharTooLong"
	DefCharRangeEndSmaller = "DefCharRangeEndSmaller"

	ExprMapNoneDefined  = "ExprMapNoneDefined"
	ExprMapNotExist     = "ExprMapNotExist"
	ExprMapDuplicateKey = "ExprMapDuplicateKey"
	ExprMapUndefExpr    = "ExprMapUndefExpr"
	ToExprMapUndefExpr  = "ToExprMapUndefExpr"

	InterpreterNoParentLabel            = "InterpreterNoParentLabel"
	InterpreterUnaryNotBoolean          = "InterpreterUnaryNotBoolean"
	InterpreterUnaryNotNumeric          = "InterpreterUnaryNotNumeric"
	InterpreterBinaryMismatchedTypes    = "InterpreterBinaryMismatchedTypes"
	InterpreterFuncUndefined            = "InterpreterFuncUndefined"
	InterpreterFuncTooFewArgs           = "InterpreterFuncTooFewArgs"
	InterpreterFuncTooManyArgs          = "InterpreterFuncTooManyArgs"
	InterpreterFuncArgWrongType         = "InterpreterFuncArgWrongType"
	InterpreterAlreadyDefined           = "InterpreterAlreadyDefined"
	InterpreterSymbolNotFound           = "InterpreterSymbolNotFound"
	InterpreterIdentifierNotValueSymbol = "InterpreterIdentifierNotValueSymbol"

	BlockIsEmpty                = "BlockIsEmpty"
	BlockOpUncapturableByParent = "BlockOpUncapturableByParent"

	ResolvedValueNot8Bit         = "ResolvedValueNot8Bit"
	ResolvedValueNot16Bit        = "ResolvedValueNot16Bit"
	ResolvedValueIsStringForInst = "ResolvedValueIsStringForInst"
	ResolvedValueIsBool          = "ResolvedValueIsBool"
	ResolvedValue16BitBool       = "ResolvedValue16BitBool"
	ResolvedValueMultiByteChar   = "ResolvedValueMultiByteChar"
	ResolvedValue16BitString     = "ResolvedValue16BitString"
	ResolvedValueTooBig          = "ResolvedValueTooBig"
	ResolveImpossible            = "ResolveImpossible"
	ResolveDeadlock              = "ResolveDeadlock"

	RsNotSet = "RsNotSet"
)

type Severity int

const (
	// Will just be a warning - stops nothing
	Warning Severity = iota
	// Will try re-resolving after main pass is done
	UnresolvedIdentifier
	// General error. Will not stop build, but will fail and not be re-evaluated
	Error
	// Stop building entirely
	Fatal
)
