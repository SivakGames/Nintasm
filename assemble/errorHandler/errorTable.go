package errorHandler

import (
	"fmt"
	enumErrorCodes "misc/nintasm/constants/enums/errorCodes"
	enumTerminalColors "misc/nintasm/constants/enums/terminalColors"
	"misc/nintasm/util"
)

type ErrorTableEntry struct {
	description string
	severity    enumErrorCodes.Severity
}

const SEVERITY_PREFIX = "SEVERITY_"

func newErrorTableEntry(severity enumErrorCodes.Severity, description string) ErrorTableEntry {
	return ErrorTableEntry{
		description: description,
		severity:    severity,
	}
}

func highlight(s string) string {
	return util.Colorize(s, enumTerminalColors.Yellow, false)
}
func coloredSymbol(s string) string {
	return util.Colorize(s, enumTerminalColors.LightCyan, false)
}
func coloredIdentifier(s string) string {
	return util.Colorize(s, enumTerminalColors.Magenta, false)
}
func coloredDirective(s string) string {
	return util.Colorize(s, enumTerminalColors.LightMagenta, false)
}
func coloredNumber(s string) string {
	return util.Colorize(s, enumTerminalColors.Yellow, false)
}
func coloredString(s string) string {
	return util.Colorize(s, enumTerminalColors.AnsiGreen, false)
}

var errorTable = map[enumErrorCodes.Def]ErrorTableEntry{
	enumErrorCodes.Other:      newErrorTableEntry(enumErrorCodes.Error, "%v"),
	enumErrorCodes.OtherFatal: newErrorTableEntry(enumErrorCodes.Fatal, "🛑 %v"),

	enumErrorCodes.IncludeFileNotExist: newErrorTableEntry(enumErrorCodes.Fatal, fmt.Sprintf("Source file %v does not exist!", coloredString("%v"))),
	enumErrorCodes.FailOpenFile:        newErrorTableEntry(enumErrorCodes.Fatal, fmt.Sprintf("Failed to open source file: %v", coloredString("%v"))),
	enumErrorCodes.FailScanFile:        newErrorTableEntry(enumErrorCodes.Fatal, "Failed to scan file!\n%v"),

	enumErrorCodes.BinFileSeekAtEnd:          newErrorTableEntry(enumErrorCodes.Fatal, "Seek value of %d is at the very end of file so no bytes can be read!"),
	enumErrorCodes.BinFileSeekAfterEnd:       newErrorTableEntry(enumErrorCodes.Fatal, "Seek value of %d goes beyond the size of file by %d byte(s)"),
	enumErrorCodes.BinFileReadBeyondFileSize: newErrorTableEntry(enumErrorCodes.Fatal, "Read value of %d goes beyond the size of file by %d byte(s)"),

	enumErrorCodes.OrphanedMultilineBrace:       newErrorTableEntry(enumErrorCodes.Error, "Orphaned multiline brace"),
	enumErrorCodes.TokenizerUnknownIllegalToken: newErrorTableEntry(enumErrorCodes.Error, "Unknown/Illegal token: %v"),

	enumErrorCodes.ParserEndOfInput:                  newErrorTableEntry(enumErrorCodes.Error, "Parsing error - Unexpected end of input!"),
	enumErrorCodes.OperandStatementEmpty:             newErrorTableEntry(enumErrorCodes.Error, "Parsing error - Operand is missing!"),
	enumErrorCodes.ParserUnexpectedToken:             newErrorTableEntry(enumErrorCodes.Error, fmt.Sprintf("Parsing error - Unexpected token: %v", highlight("%v"))),
	enumErrorCodes.ParserTemplateStringNotIdentifier: newErrorTableEntry(enumErrorCodes.Error, fmt.Sprintf("Template string resolved to %v which is not a valid identifier!", coloredSymbol("%v"))),
	enumErrorCodes.ParserTemplateStringEmpty:         newErrorTableEntry(enumErrorCodes.Error, fmt.Sprintf("Template string cannot be empty!")),

	enumErrorCodes.OperationUNKNOWN:                 newErrorTableEntry(enumErrorCodes.Error, "UNKNOWN OPERATION"),
	enumErrorCodes.OperationBadTokenAfter:           newErrorTableEntry(enumErrorCodes.Error, "ILLEGAL token(s) after operation: %v"),
	enumErrorCodes.OperationDirectiveUnknown:        newErrorTableEntry(enumErrorCodes.Error, fmt.Sprintf("Operation parsing failed - Unknown directive! .%v", coloredDirective("%v"))),
	enumErrorCodes.OperationEmpty:                   newErrorTableEntry(enumErrorCodes.Error, "Operation parsing failed - UNEXPECTED EMPTY OPERATION"),
	enumErrorCodes.OperationLabelBadTokenAfter:      newErrorTableEntry(enumErrorCodes.Error, "ILLEGAL token(s) after labeled operation: %v"),
	enumErrorCodes.OperationLabelMissingColon:       newErrorTableEntry(enumErrorCodes.Error, "Incomplete label/labeled operation - either needs colon or directive"),
	enumErrorCodes.OperationLabeledDirectiveNoSpace: newErrorTableEntry(enumErrorCodes.Error, "Operation parsing failed - Labeled directive needs space between label and directive"),
	enumErrorCodes.OperationLabeledDirectiveUnknown: newErrorTableEntry(enumErrorCodes.Error, "Operation parsing failed - Unknown labeled directive!"),

	enumErrorCodes.OperandListStartingComma:       newErrorTableEntry(enumErrorCodes.Error, "Operand list cannot start with a comma!"),
	enumErrorCodes.OperandListTooMany:             newErrorTableEntry(enumErrorCodes.Error, fmt.Sprintf("Too many operands for operation! Takes at most: %v", coloredNumber("%d"))),
	enumErrorCodes.OperandListTooFew:              newErrorTableEntry(enumErrorCodes.Error, fmt.Sprintf("Too few operands for operation! Needs at least: %v", coloredNumber("%d"))),
	enumErrorCodes.OperandBadCalleeName:           newErrorTableEntry(enumErrorCodes.Error, "Illegal functional callee name: %v"),
	enumErrorCodes.OperandMisplacedLiteral:        newErrorTableEntry(enumErrorCodes.Error, "Misplaced literal - %v"),
	enumErrorCodes.OperandMisplacedIdentifier:     newErrorTableEntry(enumErrorCodes.Error, "Misplaced identifier - %v"),
	enumErrorCodes.OperandMissingPrimaryExpr:      newErrorTableEntry(enumErrorCodes.Error, "No primary expression found"),
	enumErrorCodes.OperandPeriodMissingIdentifier: newErrorTableEntry(enumErrorCodes.Error, "Identifier must follow period!"),
	enumErrorCodes.OperandBadPrimaryExpr:          newErrorTableEntry(enumErrorCodes.Error, fmt.Sprintf("Bad primary expression - %v", highlight("%v"))),

	enumErrorCodes.NodeTypeNotBool:           newErrorTableEntry(enumErrorCodes.Error, "Value must be a boolean!"),
	enumErrorCodes.NodeTypeNotIdentifier:     newErrorTableEntry(enumErrorCodes.Error, "Value must be an identifier!"),
	enumErrorCodes.NodeTypeNotString:         newErrorTableEntry(enumErrorCodes.Error, "Value must be a string!"),
	enumErrorCodes.NodeTypeNotSubstitutionID: newErrorTableEntry(enumErrorCodes.Error, "Value must be a substitution ID!"),
	enumErrorCodes.NodeTypeNotArray:          newErrorTableEntry(enumErrorCodes.Error, "Value must be an array!"),

	enumErrorCodes.NodeTypeNotNumeric:   newErrorTableEntry(enumErrorCodes.Error, "Value must be numeric!"),
	enumErrorCodes.NodeValueNotPositive: newErrorTableEntry(enumErrorCodes.Error, "Value must be positive!"),
	enumErrorCodes.NodeValueNot8Bit:     newErrorTableEntry(enumErrorCodes.Error, "Value must be 8 bit!"),
	enumErrorCodes.NodeValueNot16Bit:    newErrorTableEntry(enumErrorCodes.Error, "Value must be 16 bit!"),
	enumErrorCodes.NodeValueNotPowerOf2: newErrorTableEntry(enumErrorCodes.Error, "Value must be a power of 2!"),

	enumErrorCodes.NodeValueNotGT:        newErrorTableEntry(enumErrorCodes.Error, fmt.Sprintf("Value must be %v %v", coloredSymbol(">"), coloredNumber("%d"))),
	enumErrorCodes.NodeValueNotLT:        newErrorTableEntry(enumErrorCodes.Error, fmt.Sprintf("Value must be %v %v", coloredSymbol("<"), coloredNumber("%d"))),
	enumErrorCodes.NodeValueNotGTE:       newErrorTableEntry(enumErrorCodes.Error, fmt.Sprintf("Value must be %v %v", coloredSymbol(">="), coloredNumber("%d"))),
	enumErrorCodes.NodeValueNotLTE:       newErrorTableEntry(enumErrorCodes.Error, fmt.Sprintf("Value must be %v %v", coloredSymbol("<="), coloredNumber("%d"))),
	enumErrorCodes.NodeValueNotGTEandLTE: newErrorTableEntry(enumErrorCodes.Error, fmt.Sprintf("Value must be %v %v and %v %v", coloredSymbol(">="), coloredNumber("%d"), coloredSymbol("<="), coloredNumber("%d"))),

	enumErrorCodes.InterpreterComputedMemberNegativeIndex: newErrorTableEntry(enumErrorCodes.Error, fmt.Sprintf("Array index %v must be positive", coloredNumber("%d"))),
	enumErrorCodes.InterpreterComputedMemberIndexTooBig:   newErrorTableEntry(enumErrorCodes.Error, fmt.Sprintf("Array index %v is out of range! Max range is: %v", coloredNumber("%d"), coloredNumber("%d"))),

	enumErrorCodes.InvalidValueAlias: newErrorTableEntry(enumErrorCodes.Error, "Invalid value alias!"),
	enumErrorCodes.UnacceptableAlias: newErrorTableEntry(enumErrorCodes.Error, "Unacceptable value alias!"),

	enumErrorCodes.INESValueAlreadySet: newErrorTableEntry(enumErrorCodes.Error, "%v value has already been set!!!"),

	enumErrorCodes.BankSizeUneven:      newErrorTableEntry(enumErrorCodes.Fatal, "Bank size is not evenly distributable"),
	enumErrorCodes.BankOverflow:        newErrorTableEntry(enumErrorCodes.Fatal, "Bank will overflow by: %d byte(s) here"),
	enumErrorCodes.BankNotSequential:   newErrorTableEntry(enumErrorCodes.Fatal, "Bank declarations must be sequentially incrementing"),
	enumErrorCodes.BankNumberTooHigh:   newErrorTableEntry(enumErrorCodes.Fatal, "Too high of a bank number!"),
	enumErrorCodes.OrgTooSmall:         newErrorTableEntry(enumErrorCodes.Error, "ORG is too small! Attempted: %d / Minimum Allowed: %d"),
	enumErrorCodes.OrgTooBig:           newErrorTableEntry(enumErrorCodes.Error, "ORG is too big! Attempted: %d / Max Allowed: %d"),
	enumErrorCodes.OrgLTProgramCounter: newErrorTableEntry(enumErrorCodes.Error, "Cannot set ORG to a value less than where the program counter currently is!\nThis would overwrite data!\n Attempted: %d / Currently at: %d"),

	enumErrorCodes.InstUnsupportedMode:      newErrorTableEntry(enumErrorCodes.Error, fmt.Sprintf("%v mode is not supported by instruction!", coloredNumber("%v"))),
	enumErrorCodes.InstTokenAfterOperand:    newErrorTableEntry(enumErrorCodes.Error, "No more tokens can follow this instruction's operands! %v"),
	enumErrorCodes.InstBadAccumMode:         newErrorTableEntry(enumErrorCodes.Error, "No tokens can follow A for accumulator mode."),
	enumErrorCodes.InstXYUnusableMode:       newErrorTableEntry(enumErrorCodes.Error, "X or Y indexes cannot be used with target mode"),
	enumErrorCodes.InstBadIndexValue:        newErrorTableEntry(enumErrorCodes.Error, "Bad/Unknown index value! %v"),
	enumErrorCodes.InstIndirectIndexMustBeX: newErrorTableEntry(enumErrorCodes.Error, "Must use X index for this kind of indirect addressing"),
	enumErrorCodes.InstIndirectIndexMustBeY: newErrorTableEntry(enumErrorCodes.Error, "Must use Y index for this kind of indirect addressing"),

	enumErrorCodes.DirectiveUnopenedEndBlock:  newErrorTableEntry(enumErrorCodes.Error, "%v with no opening operation found!"),
	enumErrorCodes.DirectiveUnmatchedEndBlock: newErrorTableEntry(enumErrorCodes.Error, fmt.Sprintf("Non-matching closing block for parent %v operation", coloredIdentifier("%v"))),
	enumErrorCodes.DirectiveNestedLabelBlock:  newErrorTableEntry(enumErrorCodes.Error, "Cannot define a labeled block when in another block statement!"),

	enumErrorCodes.MacroNotExist:                  newErrorTableEntry(enumErrorCodes.Error, fmt.Sprintf("Specified macro %v doesn't exist!", coloredIdentifier("%v"))),
	enumErrorCodes.MacroInvokeDoubleCurlyBrace:    newErrorTableEntry(enumErrorCodes.Error, "Macro invoking error - Must close curly brace before opening another!"),
	enumErrorCodes.MacroInvokeUnclosedCurlyBrace:  newErrorTableEntry(enumErrorCodes.Error, "Macro invoking error - Unclosed curly brace!"),
	enumErrorCodes.MacroInvokeTooFewArgs:          newErrorTableEntry(enumErrorCodes.Error, fmt.Sprintf("Macro invoking error - Too few arguments! Received: %v / Min needed: %v", coloredNumber("%d"), coloredNumber("%d"))),
	enumErrorCodes.MacroInvokeTooManyArgs:         newErrorTableEntry(enumErrorCodes.Error, fmt.Sprintf("Macro invoking error - Too many arguments! Received: %v / Max allowed: %v", coloredNumber("%d"), coloredNumber("%d"))),
	enumErrorCodes.MacroMisplacedExitMacro:        newErrorTableEntry(enumErrorCodes.Error, "Exit macro without actually being in a macro..."),
	enumErrorCodes.MacroSubstitutionAlreadyExists: newErrorTableEntry(enumErrorCodes.Error, fmt.Sprintf("Substition value %v already exists!", coloredSymbol("%v"))),
	enumErrorCodes.SubstitutionNotLabelLike:       newErrorTableEntry(enumErrorCodes.Error, fmt.Sprintf("Substition value %v must start with a letter or underscore after the backslash and be followed by numbers, underscores, or numbers", coloredSymbol("%v"))),

	enumErrorCodes.AssignmentMissingOperand:      newErrorTableEntry(enumErrorCodes.Error, "Missing operand for assignment!"),
	enumErrorCodes.AssignmentLocalNotInNamespace: newErrorTableEntry(enumErrorCodes.Error, "Cannot use local assignment outside of namespace"),
	enumErrorCodes.AssignmentNamespaceNotLocal:   newErrorTableEntry(enumErrorCodes.Error, "Must use local assignment within namespace"),

	enumErrorCodes.ReassignmentIdentifierNotFound:  newErrorTableEntry(enumErrorCodes.Fatal, fmt.Sprintf("Reassigned identifier %v not found and MUST be resolved to proceed", coloredIdentifier("%v"))),
	enumErrorCodes.ReassignmentIdentifierNotSymbol: newErrorTableEntry(enumErrorCodes.Error, fmt.Sprintf("Reassigned identifier %v must be a symbol. Got: %v", coloredIdentifier("%v"), coloredSymbol("%v"))),

	enumErrorCodes.RemovedIdentifierNotFound:  newErrorTableEntry(enumErrorCodes.Warning, fmt.Sprintf("Symbol %v not found so deleting nothing...", coloredIdentifier("%v"))),
	enumErrorCodes.RemovedIdentifierWrongType: newErrorTableEntry(enumErrorCodes.Error, fmt.Sprintf("Symbol %v is wrong type for deletion. Got: %v / Wanted: %v", coloredIdentifier("%v"), coloredSymbol("%v"), coloredSymbol("%v"))),

	enumErrorCodes.IfStatementElseIfNoParentIf:       newErrorTableEntry(enumErrorCodes.Error, "Cannot use elseif without parent if"),
	enumErrorCodes.IfStatementElseIfAfterElse:        newErrorTableEntry(enumErrorCodes.Error, "Cannot have elseif after else"),
	enumErrorCodes.IfStatementDuplicateElse:          newErrorTableEntry(enumErrorCodes.Error, "Can only have 1 else in this block"),
	enumErrorCodes.SwitchStatementBadOperand:         newErrorTableEntry(enumErrorCodes.Error, "Switch operand type must be or resolve to a number or string"),
	enumErrorCodes.SwitchStatementCaseAfterDefault:   newErrorTableEntry(enumErrorCodes.Error, "Cannot have case after default"),
	enumErrorCodes.SwitchStatementMismatchedCaseType: newErrorTableEntry(enumErrorCodes.Error, "Case's type mismatched from original switch's type"),
	enumErrorCodes.SwitchStatementDuplicateDefault:   newErrorTableEntry(enumErrorCodes.Error, "Cannot only have 1 default in this block"),
	enumErrorCodes.CaseNoSwitch:                      newErrorTableEntry(enumErrorCodes.Error, "Cannot have case outside a switch"),
	enumErrorCodes.DefaultNoSwitch:                   newErrorTableEntry(enumErrorCodes.Error, "Cannot have default outside a switch"),

	enumErrorCodes.NamespaceNotExist:            newErrorTableEntry(enumErrorCodes.Error, fmt.Sprintf("Namespace %v does not exist!", coloredIdentifier("%v"))),
	enumErrorCodes.NamespaceToValuesNotResolved: newErrorTableEntry(enumErrorCodes.Error, "Namespace value %v is not resolved and cannot be converted to a value"),
	enumErrorCodes.GNSIsourceIsLocal:            newErrorTableEntry(enumErrorCodes.Error, fmt.Sprintf("GNSI target %v is a local label and cannot be used", coloredIdentifier("%v"))),
	enumErrorCodes.GNSIparentNotFound:           newErrorTableEntry(enumErrorCodes.Error, fmt.Sprintf("GNSI target %v not found", coloredIdentifier("%v"))),
	enumErrorCodes.GNSIparentNotLabel:           newErrorTableEntry(enumErrorCodes.Error, fmt.Sprintf("GNSI target %v is not a label!", coloredIdentifier("%v"))),
	enumErrorCodes.GNSIparentHasNoLocals:        newErrorTableEntry(enumErrorCodes.Error, fmt.Sprintf("GNSI target %v has no local labels!", coloredIdentifier("%v"))),
	enumErrorCodes.BytesWithinLabelNoEnd:        newErrorTableEntry(enumErrorCodes.Error, "No ending label found so range cannot be calculated"),

	enumErrorCodes.SubFuncStartTooBig:        newErrorTableEntry(enumErrorCodes.Error, "Starting index is too big! Highest index possible: %d"),
	enumErrorCodes.SubFuncEndTooBig:          newErrorTableEntry(enumErrorCodes.Error, "Ending index is too big! Highest index possible: %d"),
	enumErrorCodes.SubFuncEndBiggerThanStart: newErrorTableEntry(enumErrorCodes.Error, "End index must be greater than start!"),

	enumErrorCodes.IndexFuncValueNotFound: newErrorTableEntry(enumErrorCodes.Error, fmt.Sprintf("Value %v not found for index function", coloredIdentifier("%v"))),

	enumErrorCodes.CharMapNoneDefined:     newErrorTableEntry(enumErrorCodes.Error, "No character maps have been defined!"),
	enumErrorCodes.CharMapNotExist:        newErrorTableEntry(enumErrorCodes.Error, "Specified charmap doesn't exist!"),
	enumErrorCodes.CharMapDuplicateKey:    newErrorTableEntry(enumErrorCodes.Error, "Character %c has already been defined in current map %v"),
	enumErrorCodes.ToCharMapUndefChar:     newErrorTableEntry(enumErrorCodes.Error, "Character `%c` is not defined in currently used charmap %v"),
	enumErrorCodes.DefCharTooLong:         newErrorTableEntry(enumErrorCodes.Error, "Can't map %v - Character definition must be 1 character long!"),
	enumErrorCodes.DefCharRangeEndSmaller: newErrorTableEntry(enumErrorCodes.Error, "End value of defined character range value must be greater than start value"),

	enumErrorCodes.ExprMapNoneDefined:  newErrorTableEntry(enumErrorCodes.Error, "No expression maps have been defined!"),
	enumErrorCodes.ExprMapNotExist:     newErrorTableEntry(enumErrorCodes.Error, "Specified exprmap doesn't exist!"),
	enumErrorCodes.ExprMapDuplicateKey: newErrorTableEntry(enumErrorCodes.Error, fmt.Sprintf("Expression %v has already been defined in current map", coloredSymbol("`%v`"))),
	enumErrorCodes.ExprMapUndefExpr:    newErrorTableEntry(enumErrorCodes.Error, fmt.Sprintf("Expression %v is not defined in currently used exprmap %v", coloredSymbol("`%v`"), coloredIdentifier("%v"))),
	enumErrorCodes.ToExprMapUndefExpr:  newErrorTableEntry(enumErrorCodes.Error, "Expression `%v` is not defined in currently used exprmap"),

	enumErrorCodes.InterpreterNoParentLabel:            newErrorTableEntry(enumErrorCodes.Error, "Cannot use operation! No parent label!"),
	enumErrorCodes.InterpreterNoLocalLabel:             newErrorTableEntry(enumErrorCodes.Error, "Cannot use operation! No local label!"),
	enumErrorCodes.InterpreterUnaryNotBoolean:          newErrorTableEntry(enumErrorCodes.Error, "Unary expression must be boolean: %v %v"),
	enumErrorCodes.InterpreterUnaryNotNumeric:          newErrorTableEntry(enumErrorCodes.Error, "Unary expression must be numeric: %v %v"),
	enumErrorCodes.InterpreterBinaryMismatchedTypes:    newErrorTableEntry(enumErrorCodes.Error, fmt.Sprintf("Binary expression types mismatched: %v %v %v", "%v", coloredSymbol("%v"), "%v")),
	enumErrorCodes.InterpreterFuncUndefined:            newErrorTableEntry(enumErrorCodes.Error, "Function %v is undefined!"),
	enumErrorCodes.InterpreterFuncTooFewArgs:           newErrorTableEntry(enumErrorCodes.Error, "Too few arguments for %v function!"),
	enumErrorCodes.InterpreterFuncTooManyArgs:          newErrorTableEntry(enumErrorCodes.Error, "Too many arguments for %v function!"),
	enumErrorCodes.InterpreterFuncArgWrongType:         newErrorTableEntry(enumErrorCodes.Error, fmt.Sprintf("Argument is wrong type for %v function", coloredSymbol("%v"))),
	enumErrorCodes.InterpreterAlreadyDefined:           newErrorTableEntry(enumErrorCodes.Error, fmt.Sprintf("Symbol %v has been previously defined! (Defined as %v)", coloredSymbol("%v"), coloredNumber("%v"))),
	enumErrorCodes.InterpreterSymbolNotFound:           newErrorTableEntry(enumErrorCodes.Error, fmt.Sprintf("Symbol %v was not found and must be resolved!", coloredSymbol("%v"))),
	enumErrorCodes.InterpreterIdentifierNotValueSymbol: newErrorTableEntry(enumErrorCodes.Error, "Identifier %v is not usable as a numeric value!"),

	enumErrorCodes.BlockIsEmpty:                newErrorTableEntry(enumErrorCodes.Warning, "Block is empty..."),
	enumErrorCodes.BlockOpUncapturableByParent: newErrorTableEntry(enumErrorCodes.Error, "%v - This operation is uncapturable by block"),

	enumErrorCodes.ResolvedValueNot8Bit:         newErrorTableEntry(enumErrorCodes.Error, fmt.Sprintf("Operand must resolve to an 8 bit value! Got: %v", coloredNumber("%d"))),
	enumErrorCodes.ResolvedValueNot16Bit:        newErrorTableEntry(enumErrorCodes.Error, fmt.Sprintf("Operand must resolve to a 16 bit value! Got: %v", coloredNumber("%d"))),
	enumErrorCodes.ResolvedValueIsImplicitUndef: newErrorTableEntry(enumErrorCodes.Error, "Operand has resolved to an implicitly undefined value!"),
	enumErrorCodes.ResolvedValueIsStringForInst: newErrorTableEntry(enumErrorCodes.Warning, "Operand has resolved as a string and not recommended for target instruction node"),
	enumErrorCodes.ResolvedValueIsBool:          newErrorTableEntry(enumErrorCodes.Warning, "Operand has resolved as bool; Will be converted to: %d"),
	enumErrorCodes.ResolvedValue16BitBool:       newErrorTableEntry(enumErrorCodes.Error, "Boolean value cannot be used as a 16 bit operand"),
	enumErrorCodes.ResolvedValueMultiByteChar:   newErrorTableEntry(enumErrorCodes.Warning, "Character %v encoding requires more than a single byte. Using %d bytes"),
	enumErrorCodes.ResolvedValue16BitString:     newErrorTableEntry(enumErrorCodes.Error, "String value cannot be used as a 16 bit operand"),
	enumErrorCodes.ResolvedValueTooBig:          newErrorTableEntry(enumErrorCodes.Error, fmt.Sprintf("Resolved value byte size of %v too large for target byte size of %v", coloredNumber("%d"), coloredNumber("%d"))),
	enumErrorCodes.ResolvedValueBranchTooFar:    newErrorTableEntry(enumErrorCodes.Error, "Branch destination of %d bytes away is too far!"),
	enumErrorCodes.ResolvedValueNotInt:          newErrorTableEntry(enumErrorCodes.Warning, "Resolved numeric value of %.4f is not an int. Truncating to %d"),
	enumErrorCodes.ResolveImpossible:            newErrorTableEntry(enumErrorCodes.Error, "Symbol cannot be resolved!"),
	enumErrorCodes.ResolveDeadlock:              newErrorTableEntry(enumErrorCodes.Fatal, "🛑 Resolve Deadlock - Symbols are set up in a way in which they will NEVER resolve!"),

	enumErrorCodes.RsNotSet: newErrorTableEntry(enumErrorCodes.Error, "RS has not yet been set!"),
}

var errorHintTable = map[enumErrorCodes.Def]string{
	enumErrorCodes.InstUnsupportedMode: fmt.Sprintf("Supported modes are: %v", coloredNumber("%v")),
}
