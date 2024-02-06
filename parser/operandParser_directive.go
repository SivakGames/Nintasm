package parser

import (
	"errors"
	enumNodeTypes "misc/nintasm/enums/nodeTypes"
	enumTokenTypes "misc/nintasm/enums/tokenTypes"
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

var directiveMinMaxOperands = map[string][2]int{
	"DB":   {1, 128},
	"DW":   {1, 128},
	"DWBE": {1, 128},
	"RDB":  {1, 128},
	"DS":   {1, 2},
}

var directiveOperandByteSizes = map[string]int{
	"DB":   1,
	"DW":   2,
	"DWBE": 2,
	"RDB":  1,
}

func (p *DirectiveOperandParser) Process(operationType tokenEnum, operationValue string) error {
	var evalOperandSize int
	var err error

	directiveName := strings.ToUpper(operationValue)
	aliasValue, aliasExists := directiveAliases[directiveName]
	if aliasExists {
		directiveName = aliasValue
	}

	isBigEndian := false
	minMaxOperands := directiveMinMaxOperands[directiveName]
	minOperands := minMaxOperands[0]
	maxOperands := minMaxOperands[1]

	operandList, err := p.GetOperandList(minOperands, maxOperands)
	if err != nil {
		return err // ❌ Fails
	}

	switch operationType {

	case enumTokenTypes.DIRECTIVE_data:
		if directiveName == "DS" {
			err = evalDataPadOperands(&operandList)
			return err // 🟢/❌ Could be either

		}

		evalOperandSize = directiveOperandByteSizes[directiveName]
		if directiveName == "DWBE" {
			isBigEndian = true
		}
		err = evalDataInsertionOperands(&operandList, evalOperandSize, isBigEndian)
		if err != nil {
			return err // ❌ Fails
		}

	default:
		return errors.New("BAD DIRECTIVE OPERATION TYPE!!!")
	}

	return nil

}

// For .db, .dw, etc.
func evalDataInsertionOperands(operandList *[]Node, operandSize int, isBigEndian bool) error {
	var asRomData = make([]uint8, 0)
	var err error

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

// For .ds
func evalDataPadOperands(operandList *[]Node) error {
	padValue := 0xff

	repetitionNode := &(*operandList)[0]
	if repetitionNode.NodeType != enumNodeTypes.NumericLiteral {
		return errors.New("DS/PAD directive Must be a number!")
	}
	repetitionNumber := repetitionNode.AsNumber
	if repetitionNumber < 1 {
		return errors.New("DS/PAD directive cannot less than 1!")
	}
	numRepetitions := repetitionNumber

	if len(*operandList) == 2 {
		padNode := &(*operandList)[1]
		if padNode.NodeType != enumNodeTypes.NumericLiteral {
			return errors.New("DS/PAD directive Must be a number!")
		}
		padNodeNumber := padNode.AsNumber
		if padNodeNumber < 0 {
			return errors.New("DS/PAD directive cannot be negative!")
		}
		padValue = padNodeNumber
	}

	asRomData := make([]uint8, numRepetitions)
	for i := range asRomData {
		asRomData[i] = uint8(padValue)
	}

	err := romBuilder.AddBytesToRom(asRomData)
	if err != nil {
		return err // ❌ Fails
	}

	return nil
}
