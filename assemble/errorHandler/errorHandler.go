package errorHandler

import (
	"errors"
	"fmt"
	"misc/nintasm/assemble/fileStack"
	enumErrorCodes "misc/nintasm/constants/enums/errorCodes"
	"misc/nintasm/util"
	"strconv"
	"strings"
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

func coloredSymbol(s string) string {
	return util.Colorize(s, "lightcyan", false)
}
func coloredNumber(s string) string {
	return util.Colorize(s, "yellow", false)
}

var errorTable = map[enumErrorCodes.Def]ErrorTableEntry{
	enumErrorCodes.IncludeFileNotExist: newErrorTableEntry(enumErrorCodes.Fatal, "Source file \x1b[92m%v\x1b[0m does not exist!"),
	enumErrorCodes.FailOpenFile:        newErrorTableEntry(enumErrorCodes.Fatal, "Failed to open source file: %v"),
	enumErrorCodes.FailScanFile:        newErrorTableEntry(enumErrorCodes.Fatal, "Failed to scan file!\n%v"),

	enumErrorCodes.BinFileSeekAtEnd:          newErrorTableEntry(enumErrorCodes.Fatal, "Seek value of %d is at the very end of file so no bytes can be read!"),
	enumErrorCodes.BinFileSeekAfterEnd:       newErrorTableEntry(enumErrorCodes.Fatal, "Seek value of %d goes beyond the size of file by %d byte(s)"),
	enumErrorCodes.BinFileReadBeyondFileSize: newErrorTableEntry(enumErrorCodes.Fatal, "Read value of %d goes beyond the size of file by %d byte(s)"),

	enumErrorCodes.TokenizerUnknownIllegalToken: newErrorTableEntry(enumErrorCodes.Error, "Unknown/Illegal token: %v"),

	enumErrorCodes.ParserEndOfInput:      newErrorTableEntry(enumErrorCodes.Error, "Parsing error - Unexpected end of input!"),
	enumErrorCodes.ParserUnexpectedToken: newErrorTableEntry(enumErrorCodes.Error, fmt.Sprintf("Parsing error - Unexpected token: %v", coloredSymbol("%v"))),

	enumErrorCodes.OperationUNKNOWN:                 newErrorTableEntry(enumErrorCodes.Error, "UNKNOWN OPERATION"),
	enumErrorCodes.OperationBadTokenAfter:           newErrorTableEntry(enumErrorCodes.Error, "ILLEGAL token(s) after operation: %v"),
	enumErrorCodes.OperationDirectiveUnknown:        newErrorTableEntry(enumErrorCodes.Error, "Operation parsing failed - Unknown directive!"),
	enumErrorCodes.OperationEmpty:                   newErrorTableEntry(enumErrorCodes.Error, "Operation parsing failed - UNEXPECTED EMPTY OPERATION"),
	enumErrorCodes.OperationLabelBadTokenAfter:      newErrorTableEntry(enumErrorCodes.Error, "ILLEGAL token(s) after labeled operation: %v"),
	enumErrorCodes.OperationLabelMissingColon:       newErrorTableEntry(enumErrorCodes.Error, "Incomplete label/labeled operation - either needs colon or directive"),
	enumErrorCodes.OperationLabeledDirectiveNoSpace: newErrorTableEntry(enumErrorCodes.Error, "Operation parsing failed - Labeled directive needs space between label and directive"),
	enumErrorCodes.OperationLabeledDirectiveUnknown: newErrorTableEntry(enumErrorCodes.Error, "Operation parsing failed - Unknown labeled directive!"),

	enumErrorCodes.OperandListStartingComma:       newErrorTableEntry(enumErrorCodes.Error, "Operand list cannot start with a comma!"),
	enumErrorCodes.OperandListTooMany:             newErrorTableEntry(enumErrorCodes.Error, "Too many operands for operation! Takes at most: %d"),
	enumErrorCodes.OperandListTooFew:              newErrorTableEntry(enumErrorCodes.Error, "Too few operands for operation! Needs at least: %d"),
	enumErrorCodes.OperandBadCalleeName:           newErrorTableEntry(enumErrorCodes.Error, "Illegal functional callee name: %v"),
	enumErrorCodes.OperandMisplacedLiteral:        newErrorTableEntry(enumErrorCodes.Error, "Misplaced literal - %v"),
	enumErrorCodes.OperandMisplacedIdentifier:     newErrorTableEntry(enumErrorCodes.Error, "Misplaced identifier - %v"),
	enumErrorCodes.OperandMissingPrimaryExpr:      newErrorTableEntry(enumErrorCodes.Error, "No primary expression found"),
	enumErrorCodes.OperandPeriodMissingIdentifier: newErrorTableEntry(enumErrorCodes.Error, "Identifier must follow period!"),
	enumErrorCodes.OperandBadPrimaryExpr:          newErrorTableEntry(enumErrorCodes.Error, "Bad primary expression - %v"),

	enumErrorCodes.NodeTypeNotBool:           newErrorTableEntry(enumErrorCodes.Error, "Value must be a boolean!"),
	enumErrorCodes.NodeTypeNotIdentifier:     newErrorTableEntry(enumErrorCodes.Error, "Value must be an identifier!"),
	enumErrorCodes.NodeTypeNotString:         newErrorTableEntry(enumErrorCodes.Error, "Value must be a string!"),
	enumErrorCodes.NodeTypeNotSubstitutionID: newErrorTableEntry(enumErrorCodes.Error, "Value must be a substitution ID!"),

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

	enumErrorCodes.InvalidValueAlias: newErrorTableEntry(enumErrorCodes.Error, "Invalid value alias!"),
	enumErrorCodes.UnacceptableAlias: newErrorTableEntry(enumErrorCodes.Error, "Unacceptable value alias!"),

	enumErrorCodes.INESValueAlreadySet: newErrorTableEntry(enumErrorCodes.Error, "%v value has already been set!!!"),

	enumErrorCodes.BankSizeUneven:       newErrorTableEntry(enumErrorCodes.Fatal, "Bank size is not evenly distributable"),
	enumErrorCodes.BankOverflow:         newErrorTableEntry(enumErrorCodes.Fatal, "Bank will overflow by: %d byte(s) here"),
	enumErrorCodes.BankNotSequential:    newErrorTableEntry(enumErrorCodes.Fatal, "Bank declarations must be sequentially incrementing"),
	enumErrorCodes.BankNumberTooHigh:    newErrorTableEntry(enumErrorCodes.Fatal, "Too high of a bank number!"),
	enumErrorCodes.OrgTooSmall:          newErrorTableEntry(enumErrorCodes.Error, "ORG is too small! Attempted: %d / Minimum Allowed: %d"),
	enumErrorCodes.OrgTooBig:            newErrorTableEntry(enumErrorCodes.Error, "ORG is too big! Attempted: %d / Max Allowed: %d"),
	enumErrorCodes.OrgLTEProgramCounter: newErrorTableEntry(enumErrorCodes.Error, "Cannot set ORG to a value less than where the program counter currently is!\nThis would overwrite data!\n Attempted: %d / Currently at: %d"),

	enumErrorCodes.InstUnsupportedMode:      newErrorTableEntry(enumErrorCodes.Error, "Mode is not supported by instruction!"),
	enumErrorCodes.InstTokenAfterOperand:    newErrorTableEntry(enumErrorCodes.Error, "No more tokens can follow this instruction's operands! %v"),
	enumErrorCodes.InstBadAccumMode:         newErrorTableEntry(enumErrorCodes.Error, "No tokens can follow A for accumulator mode."),
	enumErrorCodes.InstXYUnusableMode:       newErrorTableEntry(enumErrorCodes.Error, "X or Y indexes cannot be used with target mode"),
	enumErrorCodes.InstBadIndexValue:        newErrorTableEntry(enumErrorCodes.Error, "Bad/Unknown index value! %v"),
	enumErrorCodes.InstIndirectIndexMustBeX: newErrorTableEntry(enumErrorCodes.Error, "Must use X index for this kind of indirect addressing"),
	enumErrorCodes.InstIndirectIndexMustBeY: newErrorTableEntry(enumErrorCodes.Error, "Must use Y index for this kind of indirect addressing"),

	enumErrorCodes.DirectiveUnopenedEndBlock:  newErrorTableEntry(enumErrorCodes.Error, "%v with no opening operation found!"),
	enumErrorCodes.DirectiveUnmatchedEndBlock: newErrorTableEntry(enumErrorCodes.Error, "Non-matching closing block with parent operation, %v"),
	enumErrorCodes.DirectiveNestedLabelBlock:  newErrorTableEntry(enumErrorCodes.Error, "Cannot define a labeled block when in another block statement!"),

	enumErrorCodes.MacroNotExist:                 newErrorTableEntry(enumErrorCodes.Error, "Specified macro %v doesn't exist!"),
	enumErrorCodes.MacroInvokeDoubleCurlyBrace:   newErrorTableEntry(enumErrorCodes.Error, "Macro invoking error - Must close curly brace before opening another!"),
	enumErrorCodes.MacroInvokeUnclosedCurlyBrace: newErrorTableEntry(enumErrorCodes.Error, "Macro invoking error - Unclosed curly brace!"),

	enumErrorCodes.AssignmentMissingOperand: newErrorTableEntry(enumErrorCodes.Error, "Missing operand for assignment!"),

	enumErrorCodes.IfStatementElseIfAfterElse: newErrorTableEntry(enumErrorCodes.Error, "Cannot have elseif after else"),
	enumErrorCodes.IfStatementDuplicateElse:   newErrorTableEntry(enumErrorCodes.Error, "Cannot only have 1 else in this block"),

	enumErrorCodes.CharMapNoneDefined:     newErrorTableEntry(enumErrorCodes.Error, "No character maps have been defined!"),
	enumErrorCodes.CharMapNotExist:        newErrorTableEntry(enumErrorCodes.Error, "Specified charmap doesn't exist!"),
	enumErrorCodes.CharMapDuplicateKey:    newErrorTableEntry(enumErrorCodes.Error, "Character %c has already been defined in current map"),
	enumErrorCodes.ToCharMapUndefChar:     newErrorTableEntry(enumErrorCodes.Error, "Character `%c` is not defined in currently used charmap %v"),
	enumErrorCodes.DefCharTooLong:         newErrorTableEntry(enumErrorCodes.Error, "Can't map %v - Character definition must be 1 character long!"),
	enumErrorCodes.DefCharRangeEndSmaller: newErrorTableEntry(enumErrorCodes.Error, "End value of defined character range value must be greater than start value"),

	enumErrorCodes.ExprMapNoneDefined:  newErrorTableEntry(enumErrorCodes.Error, "No expression maps have been defined!"),
	enumErrorCodes.ExprMapNotExist:     newErrorTableEntry(enumErrorCodes.Error, "Specified exprmap doesn't exist!"),
	enumErrorCodes.ExprMapDuplicateKey: newErrorTableEntry(enumErrorCodes.Error, "Expression %v has already been defined in current map"),
	enumErrorCodes.ToExprMapUndefExpr:  newErrorTableEntry(enumErrorCodes.Error, "Expression `%v` is not defined in currently used exprmap"),

	enumErrorCodes.InterpreterNoParentLabel:            newErrorTableEntry(enumErrorCodes.Error, "Cannot use operation! No parent label!"),
	enumErrorCodes.InterpreterBinaryMismatchedTypes:    newErrorTableEntry(enumErrorCodes.Error, "Binary expression types mismatched: %v %v %v"),
	enumErrorCodes.InterpreterUnaryNotNumeric:          newErrorTableEntry(enumErrorCodes.Error, "Unary expression must be numeric: %v %v"),
	enumErrorCodes.InterpreterFuncTooFewArgs:           newErrorTableEntry(enumErrorCodes.Error, "Too few arguments for function!"),
	enumErrorCodes.InterpreterFuncTooManyArgs:          newErrorTableEntry(enumErrorCodes.Error, "Too many arguments for function!"),
	enumErrorCodes.InterpreterFuncArgWrongType:         newErrorTableEntry(enumErrorCodes.Error, "Argument is wrong type"),
	enumErrorCodes.InterpreterAlreadyDefined:           newErrorTableEntry(enumErrorCodes.Error, "Symbol %v has been previously defined! (Defined as %v)"),
	enumErrorCodes.InterpreterSymbolNotFound:           newErrorTableEntry(enumErrorCodes.Error, "Symbol %v was not found and must be resolved!"),
	enumErrorCodes.InterpreterIdentifierNotValueSymbol: newErrorTableEntry(enumErrorCodes.Error, "Identifier %v is not usable as a numeric value!"),

	enumErrorCodes.BlockIsEmpty:                newErrorTableEntry(enumErrorCodes.Warning, "Block is empty..."),
	enumErrorCodes.BlockOpUncapturableByParent: newErrorTableEntry(enumErrorCodes.Error, "%v - This operation is uncapturable by block"),

	enumErrorCodes.ResolvedValueNot8Bit:       newErrorTableEntry(enumErrorCodes.Error, "Operand must resolve to an 8 bit value!"),
	enumErrorCodes.ResolvedValueNot16Bit:      newErrorTableEntry(enumErrorCodes.Error, "Operand must resolve to a 16 bit value!"),
	enumErrorCodes.ResolvedValueIsBool:        newErrorTableEntry(enumErrorCodes.Warning, "Operand has resolved as bool; Will be converted to: %d"),
	enumErrorCodes.ResolvedValue16BitBool:     newErrorTableEntry(enumErrorCodes.Error, "Boolean value cannot be used as a 16 bit operand"),
	enumErrorCodes.ResolvedValueMultiByteChar: newErrorTableEntry(enumErrorCodes.Warning, "Character %v encoding requires more than a single byte. Using %d bytes"),
	enumErrorCodes.ResolvedValue16BitString:   newErrorTableEntry(enumErrorCodes.Error, "String value cannot be used as a 16 bit operand"),

	enumErrorCodes.RsNotSet: newErrorTableEntry(enumErrorCodes.Error, "RS has not yet been set!"),
}

// ++++++++++++++++++++++++++++++++++++++++

type ErrorEntry struct {
	code        enumErrorCodes.Def
	lineNumber  uint
	lineContent string
	fileName    string
	message     string
	hint        string
	severity    enumErrorCodes.Severity
}

func NewErrorEntry(code enumErrorCodes.Def, message string, severity enumErrorCodes.Severity) ErrorEntry {
	fileData := fileStack.GetTopOfFileStack()

	if fileData != nil {
		return ErrorEntry{
			code:        code,
			message:     message,
			fileName:    fileData.FileName,
			lineNumber:  fileData.CurrentLineNumber,
			lineContent: fileData.ProcessedLines[fileData.CurrentLineNumber-1],
			severity:    severity,
		}
	}
	return ErrorEntry{
		code:        code,
		message:     message,
		fileName:    "NO FILE",
		lineNumber:  0,
		lineContent: "",
		severity:    severity,
	}
}

// If severity is >= threshold it should stop propagating up
func CheckErrorContinuesUpwardPropagation(err error, threshold enumErrorCodes.Severity) error {
	severityValue := err.Error()

	modded, ok := strings.CutPrefix(severityValue, SEVERITY_PREFIX)
	if ok {
		severityAmt, err := strconv.Atoi(modded)
		if err != nil {
			return err
		}
		if severityAmt <= int(threshold) {
			return nil
		}
	}
	return err
}

/*

░ >> D:\Emulate\NES\Disassemblies\Lolo 3\prg\fixed.6502
▓   2711   .include "prg/music-engine/dpcm-samples.6502a"
▓  FATAL ERROR  Source  .INCLUDE  file prg/music-engine/dpcm-samples.6502a does not exist!!!

 >>> Assembly WILL NOT continue due to fatal errors! <<<
Assembly could not be completed due to errors!
Total Error Count: 1 / Total Warning Count: 0

*/

// xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx

func AddNew(errorTableKey enumErrorCodes.Def, args ...interface{}) error {
	errData, tableKeyExists := errorTable[errorTableKey]
	if tableKeyExists {
		errMsg := fmt.Sprintf(errData.description, args...)
		entry := NewErrorEntry(errorTableKey, errMsg, errData.severity)

		colorizedFileName := util.Colorize(fmt.Sprintf(" >> %v ", entry.fileName), "red", true)
		fmt.Println("░", colorizedFileName)

		//Line number and content
		colorizedLineNumber := util.Colorize(util.PadStringLeft(fmt.Sprintf(" %d ", entry.lineNumber), 7, ' '), "blue", true)
		fmt.Println("▓", colorizedLineNumber, entry.lineContent)

		severityDescription, severityColor := "", ""

		switch entry.severity {
		case enumErrorCodes.Warning:
			severityColor = "yellow"
			severityDescription = "WARN"
		case enumErrorCodes.Error:
			severityColor = "red"
			severityDescription = "ERROR"
		case enumErrorCodes.Fatal:
			severityColor = "magenta"
			severityDescription = "FATAL ERROR"
		}
		colorizedSeverity := util.Colorize(util.PadStringLeft(fmt.Sprintf(" %v ", severityDescription), 7, ' '), severityColor, true)
		fmt.Println("▓", colorizedSeverity, errMsg)
		return errors.New(fmt.Sprintf("%v%d", SEVERITY_PREFIX, entry.severity))
	}
	return errors.New("Non-error-code error???")
}

// xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx

// A silent error initially...
func AddUnresolved(symbolName string) error {
	return errors.New(fmt.Sprintf("%v%d", SEVERITY_PREFIX, enumErrorCodes.Error))
}

// xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx

func IsErrorCode(err error) (bool, enumErrorCodes.Def) {
	errorCode := err.Error()
	errorTableKey := enumErrorCodes.Def(errorCode)
	_, isValidErrorCode := errorTable[errorTableKey]
	return isValidErrorCode, errorTableKey
}
