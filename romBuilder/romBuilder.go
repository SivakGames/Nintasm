package romBuilder

import (
	"math"
	"misc/nintasm/assemble/errorHandler"
	enumErrorCodes "misc/nintasm/constants/enums/errorCodes"
	"misc/nintasm/interpreter/operandFactory"
	"misc/nintasm/romBuilder/romBuildingSettings"
)

type Node = operandFactory.Node

type romType []romSegmentType
type romSegmentType []bankType
type bankType struct {
	bytes         []uint8
	orgIsSet      bool
	maxOrg        int
	minOrg        int
	occupiedBytes uint
}

func newBank(bankSize int) bankType {
	fillValue := romBuildingSettings.GetEmptyRomFillValue()
	initBytes := make([]uint8, bankSize)
	for i := range initBytes {
		initBytes[i] = fillValue
	}

	return bankType{bytes: initBytes, orgIsSet: false, minOrg: -1, maxOrg: -1, occupiedBytes: 0}
}

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++

// The final ROM that will be built
var rom = make(romType, 0)

var currentBankIndex = -1
var currentInsertionIndex = -1

//-------------------------------------------

func AddNewRomSegment(totalSize int, bankSize int) error {
	floatResult := float64(totalSize) / float64(bankSize)
	numBanks, deci := math.Modf(floatResult)
	if deci != 0 {
		return errorHandler.AddNew(enumErrorCodes.BankSizeUneven)
	}

	newSegment := make([]bankType, int(numBanks))
	for i := range newSegment {
		newSegment[i] = newBank(bankSize)
	}

	rom = append(rom, newSegment)
	currentBankIndex = -1
	currentInsertionIndex = -1
	return nil
}

// The entire ROM layout (array of segments)
func getRom() *romType {
	return &rom
}

//+++++++++++++++++++++++++++

// The current ROM segment (array of bank segments)
func GetCurrentRomSegment() *romSegmentType {
	if len(rom) > 0 {
		return &rom[len(rom)-1]
	}
	panic("🛑 Attempted to access ROM with no segments!!!")
}

// How many rom segments are currently defined
func GetTotalRomSegmentsInRom() int {
	return len(*getRom())
}

// How many rom segments are currently defined
func GetRomSegmentIndex() int {
	return GetTotalRomSegmentsInRom() - 1
}

//+++++++++++++++++++++++++++

// The current bank segment (array of uint8 )
func GetCurrentBankSegment() *bankType {
	currentRomSegment := GetCurrentRomSegment()
	return &(*currentRomSegment)[currentBankIndex]
}

// The current bank segment (array of uint8 )
func GetCurrentBankSegmentBytes() *[]uint8 {
	currentBankSegment := GetCurrentBankSegment()
	return &(*currentBankSegment).bytes
}

// How many banks are in the current ROM segment
func GetTotalBanksInCurrentRomSegment() int {
	return len(*GetCurrentRomSegment())
}

//+++++++++++++++++++++++++++

func GetBankIndex() int {
	return currentBankIndex
}
func SetBankIndex(newBankIndex int) error {
	currentBankIndex = newBankIndex
	return nil
}

//+++++++++++++++++++++++++++

func GetOrg() int {
	bank := GetCurrentBankSegment()
	return bank.minOrg + currentInsertionIndex
}

func GetCurrentInsertionIndex() int {
	return currentInsertionIndex
}

// +++++++++++++++++++++++++++++

// TODO: Set ORG upper/lower bounds checks
func SetOrg(newOrg int) error {
	bank := GetCurrentBankSegment()
	bankSize := len(bank.bytes)

	if !bank.orgIsSet {
		newMinOrg := int((newOrg / bankSize) * bankSize)
		newMaxOrg := newMinOrg + bankSize - 1
		bank.minOrg = newMinOrg
		bank.maxOrg = newMaxOrg
		bank.orgIsSet = true
	} else {
		currentOrg := bank.minOrg + currentInsertionIndex
		if newOrg < bank.minOrg {
			return errorHandler.AddNew(enumErrorCodes.OrgTooSmall, newOrg, bank.minOrg)
		} else if newOrg > bank.maxOrg {
			return errorHandler.AddNew(enumErrorCodes.OrgTooBig, newOrg, bank.maxOrg)
		} else if newOrg < currentOrg {
			return errorHandler.AddNew(enumErrorCodes.OrgLTProgramCounter, newOrg, currentOrg)
		}
	}
	currentInsertionIndex = newOrg % bankSize
	return nil
}

//============================================================

// Take an array of uint8s and put it the right spot
func AddBytesToRom(insertions []uint8) error {
	currentBankSegment := GetCurrentBankSegment()
	currentBankSegmentBytes := GetCurrentBankSegmentBytes()

	toInsertSpace := currentInsertionIndex + len(insertions)
	overflowByteTotal := toInsertSpace - len(*currentBankSegmentBytes)

	if overflowByteTotal > 0 {
		return errorHandler.AddNew(enumErrorCodes.BankOverflow, overflowByteTotal) //❌☠️ FATAL ERROR
	}
	for i := range insertions {
		(*currentBankSegmentBytes)[currentInsertionIndex] = insertions[i]
		currentInsertionIndex++
		(*currentBankSegment).occupiedBytes++
	}

	return nil
}

//============================================================

// After the main pass, take an array of uint8s and write to a specific spot
func OverwriteResolvedBytesInRom(romSegment int, bank int, offset int, insertions []uint8) {
	byteTarget := &rom[romSegment][bank].bytes
	for i, insertion := range insertions {
		(*byteTarget)[i+offset] = insertion
	}
}

//xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx

func ClearRom() {
	rom = make(romType, 0)
	for i := range rom {
		(rom)[i] = make(romSegmentType, 0)
		for j := range (rom)[i] {
			(rom)[i][j].bytes = nil
			(rom)[i][j].orgIsSet = false
			(rom)[i][j].maxOrg = 0
			(rom)[i][j].minOrg = 0
		}
	}
}
