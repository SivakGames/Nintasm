package parser

import (
	"errors"
	enumTokenTypes "misc/nintasm/enums/tokenTypes"
	"misc/nintasm/parser/operandFactory"
	"misc/nintasm/romBuilder"
	"strings"
)

type DirectiveOperandParser struct {
	OperandParser
	parentLabel string
}

func NewDirectiveOperandParser() DirectiveOperandParser {
	return DirectiveOperandParser{}
}

var directiveAliases = map[string]string{
	"BYTE":        "DB",
	"WORD":        "DW",
	"WORDBE":      "DWBE",
	"REVERSEBYTE": "RDB",
	"PAD":         "DS",
}

var directiveMinMaxOperands = map[enumTokenTypes.Def][2]int{
	enumTokenTypes.DIRECTIVE_dataBytes:  {1, 128},
	enumTokenTypes.DIRECTIVE_mixedData:  {1, 128},
	enumTokenTypes.DIRECTIVE_dataSeries: {1, 2},
	enumTokenTypes.DIRECTIVE_INES:       {1, 1},
}

// Main directive parser
func (p *DirectiveOperandParser) Process(operationType tokenEnum, operationValue string) error {
	var err error

	directiveName := strings.ToUpper(operationValue)
	aliasValue, aliasExists := directiveAliases[directiveName]
	if aliasExists {
		directiveName = aliasValue
	}

	minMaxOperands := directiveMinMaxOperands[operationType]
	minOperands := minMaxOperands[0]
	maxOperands := minMaxOperands[1]

	operandList, err := p.GetOperandList(minOperands, maxOperands)
	if err != nil {
		return err // ❌ Fails
	}

	switch operationType {

	case enumTokenTypes.DIRECTIVE_dataBytes:
		err = evalDataBytesOperands(directiveName, &operandList)
		return err // 🟢/❌ Could be either

	case enumTokenTypes.DIRECTIVE_dataSeries:
		err = evalDataSeriesOperands(directiveName, &operandList)
		return err // 🟢/❌ Could be either

	case enumTokenTypes.DIRECTIVE_mixedData:
		err = evalMixedDataBytesOperands(directiveName, &operandList)
		return err // 🟢/❌ Could be either

	case enumTokenTypes.DIRECTIVE_INES:
		err = evalInesBytesOperands(directiveName, &operandList)
		return err // 🟢/❌ Could be either

	default:
		return errors.New("BAD DIRECTIVE OPERATION TYPE!!!")
	}
}

// +++++++++++++++++++++++++

var directiveOperandByteSizes = map[string]int{
	"DB":   1,
	"DW":   2,
	"DWBE": 2,
	"RDB":  1,
}

// For .db, .dw, .dwbe
func evalDataBytesOperands(directiveName string, operandList *[]Node) error {
	var asRomData = make([]uint8, 0)
	var err error

	isBigEndian := false
	operandSize := directiveOperandByteSizes[directiveName]

	if directiveName == "DWBE" {
		isBigEndian = true
	}

	for _, operand := range *operandList {
		asRomData, err = romBuilder.ConvertNodeValueToUInts(operand, operandSize, isBigEndian)
		if err != nil {
			return err // ❌ Fails
		}
		err = romBuilder.AddBytesToRom(asRomData)
		if err != nil {
			return err // ❌ Fails
		}
	}
	return nil
}

// +++++++++++++++++++++++++

// For .ds
func evalDataSeriesOperands(directiveName string, operandList *[]Node) error {
	seriesValue := uint8(0xff)
	repetitionNode := &(*operandList)[0]
	if !(operandFactory.ValidateNodeIsNumeric(repetitionNode) &&
		operandFactory.ValidateNumericNodeIsGTZero(repetitionNode)) {
		return errors.New("DS/PAD directive repeat value must be a number that is > 0")
	}

	repetitionNumber := repetitionNode.AsNumber
	numRepetitions := repetitionNumber

	if len(*operandList) == 2 {
		padNode := &(*operandList)[1]
		if !(operandFactory.ValidateNodeIsNumeric(padNode) &&
			operandFactory.ValidateNumericNodeIsPositive(padNode) &&
			operandFactory.ValidateNumericNodeIs8BitValue(padNode)) {
			return errors.New("DS/PAD directive fill value must be a non-negative 8-bit number")
		}
		seriesValue = uint8(padNode.AsNumber)
	}

	asRomData := make([]uint8, numRepetitions)
	for i := range asRomData {
		asRomData[i] = seriesValue
	}

	err := romBuilder.AddBytesToRom(asRomData)
	if err != nil {
		return err // ❌ Fails
	}

	return nil
}

// +++++++++++++++++++++++++

type mixedDataDirectiveBytesKeyFormat struct {
	numBytes  int
	bigEndian bool
}

var mixedDataDirectiveBytesKeys = map[byte]mixedDataDirectiveBytesKeyFormat{
	'B': {1, false},
	'W': {2, false},
	'E': {2, true},
}

// For .d_***_
func evalMixedDataBytesOperands(directiveName string, operandList *[]Node) error {
	var asRomData = make([]uint8, 0)
	var err error

	mixedPattern := directiveName[2:]
	lastRepeats := mixedPattern[len(mixedPattern)-1] == '_'
	if lastRepeats {
		mixedPattern = mixedPattern[:len(mixedPattern)-1]
	}

	for i, operand := range *operandList {
		var mixedPatternIndex int
		if lastRepeats && i >= len(mixedPattern) {
			mixedPatternIndex = len(mixedPattern) - 1
		} else {
			mixedPatternIndex = i % len(mixedPattern)
		}

		currentPatternKey := mixedPattern[mixedPatternIndex]
		operandSize := mixedDataDirectiveBytesKeys[currentPatternKey].numBytes
		isBigEndian := mixedDataDirectiveBytesKeys[currentPatternKey].bigEndian

		asRomData, err = romBuilder.ConvertNodeValueToUInts(operand, operandSize, isBigEndian)
		if err != nil {
			return err // ❌ Fails
		}
		err = romBuilder.AddBytesToRom(asRomData)
		if err != nil {
			return err // ❌ Fails
		}
	}
	return nil
}

// +++++++++++++++++++++++++

func evalInesBytesOperands(directiveName string, operandList *[]Node) error {
	var err error
	inesNode := &(*operandList)[0]

	switch directiveName {
	case "INESPRG":
		err = romBuilder.ValidateInesPrg(inesNode)
	case "INESCHR":
		err = romBuilder.ValidateInesChr(inesNode)
	case "INESMAP":
		err = romBuilder.ValidateInesMap(inesNode)
	case "INESMIR":
		err = romBuilder.ValidateInesMirroring(inesNode)
	case "INESBAT":
	default:
		panic("Something is very wrong with ines directive")
	}

	return err
}
