package directiveHandler

import (
	"misc/nintasm/interpreter/environment/unresolvedTable"
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

	isBigEndian := directiveName == "DWBE"
	operandSize := directiveOperandByteSizes[directiveName]

	for _, operand := range *operandList {
		asRomData, err = nodesToBytes.ConvertNodeValueToUInts(operand, operandSize, isBigEndian)
		if err != nil {
			return err // ❌ Fails
		}
		if !operand.Resolved {
			unresolvedTable.AddUnresolvedRomEntry(operand, operandSize)
		}
		err = romBuilder.AddBytesToRom(asRomData)
		if err != nil {
			return err // ❌ Fails
		}
	}
	return nil
}
