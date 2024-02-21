package directiveHandler

import (
	"misc/nintasm/romBuilder"
	"misc/nintasm/romBuilder/nodesToBytes"
)

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
