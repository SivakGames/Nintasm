package romBinary

import (
	"errors"
	"fmt"
)

var currentBankIndex = 6
var currentInsertionIndex = 0
var emptyBank = make([]uint8, 128)
var romBinary = make([][]uint8, 8)

func init() {
	for i := range romBinary {
		romBinary[i] = append(romBinary[i], emptyBank...)
	}
}

func AddToRom(insertions []uint8) error {
	currentBank := &romBinary[currentBankIndex]
	toInsertSpace := currentInsertionIndex + len(insertions)
	overflowByteTotal := toInsertSpace - len(*currentBank)

	if overflowByteTotal > 0 {
		errMsg := fmt.Sprintf("Will overflow by: %d byte(s) here", overflowByteTotal)
		return errors.New(errMsg)
	}
	for i := range insertions {
		(*currentBank)[currentInsertionIndex] = insertions[i]
		currentInsertionIndex++
	}
	return nil
}
