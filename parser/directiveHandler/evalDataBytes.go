package directiveHandler

import (
	"misc/nintasm/romBuilder"
	"misc/nintasm/romBuilder/nodesToBytes"
)

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
		asRomData, err = nodesToBytes.ConvertNodeValueToUInts(operand, operandSize, isBigEndian)
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
