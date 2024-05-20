package blockStack

import (
	enumTokenTypes "misc/nintasm/constants/enums/tokenTypes"
	"misc/nintasm/util"
	"strings"
)

// +++++++++++++++++++++++++++++++++++++++++++++
// Whether things should start processing the captured lines
var GoToProcessingFlag bool = false

//================================================

func CheckIfEndOpMatchesOpeningOp(desiredEndOpName string) bool {
	currentStackOp := getCurrentCaptureBlockListNodeCaptureBlockStackTopEntry()
	endOpName, _ := correspondingEndBlockOperations[currentStackOp.BlockOperationName]
	return endOpName == strings.ToUpper(desiredEndOpName)
}

func NEW_CheckEndOpVsStartOp(providedOpName string, startOpName string) bool {
	endOpName, _ := correspondingEndBlockOperations[startOpName]
	return endOpName == strings.ToUpper(providedOpName)
}

// Shared method for most ending labeled directives. Clears the label and pops from the capture block stack
func ProcessEndLabeledDirective() {
	ClearCurrentOperationLabel()
	popFromCurrentCaptureBlockListCaptureBlockStack()
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++

func NEW_IsStartOperation(lineOperationParsedValues *util.LineOperationParsedValues) bool {
	return lineOperationParsedValues.OperationTokenEnum == enumTokenTypes.DIRECTIVE_blockStart ||
		lineOperationParsedValues.OperationTokenEnum == enumTokenTypes.DIRECTIVE_labeledBlockStart
}

func NEW_IsEndOperation(lineOperationParsedValues *util.LineOperationParsedValues) bool {
	return lineOperationParsedValues.OperationTokenEnum == enumTokenTypes.DIRECTIVE_blockEnd ||
		lineOperationParsedValues.OperationTokenEnum == enumTokenTypes.DIRECTIVE_labeledBlockEnd
}

func CheckIfNewStartEndOperation(lineOperationParsedValues *util.LineOperationParsedValues) bool {
	isStartEndEnum := (lineOperationParsedValues.OperationTokenEnum == enumTokenTypes.DIRECTIVE_blockStart ||
		lineOperationParsedValues.OperationTokenEnum == enumTokenTypes.DIRECTIVE_labeledBlockStart ||
		lineOperationParsedValues.OperationTokenEnum == enumTokenTypes.DIRECTIVE_blockEnd ||
		lineOperationParsedValues.OperationTokenEnum == enumTokenTypes.DIRECTIVE_labeledBlockEnd)

	// If op isn't start/end enum no need to proceed...
	if !isStartEndEnum {
		return false
	}

	//If in forced evaluate mode, see if there is a pair to force-close it

	if GetCaptureBlockListForcedCapturingFlag() {
		nnn := GetCurrentOpPtr()
		xx := nnn.overwriteForcedCapture
		_, bbb := (*xx)[strings.ToUpper(lineOperationParsedValues.OperationTokenValue)]
		if bbb {
			return true
		}

		if !CheckIfEndOpMatchesOpeningOp(lineOperationParsedValues.OperationTokenValue) {
			return false
		}
	}

	return true
}

//-----------------------------------------------------

func CheckOperationIsCapturable(
	originalLine string,
	lineOperationParsedValues *util.LineOperationParsedValues) error {
	err := checkOperationIsCapturableByCurrentBlockOperation(lineOperationParsedValues)
	if err != nil {
		return err
	}
	return nil
}

func CheckOperationIsCapturableAndAppend(
	originalLine string,
	lineOperationParsedValues *util.LineOperationParsedValues,
) error {
	err := CheckOperationIsCapturable(originalLine, lineOperationParsedValues)
	if err != nil {
		return err
	}

	currentStackOp := getCurrentCaptureBlockListNodeCaptureBlockStackTopEntryFurthestAlternate()
	currentStackOp.CapturedLines = append(currentStackOp.CapturedLines, newCapturedLine(
		originalLine,
		lineOperationParsedValues.OperationLabel,
		lineOperationParsedValues.OperationTokenEnum,
		lineOperationParsedValues.OperationTokenValue,
		lineOperationParsedValues.OperandStartPosition,
		lineOperationParsedValues.ParentParserEnum,
	))
	return nil
}

//--------------------------------

func CheckIfEndOperationAndGoesToProcessing(lineOperationParsedValues *util.LineOperationParsedValues) bool {
	if (lineOperationParsedValues.OperationTokenEnum == enumTokenTypes.DIRECTIVE_blockEnd ||
		lineOperationParsedValues.OperationTokenEnum == enumTokenTypes.DIRECTIVE_labeledBlockEnd) &&
		GoToProcessingFlag {
		GoToProcessingFlag = false
		return true
	}
	return false
}
